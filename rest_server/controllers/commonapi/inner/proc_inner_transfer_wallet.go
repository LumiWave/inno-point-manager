package inner

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/api_inno_log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/util"
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

			swapInfo.TxStatus = context.SWAP_status_token_transfer_fail
			if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_TxStatus(swapInfo.BaseCoinID, fe, swapInfo); err == nil {
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

		swapInfo.TxStatus = context.SWAP_status_token_transfer_success
		if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_TxStatus(swapInfo.BaseCoinID, fe, swapInfo); err == nil {
			model.GetDB().CacheDelSwapWallet(toAddr)
			if err := model.GetDB().USPAU_Cmplt_ExchangeGoods(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
				log.Errorf("USPAU_Cmplt_ExchangeGoods err : %v", err)
			}
		}

		apiParams := &api_inno_log.AccountCoinLog{
			LogDt:         time.Now().Format("2006-01-02 15:04:05.000"),
			LogID:         int64(context.LogID_exchange),
			EventID:       int64(context.EventID_sub),
			TxHash:        txHash,
			AUID:          swapInfo.AUID,
			CoinID:        swapInfo.CoinID,
			BaseCoinID:    swapInfo.BaseCoinID,
			WalletAddress: swapInfo.ToWalletAddress,
			AdjQuantity:   "-" + strconv.FormatFloat(swapInfo.SwapCoin.AdjustCoinQuantity, 'f', -1, 64),
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
	if swapInfo.TxType == context.EventID_toCoin {
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
			WalletAddress: swapInfo.ToWalletAddress,
			AdjQuantity:   strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64),
		}
		go api_inno_log.GetInstance().PostAccountCoins(apiParams)
	} else if swapInfo.TxType == context.EventID_toPoint {
		// coin -> point 인 경우 토큰 입금 확인이 되면 포인트 DB 처리 해준다.
		fe := util.ToDecimalEncf(value, int64(decimal))

		basecoinID := model.GetDB().CoinsBySymbol[symbol].BaseCoinID
		basecoinSymbol := model.GetDB().BaseCoinMapByCoinID[basecoinID].BaseCoinSymbol
		swapInfo.TxGasFee = util.ToDecimalEncf(strconv.FormatInt(gasFee, 10), model.GetDB().CoinsBySymbol[basecoinSymbol].Decimal)

		if strings.EqualFold(swapInfo.SwapCoin.TokenTxHash, txHash) {
			// 시퀀스에 맞게 입금 콜백이 온경우
			swapInfo.TxStatus = context.SWAP_status_token_transfer_success
			if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_TxStatus(swapInfo.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64), swapInfo); err == nil {
				if err := model.GetDB().USPAU_Cmplt_ExchangeGoods(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
					log.Errorf("USPAU_Cmplt_ExchangeGoods err : %v", err)
				} else {
					model.GetDB().CacheDelSwapWallet(fromAddr)
				}
			}
		} else if len(swapInfo.SwapCoin.TokenTxHash) == 0 {
			if swapInfo.TxStatus != context.SWAP_status_init {
				log.Errorf("invalid status txhash:%v, from:%v", txHash, fromAddr)
				return resp
			}
			// 토큰 전송 유저 정보보다 콜백이 먼저 들어온경우 전송 량을 비교해서 같은지 판단한다.
			if math.Abs(swapInfo.SwapCoin.AdjustCoinQuantity) == fe {
				swapInfo.TxStatus = context.SWAP_status_token_transfer_success
				swapInfo.SwapCoin.TokenTxHash = txHash
				if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_Coin(swapInfo.TxID, swapInfo.TxStatus, txHash, time.Now().Format("2006-01-02 15:04:05.000"), swapInfo.BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
					log.Errorf("USPAU_Mod_TransactExchangeGoods_Coin err : %v", err)
				} else {
					if err := model.GetDB().USPAU_Cmplt_ExchangeGoods(swapInfo, time.Now().Format("2006-01-02 15:04:05.000"), true); err != nil {
						log.Errorf("USPAU_Cmplt_ExchangeGoods err : %v", err)
					} else {
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
			CoinID:        swapInfo.CoinID,
			BaseCoinID:    swapInfo.BaseCoinID,
			WalletAddress: swapInfo.ToWalletAddress,
			AdjQuantity:   strconv.FormatFloat(swapInfo.SwapCoin.AdjustCoinQuantity, 'f', -1, 64),
		}
		go api_inno_log.GetInstance().PostAccountCoins(apiParams)

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

		if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_Exchangefee(swapInfo.TxID, swapInfo.TxStatus, swapInfo.TxHash, strconv.FormatFloat(swapInfo.SwapFee, 'f', -1, 64), model.GetDB().Coins[swapInfo.SwapFeeCoinID].BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64)); err != nil {
			return err
		}
		if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
			log.Errorf("CacheSetSwapWallet err :%v", err)
			return err
		}
	} else if swapInfo.TxStatus == context.SWAP_status_fee_transfer_start {
		// 수수료 전송 시작 상태에서 완료 콜백이 들어온경우
		swapInfo.TxStatus = context.SWAP_status_fee_transfer_success

		if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_TxStatus(model.GetDB().Coins[swapInfo.SwapFeeCoinID].BaseCoinID, strconv.FormatFloat(swapInfo.TxGasFee, 'f', -1, 64), swapInfo); err != nil {
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
	// swap 토큰 전송 처리 시작
	reqFromParent := &context.ReqCoinTransferFromParentWallet{
		AUID:       swapInfo.AUID,
		CoinID:     swapInfo.CoinID,
		CoinSymbol: swapInfo.CoinSymbol,
		ToAddress:  swapInfo.WalletAddress,
		Quantity:   swapInfo.AdjustCoinQuantity,
	}

	resp := TransferFromParentWallet(reqFromParent, false)
	if resp.Return != 0 {
		log.Errorf("Rceived a fee but failed to send coins => return:%v message:%v", resp.Return, resp.Message)
		return resp
	}

	swapInfo.TxStatus = context.SWAP_status_token_transfer_start
	swapInfo.TokenTxHash = reqFromParent.TransactionId

	// swap 토큰 전송 시작 기록
	if err := model.GetDB().USPAU_Mod_TransactExchangeGoods_Coin(swapInfo.TxID, swapInfo.TxStatus, swapInfo.TokenTxHash, time.Now().Format("2006-01-02 15:04:05.000"), 0, ""); err == nil {
		if err := model.GetDB().CacheSetSwapWallet(swapInfo); err != nil {
			log.Errorf("CacheSetSwapWallet err :%v", err)
			return resp
		}
	}

	return resp
}
