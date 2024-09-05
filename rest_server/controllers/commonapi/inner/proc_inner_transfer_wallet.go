package inner

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/api_inno_log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
	"github.com/LumiWave/inno-point-manager/rest_server/util"
)

func TransferResultWithdrawalWallet(fromAddr, toAddr, value, fee, symbol, txHash, status string, decimal int) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	tKey := model.MakeCoinTransferKeyByTxID(txHash)
	txType, err := model.GetDB().GetCacheCoinTransferTx(tKey)
	if err != nil {
		// 존재 하지 않는 출금 정보 콜백을 받았다.
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Invalid_transfer_txid]+" txid:%v from:%v to:%v amount:%v",
			txHash, fromAddr, toAddr, value)
		resp.SetReturn(resultcode.Result_Invalid_transfer_txid)
		return resp
	}

	// 실패한 tx면 관련 redis 찾아서 삭제 한다.
	if !strings.EqualFold(status, "success") {
		Lockkey := model.MakeCoinTransferFromUserWalletLockKey(txType.AUID)
		mutex := model.GetDB().RedSync.NewMutex(Lockkey)
		if err := mutex.Lock(); err != nil {
			log.Error("redis lock err:%v", err)
			resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
			return resp
		} else {
			defer func() {
				// 1-1. redis unlock
				if ok, err := mutex.Unlock(); !ok || err != nil {
					if err != nil {
						log.Errorf("unlock err : %v", err)
					}
				}
			}()
		}

		log.Errorf("transfer txid fail txid:%v", txHash)

		if txType.Target == context.From_parent_to_other_wallet {
			parentKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
			_, err := model.GetDB().GetCacheCoinTransferFromParentWallet(parentKey)
			if err != nil {
				log.Errorf("GetCacheCoinTransferFromParentWallet err2 => %v auid:%v txid:%v", err, txType.AUID, txHash)
				return resp
			}

			// swap redis 찾아서 완료 처리 하기
			swapInfo, err := model.GetDB().CacheGetSwapWallet(toAddr)
			if err != nil {
				log.Warnf("not exist fromAddr : %v, txHash:%v", fromAddr, txHash)
				return resp
			}

			fe := util.ToDecimalEncStr(fee, int64(decimal))

			swapInfo.TxStatus = context.SWAP_status_token_transfer_withdrawal_fail

			swapCoin := &context.SwapCoin{}
			if swapInfo.TxType == context.EventID_P2C {
				swapCoin = &swapInfo.SwapToCoin
			} else if swapInfo.TxType == context.EventID_C2P {
				swapCoin = &swapInfo.SwapFromCoin
			}
			if err := model.GetDB().USPAU_Mod_TransactExchanges_TxStatus(swapCoin.BaseCoinID, fe, swapInfo); err == nil {
				// swap 토큰 전송 실패난 경우 디비에만 실패 처리 해두고 레디스 그대로 두고 cs 처리 유도한다.
			}
		}

		return resp
	}

	// 성공 분기처리
	if txType.Target == context.From_user_to_fee_wallet { // not used
	} else if txType.Target == context.From_parent_to_other_wallet { // 부모지갑에서 자식지갑으로 코인 전송 : swap point->coin, 혹은 부모 지갑에서 출금
		// 부모지갑에서 자식지갑으로 입금시 swap으로 간주하고 swap 처리 진행한다.
		parentKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
		_, err := model.GetDB().GetCacheCoinTransferFromParentWallet(parentKey)
		if err != nil {
			log.Errorf("GetCacheCoinTransferFromParentWallet err2 => %v auid:%v txid:%v", err, txType.AUID, txHash)
			return resp
		} else {
			// parent가 보낸 redis 정보는 삭제 해준다.
			model.GetDB().DelCacheCoinTransfer(tKey)                      //tx 삭제
			model.GetDB().DelCacheCoinTransferFromParentWallet(parentKey) // from parent 삭제
		}

		// swap redis 찾아서 완료 처리 하기
		swapInfo, err := model.GetDB().CacheGetSwapWallet(toAddr)
		if err != nil {
			log.Warnf("not exist fromAddr : %v, txHash:%v", fromAddr, txHash)
			return resp
		}

		fe := util.ToDecimalEncStr(fee, int64(decimal))

		swapInfo.TxStatus = context.SWAP_status_token_transfer_withdrawal_success

		swapCoin := &context.SwapCoin{}
		if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_C2C {
			swapCoin = &swapInfo.SwapToCoin
		} else if swapInfo.TxType == context.EventID_C2P {
			swapCoin = &swapInfo.SwapFromCoin
		}

		if err := model.GetDB().USPAU_Mod_TransactExchanges_TxStatus(swapCoin.BaseCoinID, fe, swapInfo); err == nil {
			if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_P2C {
				model.GetDB().CacheDelSwapWallet(toAddr)
			} else if swapInfo.TxType == context.EventID_C2C {
				// C2C는 from to 모두 캐시를 지워줘야함
				model.GetDB().CacheDelSwapWallet(toAddr)
				model.GetDB().CacheDelSwapWallet(swapInfo.SwapFromCoin.WalletAddress)
			}

			if err := model.GetDB().USPAU_Cmplt_Exchanges(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges err : %v", err)
			} else {
				sendExchangeLog(swapInfo)
			}
		}

		apiParams := &api_inno_log.AccountCoinLog{
			LogDt:         time.Now().Format("2006-01-02 15:04:05.000"),
			LogID:         int64(context.LogID_exchange),
			EventID:       int64(context.EventID_add),
			TxHash:        txHash,
			AUID:          swapInfo.AUID,
			CoinID:        swapInfo.SwapToCoin.CoinID,
			BaseCoinID:    swapInfo.SwapToCoin.BaseCoinID,
			WalletAddress: swapInfo.SwapToCoin.WalletAddress,
			WalletTypeID:  swapInfo.SwapToCoin.WalletTypeID,
			AdjQuantity:   "-" + strconv.FormatFloat(swapInfo.SwapToCoin.AdjustCoinQuantity, 'f', -1, 64),
			WalletID:      swapInfo.SwapToCoin.WalletID,
		}
		go api_inno_log.GetInstance().PostAccountCoins(apiParams)

	} else if txType.Target == context.From_user_to_parent_wallet { // 자식 지갑에서 부모 지갑으로 전송 : swap coin->point
	} else if txType.Target == context.From_user_to_other_wallet { // 자식지갑에서 다른 지갑으로 코인 전송
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
	}

	return resp
}

