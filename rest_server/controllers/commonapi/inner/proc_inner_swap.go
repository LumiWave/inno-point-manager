package inner

import (
	"math"
	"strconv"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func PutSwapStatus(params *context.ReqSwapStatus) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	if swapInfo, err := model.GetDB().CacheGetSwapWallet(params.FromWalletAddress); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_GetSwapInfo])
		resp.SetReturn(resultcode.Result_RedisError_GetSwapInfo)
	} else {
		switch params.TxStatus {
		case context.SWAP_status_fee_transfer_start, context.SWAP_status_fee_transfer_success: // swap 수수료 전송 시작
			// 콜백이 먼저 들어와서 상태가 진행 된 경우는 버린다.
			if swapInfo.TxStatus < params.TxStatus {
				swapInfo.TxStatus = params.TxStatus
				swapInfo.TxHash = params.TxHash
				if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
					log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
					resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
				} else {
					if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_Exchangefee(swapInfo.TxID,
						params.TxStatus,
						params.TxHash,
						strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64), swapInfo.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
						resp.SetReturn(resultcode.Result_Error_Db_TransactExchangeGoods_Gasfee)
					}
				}
			} else {
				log.Warnf("swap not equal status redis:%v, rev:%v", swapInfo.TxStatus, params.TxStatus)
			}
		case context.SWAP_status_token_transfer_start: // swap용 토큰 전송 시작 ( coin->point swap)
			swapInfo.TxStatus = params.TxStatus
			swapInfo.SwapCoin.TokenTxHash = params.TxHash
			if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
				log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
				resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
			} else {
				if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_Coin(swapInfo.TxID,
					params.TxStatus,
					params.TxHash,
					time.Now().Format("2006-01-02 15:04:05.000"), 0, ""); err != nil {
					resp.SetReturn(resultcode.Result_Error_Db_TransactExchangeGoods_Gasfee)
				}
			}
		case context.SWAP_status_fee_transfer_fail, context.SWAP_status_token_transfer_fail:
			if err := model.GetDB().CacheDelSwapWallet(params.FromWalletAddress); err != nil { // swap 수수료 전송 실패
				log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
				resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
			} else {
				// 포인트 누적이 연속적으로 처리 되지 못하도록 한다. lock
				Lockkey := model.MakeMemberPointListLockKey(swapInfo.MUID)
				mutex := model.GetDB().RedSync.NewMutex(Lockkey)
				if err := mutex.Lock(); err != nil {
					log.Error("redis lock err:%v", err)
					resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
					return resp
				}

				defer func() {
					// 1-1. redis unlock
					if ok, err := mutex.Unlock(); !ok || err != nil {
						if err != nil {
							log.Errorf("unlock err : %v", err)
						}
					}
				}()

				// 실패 완료 처리
				// 최신 포인트 수량을 가져와서 복원할 포인트 정보를 다시 계산해서 완료 처리 한다.
				if _, points, err := model.GetDB().USPPO_GetList_MemberPoints(swapInfo.MUID, swapInfo.DatabaseID); err != nil {
					log.Errorf("GetPointAppList error : %v", err)
					resp.SetReturn(resultcode.Result_Error_DB_GetPointAppList)
					return resp
				} else {
					if point, ok := points[swapInfo.PointID]; ok {
						swapInfo.PreviousPointQuantity = point.Quantity
						swapInfo.AdjustPointQuantity = -swapInfo.AdjustPointQuantity
						swapInfo.PointQuantity = swapInfo.PreviousPointQuantity + swapInfo.AdjustPointQuantity
					}
				}
				if err := model.GetDB().USPAU_Cmplt_ExchangeGoods(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), false); err != nil {
					resp.SetReturn(resultcode.Result_Error_Db_Swap_Complete)
				}
			}
		}
	}
	return resp
}

