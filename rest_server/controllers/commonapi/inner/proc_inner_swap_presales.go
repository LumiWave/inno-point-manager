package inner

import (
	"math"
	"strconv"
	"time"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

func PutPreSalesSwapStatus(params *context.ReqSwapStatus) *base.BaseResponse {
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
					if err := model.GetDB().USPAU_Mod_TransactExchanges_ExchangeFees(swapInfo.TxID,
						params.TxStatus,
						params.TxHash,
						strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64), swapInfo.SwapToCoin.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
						resp.SetReturn(resultcode.Result_Error_Db_TransactExchangeGoods_Gasfee)
					}
				}
			} else {
				log.Warnf("swap not equal status redis:%v, rev:%v", swapInfo.TxStatus, params.TxStatus)
			}
		case context.SWAP_status_token_transfer_deposit_start: // swap용 토큰 전송 시작 ( coin->point, coin->coin swap)
			swapInfo.TxStatus = params.TxStatus
			swapInfo.SwapFromCoin.TokenTxHash = params.TxHash
			if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
				log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
				resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
			} else {
				if err := model.GetDB().USPAU_Mod_TransactExchanges_Coin(swapInfo.TxID,
					params.TxStatus,
					params.TxHash,
					time.Now().Format("2006-01-02 15:04:05.000"), 0, ""); err != nil {
					resp.SetReturn(resultcode.Result_Error_Db_TransactExchangeGoods_Gasfee)
				}
			}
		case context.SWAP_status_fee_transfer_fail, context.SWAP_status_token_transfer_deposit_fail:
			// 유져 지갑을 통해 전송 시도가 실패한 경우를 수신 받았을때 포인트 원복을 해줘야 함
			swapInfo, err := model.GetDB().CacheGetSwapWallet(params.FromWalletAddress)
			if err != nil {
				log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
				resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
			} else {
				if err := model.GetDB().CacheDelSwapWallet(params.FromWalletAddress); err != nil { // swap 수수료 전송 실패
					log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetSwapInfo])
					resp.SetReturn(resultcode.Result_RedisError_SetSwapInfo)
				} else {
					// 포인트 누적이 연속적으로 처리 되지 못하도록 한다. lock
					swapPoint := &context.SwapPoint{}
					if swapInfo.TxType == context.EventID_P2C {
						swapPoint = &swapInfo.SwapFromPoint
					} else if swapInfo.TxType == context.EventID_C2P {
						swapPoint = &swapInfo.SwapToPoint
					}
					Lockkey := model.MakeMemberPointListLockKey(swapPoint.MUID)
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
					if _, points, err := model.GetDB().USPPO_GetList_MemberPoints(swapPoint.MUID, swapPoint.DatabaseID); err != nil {
						log.Errorf("GetPointAppList error : %v", err)
						resp.SetReturn(resultcode.Result_Error_DB_GetPointAppList)
						return resp
					} else {
						if point, ok := points[swapPoint.PointID]; ok {
							swapPoint.PreviousPointQuantity = point.Quantity
							swapPoint.AdjustPointQuantity = -swapPoint.AdjustPointQuantity
							swapPoint.PointQuantity = swapPoint.PreviousPointQuantity + swapPoint.AdjustPointQuantity
						}
					}
					if err := model.GetDB().USPAU_Cmplt_Exchanges(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), false); err != nil {
						resp.SetReturn(resultcode.Result_Error_Db_Swap_Complete)
					}
				}
			}
		}
	}
	return resp
}