func TransferResultDepositWallet(fromAddr, toAddr, value, symbol, txHash string, gasFee int64, decimal int) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	// 부모지갑으로 입금된것만 처리한다.
	parentWallet, ok := config.GetInstance().ParentWalletsMap[toAddr]
	if !ok || !strings.EqualFold(toAddr, parentWallet.ParentWalletAddr) {
		return resp
	}

	// swap 수수료 입금 확인
	// 1. fromAddr redis에서 정보 추출
	// 1-1. point->coin, coin->point 분기
	// 2. txHash로 동일 한지 확인
	// 2-1. 메인넷 콜백이 더 빨라서 redis에 txHash가 아직 저장되지 않았다면 (TxStatus가 1(초기화) 상태) 경우에만 전송 value가 동일한지 check해서 처리

	// 1. fromAddr redis에서 정보 추출
	swapInfo, err := model.GetDB().CacheGetSwapWallet(fromAddr)
	if err != nil {
		log.Warnf("not exist fromAddr : %v, txHash:%v", fromAddr, txHash)
		return resp
	}

	// 1-1. point->coin, coin->point 분기
	if swapInfo.TxType == context.EventID_P2C {
		// 2. txHash로 동일 한지 확인
		if len(swapInfo.TxHash) != 0 && !strings.EqualFold(swapInfo.TxHash, txHash) {
			log.Warnf("not equal txhash => redis:%v, rev:%v, from:%v", swapInfo.TxHash, txHash, fromAddr)
			return resp
		} else if len(swapInfo.TxHash) == 0 {
			// 2-1. 메인넷 콜백이 더 빨라서 redis에 txHash가 아직 저장되지 않았다면 (TxStatus가 1(초기화) 상태) 경우에만 전송 value가 동일한지 check해서 처리
			if swapInfo.TxStatus == context.SWAP_status_init {
				fe := util.ToDecimalEncf(value, int64(decimal))

				basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
				basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
				swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

				if swapInfo.SwapFee == fe { // swap 수수료로 전송 받은 양이 동일하다면 정상 swap 수수료 수신으로 인식하고 정상 처리 해준다.
					// db에 수수료 전송 성공 저장
					swapInfo.IsFeeComplete = true
					if err := SwapFeeSuccess(swapInfo, txHash); err == nil {
						// swap 토큰 전송
						SwapTokenTransfer(swapInfo)
					}
				} else {
					log.Errorf("not equal swap fee redis:%v, rev:%v, txid:%v", swapInfo.SwapFee, fe, swapInfo.TxID)
				}
			} else {
				log.Errorf("invalid status txhash:%v, from:%v", txHash, fromAddr)
			}
		} else if strings.EqualFold(swapInfo.TxHash, txHash) {
			fe := util.ToDecimalEncf(value, int64(decimal))

			basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
			basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
			swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

			if swapInfo.SwapFee == fe {
				// db에 수수료 전송 성공 저장
				swapInfo.IsFeeComplete = true
				if err := SwapFeeSuccess(swapInfo, txHash); err == nil {
					// swap 토큰 전송
					SwapTokenTransfer(swapInfo)
				}
			} else {
				log.Errorf("not equal swap fee redis:%v, rev:%v, txid:%v", swapInfo.SwapFee, fe, swapInfo.TxID)
			}
		}

		// 수수료 입금 로그 전송
		apiParams := &api_inno_log.AccountCoinLog{
			LogDt:         time.Now().Format("2006-01-02 15:04:05.000"),
			LogID:         int64(context.LogID_exchange),
			EventID:       int64(context.EventID_add_fee),
			TxHash:        txHash,
			AUID:          swapInfo.AUID,
			CoinID:        swapInfo.SwapFeeCoinID,
			BaseCoinID:    model.GetDB().CoinsBySymbol[symbol].BaseCoinID,
			WalletAddress: swapInfo.SwapToCoin.WalletAddress,
			WalletTypeID:  swapInfo.SwapToCoin.WalletTypeID,
			AdjQuantity:   strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64),
			WalletID:      swapInfo.SwapToCoin.WalletID,
		}
		go api_inno_log.GetInstance().PostAccountCoins(apiParams)
	} else if swapInfo.TxType == context.EventID_C2P {
		// coin -> point 인 경우 토큰 입금 확인이 되면 포인트 DB 처리 해준다.
		fe := util.ToDecimalEncf(value, int64(decimal))

		basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
		basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
		swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

		if strings.EqualFold(swapInfo.SwapFromCoin.TokenTxHash, txHash) {
			// 시퀀스에 맞게 입금 콜백이 온경우
			swapInfo.TxStatus = context.SWAP_status_token_transfer_deposit_success
			if err := model.GetDB().USPAU_Mod_TransactExchanges_TxStatus(swapInfo.SwapFromCoin.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64), swapInfo); err == nil {
				if err := model.GetDB().USPAU_Cmplt_Exchanges(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
					log.Errorf("USPAU_Cmplt_Exchanges err : %v", err)
				} else {
					sendExchangeLog(swapInfo)
					model.GetDB().CacheDelSwapWallet(fromAddr)
				}
			}
		} else if len(swapInfo.SwapFromCoin.TokenTxHash) == 0 {
			if swapInfo.TxStatus != context.SWAP_status_init {
				log.Errorf("invalid status txhash:%v, from:%v", txHash, fromAddr)
				return resp
			}
			// 토큰 전송 유저 정보보다 콜백이 먼저 들어온경우 전송 량을 비교해서 같은지 판단한다.
			if math.Abs(swapInfo.SwapFromCoin.AdjustCoinQuantity) == fe {
				swapInfo.TxStatus = context.SWAP_status_token_transfer_deposit_success
				swapInfo.SwapFromCoin.TokenTxHash = txHash
				if err := model.GetDB().USPAU_Mod_TransactExchanges_Coin(swapInfo.TxID, swapInfo.TxStatus, txHash, time.Now().Format("2006-01-02 15:04:05.000"), swapInfo.SwapFromCoin.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
					log.Errorf("USPAU_Mod_TransactExchanges_Coin err : %v", err)
				} else {
					if err := model.GetDB().USPAU_Cmplt_Exchanges(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
						log.Errorf("USPAU_Cmplt_Exchanges err : %v", err)
					} else {
						sendExchangeLog(swapInfo)
						model.GetDB().CacheDelSwapWallet(fromAddr)
					}
				}
			}
		}

		// 역스왑 토큰 입금 로그 전송
		apiParams := &api_inno_log.AccountCoinLog{
			LogDt:         time.Now().Format("2006-01-02 15:04:05.000"),
			LogID:         int64(context.LogID_exchange),
			EventID:       int64(context.EventID_sub),
			TxHash:        txHash,
			AUID:          swapInfo.AUID,
			CoinID:        swapInfo.SwapFromCoin.CoinID,
			BaseCoinID:    swapInfo.SwapFromCoin.BaseCoinID,
			WalletAddress: swapInfo.SwapFromCoin.WalletAddress,
			WalletTypeID:  swapInfo.SwapFromCoin.WalletTypeID,
			AdjQuantity:   strconv.FormatFloat(swapInfo.SwapFromCoin.AdjustCoinQuantity, 'f', -1, 64),
			WalletID:      swapInfo.SwapFromCoin.WalletID,
		}
		go api_inno_log.GetInstance().PostAccountCoins(apiParams)

	} else if swapInfo.TxType == context.EventID_C2C {
		// c2c는 수수료를 먼저 입금 받은 후 유저의 토큰을 순차적으로 둘다 입금을 받고나서 swap할 토큰을 전송 해준다.
		// fromcoin의 txhash와 수수료 전송 txhash값이 모두 존재 할때 코인 전송을 해준다.

		// 1. 수수료 txhash가 유져를 통해 들어오거나 콜백이 먼저 들어오는 경우
		if strings.EqualFold(swapInfo.SwapFeeCoinSymbol, symbol) { // 수수료인지 아닌지 구별할 방법이 없어서 어짜피 swap중에는 다른 swap을 할수 없으니 symbol로 체크
			if len(swapInfo.TxHash) != 0 && !strings.EqualFold(swapInfo.TxHash, txHash) {
				// 수수료 관련 txhash가 아니다.
			} else if len(swapInfo.TxHash) == 0 {
				// 2-1. 메인넷 콜백이 더 빨라서 redis에 txHash가 아직 저장되지 않았다면 (TxStatus가 1(초기화) 상태) 경우에만 전송 value가 동일한지 check해서 처리
				if swapInfo.TxStatus == context.SWAP_status_init {
					fe := util.ToDecimalEncf(value, int64(decimal))

					basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
					basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
					swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

					if swapInfo.SwapFee == fe { // swap 수수료로 전송 받은 양이 동일하다면 정상 swap 수수료 수신으로 인식하고 정상 처리 해준다.
						swapInfo.IsFeeComplete = true
						// db에 수수료 전송 성공 저장
						if err := SwapFeeSuccess(swapInfo, txHash); err != nil {

						}
					} else {
						// 수수료가 아닌 다른 정보가 수신되었음 에러는 아님
						log.Warnf("not equal swap fee redis:%v, rev:%v, txid:%v", swapInfo.SwapFee, fe, swapInfo.TxID)
					}
				} else {
					log.Warnf("invalid status txhash:%v, from:%v", txHash, fromAddr)
				}
			} else if strings.EqualFold(swapInfo.TxHash, txHash) {
				// 유져에 의해 수수료 txhash 정보가 정상적으로 입력된 상태에서 콜백이 수신도어 있을때
				if swapInfo.TxStatus == context.SWAP_status_fee_transfer_start {
					fe := util.ToDecimalEncf(value, int64(decimal))

					basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
					basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
					swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

					if swapInfo.SwapFee == fe {
						// db에 수수료 전송 성공 저장
						swapInfo.IsFeeComplete = true
						if err := SwapFeeSuccess(swapInfo, txHash); err == nil {

						}
					} else {
						// 수수료가 아닌 다른 정보가 수신되었음 에러는 아님
						log.Warnf("not equal swap fee redis:%v, rev:%v, txid:%v", swapInfo.SwapFee, fe, swapInfo.TxID)
					}
				} else {
					log.Warnf("invalid status txhash:%v, from:%v", txHash, fromAddr)
				}
			}
		}

		// 수수로 입금 콜백이 아닌 유저 코인 전송 콜백인 경우 처리
		if strings.EqualFold(swapInfo.SwapFromCoin.CoinSymbol, symbol) {
			// 코인 입금인지 체크
			fe := util.ToDecimalEncf(value, int64(decimal))

			basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
			basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
			swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

			if strings.EqualFold(swapInfo.SwapFromCoin.TokenTxHash, txHash) {
				// 시퀀스에 맞게 입금 콜백이 온경우
				swapInfo.TxStatus = context.SWAP_status_token_transfer_deposit_success
				if err := model.GetDB().USPAU_Mod_TransactExchanges_TxStatus(swapInfo.SwapFromCoin.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64), swapInfo); err == nil {
					swapInfo.SwapFromCoin.IsComplete = true
					model.GetDB().CacheSetSwapWallet(swapInfo)
				}
			} else if len(swapInfo.SwapFromCoin.TokenTxHash) == 0 {
				// 토큰 전송 유저 정보보다 콜백이 먼저 들어온경우 전송 량을 비교해서 같은지 판단한다.
				if math.Abs(swapInfo.SwapFromCoin.AdjustCoinQuantity) == fe {
					swapInfo.TxStatus = context.SWAP_status_token_transfer_deposit_success
					swapInfo.SwapFromCoin.TokenTxHash = txHash
					if err := model.GetDB().USPAU_Mod_TransactExchanges_Coin(swapInfo.TxID, swapInfo.TxStatus, txHash, time.Now().Format("2006-01-02 15:04:05.000"), swapInfo.SwapFromCoin.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err == nil {
						swapInfo.SwapFromCoin.IsComplete = true
						model.GetDB().CacheSetSwapWallet(swapInfo)
					}
				}
			}
		}

		// 둘다 입금 처리 되었다면 코인 전송을 시작한다.
		if swapInfo.SwapFromCoin.IsComplete && swapInfo.IsFeeComplete {
			SwapTokenTransfer(swapInfo)
			// 캐시 삭제는 출금 완료 콜백 받는 부분에서 처리한다.
		}

	} else {
		log.Errorf("not exist txType : %v, from:%v, txHash:%v", swapInfo.TxType, fromAddr, txHash)
		return resp
	}

	return resp
}

