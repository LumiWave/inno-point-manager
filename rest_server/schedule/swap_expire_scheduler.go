package schedule

import (
	"sync"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

var gSwapExpireScheduler *SwapExpireScheduler
var onceSwapExpireScheduler sync.Once

type SwapExpireScheduler struct {
	Running     bool  //true:스케쥴실행중 , false:스케쥴중지
	DebugMode   bool  //중간중간 로그찍을부분이있을때 true
	ExpireCycle int64 // 만료 시간 second
}

func InitSwapExpireScheduler(conf *config.ServerConfig) *SwapExpireScheduler {
	schedule, ok := conf.ScheduleMap["swap_expire_scheduler"]
	if ok && schedule.Enable {
		onceSwapExpireScheduler.Do(func() {
			gSwapExpireScheduler = new(SwapExpireScheduler)
			gSwapExpireScheduler.Running = true
			gSwapExpireScheduler.ExpireCycle = schedule.ExpireCycle
			gSwapExpireScheduler.Run(schedule.TermSec)
		})
	}
	return gSwapExpireScheduler
}

func (o *SwapExpireScheduler) SetDebugMode(enable bool) {
	o.DebugMode = enable
}

func (o *SwapExpireScheduler) SetRunning(enable bool) {
	o.Running = enable
}

func (o *SwapExpireScheduler) Run(sec int64) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	go func() {
		for t := range ticker.C {
			if o.DebugMode {
				log.Debugf("SwapExpireScheduler : %v", t)
			}
			if o.Running {
				o.ScheduleProcess()
			}
		}
	}()
}

func (o *SwapExpireScheduler) ScheduleProcess() {
	// redis의 "SWAP-WALLET" 정보들 중에 일정 시간이 지나도 '수수료 입금 시작' 진행이 되고있지 않은 내역은 만료 처리 한다.
	//startTime := time.Now().UnixMilli()
	_, list, err := model.GetDB().CacheGetSwapWallets()
	if err != nil {
		log.Errorf("CacheGetSwapWallets err : %v", err)
		return
	}

	for _, value := range list {
		// 수수료 전송 시작 상태가 아닌 정보중에 10분이 지난 정보는 swap 종료하고 삭제 처리한다.
		if value.CreateAt+o.ExpireCycle < time.Now().UTC().Unix() && value.TxStatus < context.SWAP_status_fee_transfer_start {
			log.Debugf("swap expire addr : fromWallet:%v, toWallet:%v, time:%v", value.SwapFromCoin.WalletAddress, value.SwapToCoin.WalletAddress, time.Unix(value.CreateAt, 0).Format(time.RFC3339))

			if value.TxType == context.EventID_P2C ||
				value.TxType == context.EventID_C2P {
				// 현재 레디스에 포인트가 쌓이고 있을수 있으니 최종값으로 디비에 저장하고 스왑 포인트 복구 처리 해준다

				swapPoint := func() *context.SwapPoint {
					if value.TxType == context.EventID_P2C {
						return &value.SwapFromPoint
					} else if value.TxType == context.EventID_C2P {
						return &value.SwapToPoint
					}
					log.Errorf("invalid swap type : %v", value.TxType)
					return &context.SwapPoint{}
				}()
				pointKey := model.MakeMemberPointListKey(swapPoint.MUID)
				mePointInfo, err := model.GetDB().GetCacheMemberPointList(pointKey)
				if err != nil {
					if _, points, err := model.GetDB().USPPO_GetList_MemberPoints(swapPoint.MUID, swapPoint.DatabaseID); err != nil {
						log.Errorf("GetPointAppList error : %v", err)
					} else {
						if point, ok := points[swapPoint.PointID]; ok {
							swapPoint.PreviousPointQuantity = point.Quantity
							swapPoint.AdjustPointQuantity = -swapPoint.AdjustPointQuantity
							swapPoint.PointQuantity = swapPoint.PreviousPointQuantity + swapPoint.AdjustPointQuantity
						} else {
							log.Errorf("not file swap point id : %v, points : %v", swapPoint.PointID, points)
						}
					}
				} else {
					// redis에 존재 한다면 강제로 db에 먼저 write
					for _, point := range mePointInfo.Points {
						var eventID context.EventID_type
						if point.AdjustQuantity >= 0 {
							eventID = context.EventID_add
						} else {
							eventID = context.EventID_sub
						}

						if point.AdjustQuantity != 0 {
							if todayAcqQuantity, resetDate, err := model.GetDB().UpdateAppPoint(mePointInfo.DatabaseID, mePointInfo.MUID, point.PointID,
								point.PreQuantity, point.AdjustQuantity, point.Quantity, context.LogID_exchange, eventID); err != nil {
								log.Errorf("UpdateAppPoint error : %v", err)
							} else {
								//현재 일일 누적량, 날짜 업데이트
								point.TodayQuantity = todayAcqQuantity
								point.ResetDate = resetDate

								point.AdjustQuantity = 0
								point.PreQuantity = point.Quantity
							}
						} else {
							point.AdjustQuantity = 0
							point.PreQuantity = point.Quantity
						}

						// swap point quantity에 업데이트
						if swapPoint.PointID == point.PointID && swapPoint.MUID == mePointInfo.MUID {
							swapPoint.PreviousPointQuantity = point.Quantity
							swapPoint.AdjustPointQuantity = -swapPoint.AdjustPointQuantity
							swapPoint.PointQuantity = swapPoint.PreviousPointQuantity + swapPoint.AdjustPointQuantity
						}
					}

					model.GetDB().DelCacheMemberPointList(pointKey)
				}
			} else if value.TxType == context.EventID_C2C { // c2c는 양쪽 모두 그냥 지워주면 끝이다.
				model.GetDB().CacheDelSwapWallet(value.SwapFromCoin.WalletAddress)
				model.GetDB().CacheDelSwapWallet(value.SwapToCoin.WalletAddress)
			} else if value.TxType == context.EventID_P2P { // p2p는 콜백없이 그냥 바로 성공 실패가 정해지기 때문에 레디스에 남지 않아서 따로 처리할게 없다.
			}

			if err := model.GetDB().USPAU_Cmplt_Exchanges(value, time.Now().Format("2006-01-02 15:04:05.000"), false); err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges err : %v, txid:%v fromwallet:%v, towallet:%v", err, value.TxID, value.SwapFromCoin.WalletAddress, value.SwapToCoin.WalletAddress)
			} else {
				walletAddr := ""
				if value.TxType == context.EventID_P2C {
					walletAddr = value.SwapToCoin.WalletAddress
				} else if value.TxType == context.EventID_C2P || value.TxType == context.EventID_C2C {
					walletAddr = value.SwapFromCoin.WalletAddress
				}
				if err = model.GetDB().CacheDelSwapWallet(walletAddr); err != nil {
					log.Errorf("CacheDelSwapWallet err:%v, wallet:%v", err, walletAddr)
				}
			}
		}
	}
	//log.Debugf("swap expire checktime :%v", time.Now().UnixMilli()-startTime)
}