func SwapWallet(params *context.ReqSwapInfo, innoUID string) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	// 0. 포인트 누적이 연속적으로 처리 되지 못하도록 한다.
	// 2. 부모지갑에 수수료 전송 중인지 체크
	// 3. redis에 해당 포인트 정보 존재하는지 check, 있으면 강제로 db에 마지막 정보 업데이트 후 swap 진행
	// 4. 전환 정보 검증
	// 5. point->coin 시 부모지갑에 수수료 전송
	// 6. coin->point 시 부모지갑에 코인 전송
	// 6. swap 정보 redis 저장
	// 00. 부모입금 callback 기다림

	// 0. 포인트 누적이 연속적으로 처리 되지 못하도록 한다.
	Lockkey := model.MakeMemberPointListLockKey(params.MUID)
	mutex := model.GetDB().RedSync.NewMutex(Lockkey)
	if err := mutex.Lock(); err != nil {
		log.Error("redis lock err:%v", err)
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	}

	defer func() {
		// 1-1. redis unlock
		if ok, err := mutex.Unlock(); !ok || err != nil {
			if err != nil {
				log.Errorf("unlock err : %v", err)
			}
		}
	}()

	// 2. 부모지갑에 수수료 전송 중인지 체크 한사람당 한번에 한 swap 가능 하도록 막는다.
	if _, err := model.GetDB().CacheGetSwapWallet(params.WalletAddress); err == nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_Inprogress])
		resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
		return resp
	}

	// 3. redis에 해당 포인트 정보 존재하는지 check
	// 있으면 강제로 db에 마지막 정보 업데이트 후 swap 진행 : 게임사에서 포인트 쌓을때 충돌 방지
	pointKey := model.MakeMemberPointListKey(params.MUID)
	mePointInfo, err := model.GetDB().GetCacheMemberPointList(pointKey)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointAppList(params.MUID, params.DatabaseID); err != nil {
			log.Errorf("GetPointAppList error : %v", err)
			resp.SetReturn(resultcode.Result_Error_DB_GetPointAppList)
			return resp
		} else {
			for _, point := range points {
				if point.PointID == params.PointID {
					params.PreviousPointQuantity = point.Quantity
					params.PointQuantity = params.PreviousPointQuantity + params.AdjustPointQuantity
					break
				}
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
					point.PreQuantity, point.AdjustQuantity, point.Quantity, context.LogID_cp, eventID); err != nil {
					log.Errorf("UpdateAppPoint error : %v", err)
					resp.SetReturn(resultcode.Result_Error_DB_UpdateAppPoint)
					return resp
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
			if params.PointID == point.PointID && params.MUID == mePointInfo.MUID {
				params.PreviousPointQuantity = point.Quantity
				params.PointQuantity = params.PreviousPointQuantity + params.AdjustPointQuantity
			}
		}

		model.GetDB().DelCacheMemberPointList(pointKey)
	}

	// 로그를 위한 wallet ID 조회
	if _, wallets, err := model.GetDB().USPAU_GetList_AccountWallets(params.AUID); err == nil {
		if wallet, ok := wallets[params.WalletAddress]; ok {
			params.WalletID = wallet.WalletID
		} else {
			log.Errorf("Not exist available wallet auid: %v", params.AUID)
			resp.SetReturn(resultcode.Result_Error_Db_GetAccountWallets)
			return resp
		}
	}

	// 4. 전환 정보 검증
	pointInfo := model.GetDB().AppPointsMap[params.AppID].PointsMap[params.PointID]
	if params.TxType == context.EventID_toCoin {
		// 코인으로 전환시 체크
		// 당일 누적 코인 전환 수량이 넘었는지 체크
		if _, coinsMap, err := model.GetDB().GetAccountCoins(params.AUID); err != nil {
			log.Errorf("GetAccountCoins error : %v", err)
			resp.SetReturn(resultcode.Result_DBError)
			return resp
		} else {
			if val, ok := coinsMap[params.CoinID]; ok {
				if val.TodayAcqExchangeQuantity+params.AdjustCoinQuantity > model.GetDB().Coins[params.CoinID].DailyLimitedAcqExchangeQuantity {
					// error
					log.Errorf("Result_Error_Exceed_DailyLimitedSwapCoin auid:%v", params.AUID)
					resp.SetReturn(resultcode.Result_Error_Exceed_DailyLimitedSwapCoin)
					return resp
				}
			} else {
				// 내 지갑 코인 정보에 데이터가 없다는것은 최초 스왑인경우 이고 무조건 성공 처리해준다.
			}
		}
		// 포인트 보유수량이 전환량 보다 큰지 확인
		absAdjustPointQuantity := int64(math.Abs(float64(params.AdjustPointQuantity)))
		if params.PreviousPointQuantity <= 0 || // 보유 포인트량이 0일경우
			params.PreviousPointQuantity < params.AdjustPointQuantity || // 전환 할 수량보다 보유 수량이 적을 경우
			pointInfo.MinExchangeQuantity > absAdjustPointQuantity { // 전환 최소 수량 에러
			// 전환할 포인트 수량이 없음 에러
			log.Errorf("lack of minimum point quantity [point_id:%v][PointQuantity:%v]", params.PointID, params.PreviousPointQuantity)
			resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
			return resp
		}
		// 전환 비율 계산 후 타당성 확인
		exchangeCoin := float64(absAdjustPointQuantity) * pointInfo.ExchangeRatio
		exchangeCoin = toFixed(exchangeCoin, 4)
		if params.AdjustCoinQuantity != exchangeCoin {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return resp
		}
	} else if params.TxType == context.EventID_toPoint {
		// 당일 누적 포인트 전환 최대 수량이 넘었는지 체크
		if accountPoint, err := model.GetDB().GetListAccountPoints(0, params.MUID); err != nil {
			log.Errorf("GetListAccountPoints error : %v", err)
			resp.SetReturn(resultcode.Result_DBError)
			return resp
		} else {
			if val, ok := accountPoint[params.PointID]; ok {
				if val.TodayAcqExchangeQuantity+params.AdjustPointQuantity > model.GetDB().AppPointsMap[params.AppID].PointsMap[params.PointID].DailyLimitedAcqExchangeQuantity {
					// error
					log.Errorf("Result_Error_Exceed_DailyLimitedSwapPoint auid:%v", params.AUID)
					resp.SetReturn(resultcode.Result_Error_Exceed_DailyLimitedSwapPoint)
					return resp
				}
			}
		}

		absAdjustCoinQuantity := math.Abs(params.AdjustCoinQuantity)
		// 전환 비율 계산 후 타당성 확인
		exchangePoint := absAdjustCoinQuantity / pointInfo.ExchangeRatio
		exchangePoint = toFixed(exchangePoint, 0)
		if params.AdjustPointQuantity != int64(exchangePoint) {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToCoin)
			return resp
		}
	}

	if txID, err := model.GetDB().USPAU_Strt_ExchangeGoods(params); err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
		return resp
	} else {
		params.TxID = *txID
		params.CreateAt = time.Now().UTC().Unix()
		params.ToWalletAddress = config.GetInstance().ParentWalletsMapBySymbol[params.BaseCoinSymbol].ParentWalletAddr
		params.TxStatus = context.SWAP_status_init

		if err := model.GetDB().CacheSetSwapWallet(params); err != nil {
			log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
			resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
			return resp
		}
	}

	resp.Value = params
	return resp
}

func round(num float64) int {
	return int(num + math.Copysign(0, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