func SwapFeeSuccess(swapInfo *context.ReqSwapInfo, txHash string) error {
	if len(swapInfo.TxHash) == 0 && swapInfo.TxStatus == context.SWAP_status_init {
		// 수수료 시작 전에 먼저 수수료 전송 완료가 들어온경우
		swapInfo.TxHash = txHash
		swapInfo.TxStatus = context.SWAP_status_fee_transfer_success

		if err := model.GetDB().USPAU_Mod_TransactExchanges_ExchangeFees(swapInfo.TxID, swapInfo.TxStatus, swapInfo.TxHash, strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64), model.GetDB().Coins[swapInfo.SwapFeeCoinID].BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
			return err
		}
		if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
			log.Errorf("CacheSetSwapWallet err :%v", err)
			return err
		}
	} else if swapInfo.TxStatus == context.SWAP_status_fee_transfer_start {
		// 수수료 전송 시작 상태에서 완료 콜백이 들어온경우
		swapInfo.TxStatus = context.SWAP_status_fee_transfer_success

		if err := model.GetDB().USPAU_Mod_TransactExchanges_TxStatus(model.GetDB().Coins[swapInfo.SwapFeeCoinID].BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64), swapInfo); err != nil {
			return err
		}
		if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
			log.Errorf("CacheSetSwapWallet err :%v", err)
			return err
		}
	} else {
		log.Warnf("invalid case swap_status:%v", swapInfo.TxStatus)
	}

	return nil
}