func PreSalesSwapWallet(params *context.ReqSwapInfo, innoUID string) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	// 0. 포인트 누적이 연속적으로 처리 되지 못하도록 한다.
	// 2. 부모지갑에 수수료 전송 중인지 체크
	// 3. redis에 해당 포인트 정보 존재하는지 check, 있으면 강제로 db에 마지막 정보 업데이트 후 swap 진행
	// 4. 전환 정보 검증
	// 5. point->coin 시 부모지갑에 수수료 전송
	// 6. coin->point 시 부모지갑에 코인 전송
	// 7. coin->coin 시 부모지갑에 수수료 전송 후 부모지갑에 코인 전송
	// 8. point->point 시에는 콜백 없이 바로 처리한다.
	// 8. swap 정보 redis 저장
	// 9. 부모입금 callback 기다림

	// 0. 포인트 누적이 연속적으로 처리 되지 못하도록 한다. P2C, C2P만 해당함
	if params.TxType == context.EventID_Server_toC2P {
		muid := params.SwapToPoint.MUID

		Lockkey := model.MakeMemberPointListLockKey(muid)
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
	}

	// 2. 부모지갑에 수수료 혹은 코인 전송 중인지 체크, 한사람당 한번에 한 swap 가능 하도록 막는다.
	checkAlreadySwap(params, resp)
	if resp.Return != 0 {
		return resp
	}

	// 3. redis에 해당 포인트 정보 존재하는지 check
	// 있으면 강제로 db에 마지막 정보 업데이트 후 swap 진행 : 게임사에서 포인트 쌓을때 충돌 방지
	savePoint(params, resp)
	if resp.Return != 0 {
		return resp
	}

	// 4. 전환 정보 검증
	swapBaseInfo := model.GetDB().SwapAblePreSalesMap[params.SwapFromCoin.CoinID][params.SwapToPoint.PointID]
	// 당일 누적 포인트 전환 최대 수량이 넘었는지 체크
	if accountPoint, err := model.GetDB().GetListAccountPoints(params.AUID); err != nil {
		log.Errorf("GetListAccountPoints error : %v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return resp
	} else {
		if val, ok := accountPoint[params.SwapToPoint.PointID]; ok {
			if val.TodayExchangeAcqQuantity+params.SwapToPoint.AdjustPointQuantity > model.GetDB().AppPointsMap[params.SwapToPoint.AppID].PointsMap[params.SwapToPoint.PointID].DailyLimitExchangeAcqQuantity {
				// error
				log.Errorf("Result_Error_Exceed_DailyLimitedSwapPoint auid:%v", params.AUID)
				resp.SetReturn(resultcode.Result_Error_Exceed_DailyLimitedSwapPoint)
				return resp
			}
		}
	}

	absAdjustCoinQuantity := math.Abs(params.SwapFromCoin.AdjustCoinQuantity)
	// 전환 비율 계산 후 타당성 확인
	exchangePoint := absAdjustCoinQuantity * swapBaseInfo.ExchangeRatio
	exchangePoint = toFixed(exchangePoint, 0)
	if params.SwapToPoint.AdjustPointQuantity != int64(exchangePoint) {
		resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToCoin)
		return resp
	}

	// p2p가 아닌 작업은 시작을 해두고 콜백이 기다려서 처리해준다.
	if txID, err := model.GetDB().USPPR_Strt_PreSalesExchanges(
		swapBaseInfo.SalesID,
		params.AUID,
		params.SwapToPoint.MUID,
		params.SwapFromCoin.BaseCoinID,
		params.SwapFromCoin.WalletTypeID,
		params.SwapFromCoin.WalletID,
		params.SwapFromCoin.WalletAddress,
		params.SwapFromCoin.CoinID,
		strconv.FormatFloat(params.SwapFromCoin.AdjustCoinQuantity, 'f', -1, 64),
		params.SwapToPoint.PointID,
		strconv.FormatInt(params.SwapToPoint.AdjustPointQuantity, 10)); err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
		return resp
	} else {
		params.TxID = txID
		params.CreateAt = time.Now().UTC().Unix()

		// from coin 구조체에 부모 지갑 주소를 넣어서 해당 주소로 전송 유도 한다.
		params.SwapFromCoin.ToWalletAddress = config.GetInstance().ParentWalletsMapBySymbol[params.SwapFromCoin.BaseCoinSymbol].ParentWalletAddr

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