func SwapTokenTransfer(swapInfo *context.ReqSwapInfo) *base.BaseResponse {
	// swap 수수료 정상 입금 처리
	swapInfo.TxStatus = context.SWAP_status_fee_transfer_success

	swapCoin := &context.SwapCoin{}
	if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_C2C {
		swapCoin = &swapInfo.SwapToCoin
	} else {
		return nil
	}

	// swap 토큰 전송 처리 시작
	reqFromParent := &context.ReqCoinTransferFromParentWallet{
		AUID:       swapInfo.AUID,
		CoinID:     swapCoin.CoinID,
		CoinSymbol: swapCoin.CoinSymbol,
		ToAddress:  swapCoin.WalletAddress,
		Quantity:   swapCoin.AdjustCoinQuantity,
	}

	resp := TransferFromParentWallet(reqFromParent, false)
	if resp.Return != 0 {
		log.Errorf("Rceived a fee but failed to send coins => return:%v message:%v", resp.Return, resp.Message)
		return resp
	}

	swapInfo.TxStatus = context.SWAP_status_token_transfer_withdrawal_start
	swapCoin.TokenTxHash = reqFromParent.TransactionId

	// swap 토큰 전송 시작 기록
	if err := model.GetDB().USPAU_Mod_TransactExchanges_Coin(swapInfo.TxID, swapInfo.TxStatus, swapCoin.TokenTxHash, time.Now().Format("2006-01-02 15:04:05.000"), 0, ""); err == nil {
		if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
			log.Errorf("CacheSetSwapWallet err :%v", err)
			return resp
		}
	}

	return resp
}

func sendExchangeLog(swapInfo *context.ReqSwapInfo) {
	isC2pComplete := func() bool {
		if swapInfo.TxType == context.EventID_C2P && swapInfo.TxStatus == context.SWAP_status_token_transfer_deposit_success {
			return true
		}
		return false
	}()
	isC2CP2CComplete := func() bool {
		if (swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_C2C) && swapInfo.TxStatus == context.SWAP_status_token_transfer_withdrawal_success {
			return true
		}
		return false
	}()
	isP2PComplete := func() bool {
		return swapInfo.TxType == context.EventID_P2P
	}()

	// swap이. 종료된 경우에만 최종 로그를 남긴다.
	if isC2pComplete || isC2CP2CComplete || isP2PComplete {
		apiParams := &api_inno_log.ExchangeLogs{
			LogDT:   time.Now().Format("2006-01-02 15:04:05.000"),
			EventID: swapInfo.TxType,
			TxID:    swapInfo.TxID,
			AUID:    swapInfo.AUID,
			MUID: func() int64 {
				if swapInfo.TxType == context.EventID_P2C {
					return swapInfo.SwapFromPoint.MUID
				} else if swapInfo.TxType == context.EventID_C2P {
					return swapInfo.SwapToPoint.MUID
				} else {
					return 0 // C2C, P2P
				}
			}(),
			InnoUID: swapInfo.InnoUID,
			AppID: func() int64 {
				if swapInfo.TxType == context.EventID_P2C {
					return swapInfo.SwapFromPoint.AppID
				} else if swapInfo.TxType == context.EventID_C2P {
					return swapInfo.SwapToPoint.AppID
				} else {
					return 0 // C2C, P2P
				}
			}(),
			ExchangeFees: func() string {
				if len(swapInfo.SwapFeeD) == 0 {
					return "0"
				}
				return swapInfo.SwapFeeD
			}(),
			FromBaseCoinID:    swapInfo.SwapFromCoin.BaseCoinID,
			FromWalletTypeID:  swapInfo.SwapFromCoin.WalletTypeID,
			FromWalletID:      swapInfo.SwapFromCoin.WalletID,
			FromWalletAddress: swapInfo.SwapFromCoin.WalletAddress,
			FromID: func() int64 {
				if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_P2P {
					return swapInfo.SwapFromPoint.PointID
				} else if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_C2C {
					return swapInfo.SwapFromCoin.CoinID
				} else {
					return 0 // C2C
				}
			}(),
			FromAdjQuantity: func() string {
				if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_P2P {
					return strconv.FormatInt(swapInfo.SwapFromPoint.AdjustPointQuantity, 10)
				} else if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_C2C {
					return strconv.FormatFloat(swapInfo.SwapFromCoin.AdjustCoinQuantity, 'f', -1, 64)
				}
				return ""
			}(),
			ToBaseCoinID:    swapInfo.SwapToCoin.BaseCoinID,
			ToWalletTypeID:  swapInfo.SwapToCoin.WalletTypeID,
			ToWalletID:      swapInfo.SwapToCoin.WalletID,
			ToWalletAddress: swapInfo.SwapToCoin.WalletAddress,
			ToID: func() int64 {
				if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_P2P {
					return swapInfo.SwapToPoint.PointID
				} else if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_C2C {
					return swapInfo.SwapToCoin.CoinID
				}
				return 0
			}(),
			ToAdjQuantity: func() string {
				if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_P2P {
					return strconv.FormatInt(swapInfo.SwapToPoint.AdjustPointQuantity, 10)
				} else if swapInfo.TxType == context.EventID_P2C || swapInfo.TxType == context.EventID_C2C {
					return strconv.FormatFloat(swapInfo.SwapToCoin.AdjustCoinQuantity, 'f', -1, 64)
				}
				return ""
			}(),
			TransactedDT: time.Unix(swapInfo.CreateAt, 0).Format("2006-01-02 15:04:05.000"),
			CompletedDT:  time.Now().Format("2006-01-02 15:04:05.000"),
		}
		log.Debugf("log : %v", apiParams)
		go api_inno_log.GetInstance().PostExchange(apiParams)
	}
}
