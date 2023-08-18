package inner

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
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

		// 아래 두경우 swap 복구 루틴이 똑같다.
		if txType.Target == context.From_user_to_fee_wallet || txType.Target == context.From_parent_to_other_wallet {
			// swap 정보 복구 하고 redis 삭제
			swapKey := model.MakeSwapKey(fromAddr)
			if reqSwapInfo, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
				log.Errorf("GetCacheSwapInfo err => auid:%v txid:%v", err, txType.AUID, txHash)
				return resp
			} else {
				// 최신 코인 수량 정보 수집 해서 swap 시도 하려 했던 만큼 더해준다.
				// _, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
				// if err != nil {
				// 	log.Errorf("GetAccountCoins error : %v, auid:%v txid:%v", err, txType.AUID, txHash)
				// 	model.MakeDbError(resp, resultcode.Result_DBError, err)
				// 	return resp
				// }
				//meCoin := coinsMap[reqSwapInfo.CoinID]
				// reqSwapInfo.PreviousCoinQuantity = meCoin.Quantity
				// reqSwapInfo.CoinQuantity = meCoin.Quantity - reqSwapInfo.AdjustCoinQuantity
				reqSwapInfo.AdjustCoinQuantity = -(reqSwapInfo.AdjustCoinQuantity)
				// 최신 포인트 정보 수집
				nowPointQuantity, err := model.GetDB().GetPointApp(reqSwapInfo.MUID, reqSwapInfo.AppID, reqSwapInfo.DatabaseID)
				if err != nil {
					log.Errorf("GetPointApp error : %v, auid:%v muid:%v txid:%v", err, txType.AUID, reqSwapInfo.MUID, txHash)
					model.MakeDbError(resp, resultcode.Result_DBError, err)
					return resp
				}
				reqSwapInfo.PreviousPointQuantity = nowPointQuantity
				reqSwapInfo.PointQuantity = nowPointQuantity - reqSwapInfo.AdjustPointQuantity
				reqSwapInfo.AdjustPointQuantity = -(reqSwapInfo.AdjustPointQuantity)
				// event id 역으로 바꿔준다
				if reqSwapInfo.TxType == context.EventID_toCoin {
					reqSwapInfo.TxType = context.EventID_toPoint
				} else if reqSwapInfo.TxType == context.EventID_toPoint {
					reqSwapInfo.TxType = context.EventID_toCoin
				}

				// if err := model.GetDB().PostPointCoinSwap(reqSwapInfo, txHash); err != nil {
				// 	resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
				// 	return resp
				// }

				model.GetDB().DelCacheSwapInfo(swapKey)
			}
		}

		if txType.Target == context.From_user_to_fee_wallet { // point -> coin swap 요청을 위한 1단계 수수료 출금 정보
			model.GetDB().DelCacheCoinTransfer(tKey) // txid redis 삭제
			userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
			model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user redis 삭제
		} else if txType.Target == context.From_parent_to_other_wallet { // coin -> point swap 요청 : 부모에서 자식에서 코인 입금
			model.GetDB().DelCacheCoinTransfer(tKey) // txid redis 삭제
			parentKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
			model.GetDB().DelCacheCoinTransferFromParentWallet(parentKey) // from parent redis 삭제
		}

		return resp
	}

	// 성공 분기처리
	if txType.Target == context.From_user_to_fee_wallet { // point -> coin swap 요청을 위한 1단계 수수료 출금 정보
		// 자식 지갑에서 수수료 + 전송량 빼기
		// _, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		// if err != nil {
		// 	log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, txHash)
		// 	model.MakeDbError(resp, resultcode.Result_DBError, err)
		// 	return resp
		// }
		// meCoin, ok := coinsMap[txType.CoinID]
		// if !ok {
		// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
		// 	return resp
		// }

		// scale := new(big.Float).SetFloat64(1)
		// scale.SetString("1e" + fmt.Sprintf("%d", decimal))

		// valueAmount, _ := new(big.Float).SetString(value)
		// valueAmount = new(big.Float).Quo(valueAmount, scale)
		// amount, _ := valueAmount.Float64()

		// feeAmount, _ := new(big.Float).SetString(fee)
		// feeAmount = new(big.Float).Quo(feeAmount, scale)
		// fe, _ := feeAmount.Float64()

		// adjust := new(big.Float).Add(valueAmount, feeAmount)
		// adjustQuantity, _ := adjust.Float64()

		// lastQuantity := new(big.Float).SetFloat64(meCoin.Quantity)
		// new := new(big.Float).Sub(lastQuantity, adjust)
		// newQuantity, _ := new.Float64()

		// if err := model.GetDB().UpdateAccountCoins(
		// 	txType.AUID,
		// 	txType.CoinID,
		// 	model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		// 	meCoin.WalletAddress,
		// 	meCoin.Quantity,
		// 	-adjustQuantity, // 전송 수수료 + amount
		// 	newQuantity,
		// 	context.LogID_external_wallet,
		// 	context.EventID_sub,
		// 	txHash); err != nil {
		// 	log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v", err, resp.Return, resp.Message, txHash)
		// 	return resp
		// }
		// meCoin.Quantity = meCoin.Quantity - (amount + fe) // 남은 수량

		// redis swap 정보 찾아서 부모지갑에서 자식지갑으로 코인 전송
		swapKey := model.MakeSwapKey(fromAddr)
		if reqSwapInfo, err := model.GetDB().GetCacheSwapInfo(swapKey); err == nil {
			reqFromParent := &context.ReqCoinTransferFromParentWallet{
				AUID:       reqSwapInfo.AUID,
				CoinID:     reqSwapInfo.CoinID,
				CoinSymbol: reqSwapInfo.CoinSymbol,
				ToAddress:  reqSwapInfo.WalletAddress,
				Quantity:   reqSwapInfo.AdjustCoinQuantity,
			}

			resp = TransferFromParentWallet(reqFromParent, false)
			if resp.Return != 0 {
				log.Errorf("Rceived a fee but failed to send coins => return:%v message:%v", resp.Return, resp.Message)
				return resp
			}
		}
		// user가 보낸 redis 정보는 삭제 해준다. 이때는 이미 from parent reids 정보가 생성되어 있음
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
	} else if txType.Target == context.From_parent_to_other_wallet { // 부모지갑에서 자식지갑으로 코인 전송 : swap point->coin, 혹은 부모 지갑에서만 출금 가능한 코인(MATIC) 외부 출금
		// 부모지갑에서 자식지갑으로 입금시 swap으로 간주하고 swap 처리 진행한다.
		swapKey := model.MakeSwapKey(toAddr)
		parentKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
		_, err1 := model.GetDB().GetCacheSwapInfo(swapKey)
		_, err2 := model.GetDB().GetCacheCoinTransferFromParentWallet(parentKey)
		if err1 != nil && err2 != nil {
			log.Errorf("GetCacheSwapInfo err1 => %v auid:%v txid:%v", err1, txType.AUID, txHash)
			log.Errorf("GetCacheCoinTransferFromParentWallet err2 => %v auid:%v txid:%v", err2, txType.AUID, txHash)
			return resp
		}

		if err2 == nil {
			// 부모지갑사용 코인이 아니면 패스
			// coin := model.GetDB().Coins[txType.CoinID]
			// if baseCoin, ok := model.GetDB().BaseCoinMapByCoinID[coin.BaseCoinID]; ok {
			// 	if baseCoin.IsUsedParentWallet {
			// 		// 부모지갑출금을 사용하는 코인 외부 지갑 전송 성공으로 처리하고 그 보낸 사람의 코인을 차감 시킨다.
			// 		// 자식 지갑의 코인 빼고 수수료 코인은 빼지 말고 일단 주석처리한다.
			// 		_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
			// 		if err != nil {
			// 			log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, txHash)
			// 			model.MakeDbError(resp, resultcode.Result_DBError, err)
			// 			return resp
			// 		}

			// 		meCoin, ok := coinsMap[txType.CoinID]
			// 		if !ok {
			// 			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
			// 			return resp
			// 		}

			// 		scale := new(big.Float).SetFloat64(1)
			// 		scale.SetString("1e" + fmt.Sprintf("%d", decimal))

			// 		valueAmount, _ := new(big.Float).SetString(value)
			// 		valueAmount = new(big.Float).Quo(valueAmount, scale)
			// 		amount, _ := valueAmount.Float64()

			// 		if err := model.GetDB().UpdateAccountCoins(
			// 			txType.AUID,
			// 			txType.CoinID,
			// 			model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
			// 			meCoin.WalletAddress,
			// 			meCoin.Quantity,
			// 			-amount, // amount
			// 			meCoin.Quantity-amount,
			// 			context.LogID_external_wallet,
			// 			context.EventID_sub,
			// 			txHash); err != nil {
			// 			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v amount:%v", err, resp.Return, resp.Message, txHash, -amount)
			// 			return resp
			// 		}
			// 	}
			// }

		}
		// parent가 보낸 redis 정보는 삭제 해준다.
		model.GetDB().DelCacheCoinTransfer(tKey)                      //tx 삭제
		model.GetDB().DelCacheCoinTransferFromParentWallet(parentKey) // from user 삭제
		model.GetDB().DelCacheSwapInfo(swapKey)
	} else if txType.Target == context.From_user_to_parent_wallet { // 자식 지갑에서 부모 지갑으로 전송 : swap coin->point
		// 자식 지갑의 basecoin 수량 acutual fee 만큼 축소
		// _, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		// if err != nil {
		// 	log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, txHash)
		// 	model.MakeDbError(resp, resultcode.Result_DBError, err)
		// 	return resp
		// }

		// meCoin, ok := coinsMap[txType.CoinID]
		// if !ok {
		// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
		// 	return resp
		// }

		// baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID]
		// meBaseCoinID := int64(0)
		// for _, coin := range model.GetDB().Coins {
		// 	if strings.EqualFold(coin.CoinSymbol, baseCoinInfo.BaseCoinSymbol) {
		// 		meBaseCoinID = coin.CoinId
		// 	}
		// }

		// meBaseCoin, ok := coinsMap[meBaseCoinID]
		// if !ok {
		// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
		// 	return resp
		// }

		// scale := new(big.Float).SetFloat64(1)
		// scale.SetString("1e" + fmt.Sprintf("%d", decimal))

		// feeAmount, _ := new(big.Float).SetString(fee)
		// feeAmount = new(big.Float).Quo(feeAmount, scale)
		// fe, _ := feeAmount.Float64()

		// if err := model.GetDB().UpdateAccountCoins(
		// 	txType.AUID,
		// 	meBaseCoinID,
		// 	model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		// 	meBaseCoin.WalletAddress,
		// 	meBaseCoin.Quantity,
		// 	-fe, // amount
		// 	meBaseCoin.Quantity-fe,
		// 	context.LogID_external_wallet,
		// 	context.EventID_sub,
		// 	txHash); err != nil {
		// 	log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v", err, resp.Return, resp.Message, txHash)
		// 	return resp
		// }
		// meBaseCoin.Quantity = meCoin.Quantity - fe // 남은 수량

		// redis swap 정보 찾아서 스왑 처리
		swapKey := model.MakeSwapKey(fromAddr)
		if _, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
			log.Errorf("GetCacheSwapInfo err => audi:%v txid:%v", err, txType.AUID, txHash)
			return resp
		}
		// user가 보낸 redis 정보는 삭제 해준다. 이때는 이미 from parent reids 정보가 생성되어 있음
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
		model.GetDB().DelCacheSwapInfo(swapKey)
	} else if txType.Target == context.From_user_to_other_wallet { // 자식지갑에서 다른 지갑으로 코인 전송
		// 자식 지갑의 코인 빼고, 수수료 코인 뺀다.
		// _, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		// if err != nil {
		// 	log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, txHash)
		// 	model.MakeDbError(resp, resultcode.Result_DBError, err)
		// 	return resp
		// }

		// meCoin, ok := coinsMap[txType.CoinID]
		// if !ok {
		// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
		// 	return resp
		// }

		// baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID]
		// meBaseCoinID := int64(0)
		// for _, coin := range model.GetDB().Coins {
		// 	if strings.EqualFold(coin.CoinSymbol, baseCoinInfo.BaseCoinSymbol) {
		// 		meBaseCoinID = coin.CoinId
		// 	}
		// }

		// meBaseCoin, ok := coinsMap[meBaseCoinID]
		// if !ok {
		// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, txHash)
		// 	return resp
		// }

		// scale := new(big.Float).SetFloat64(1)
		// scale.SetString("1e" + fmt.Sprintf("%d", decimal))

		// valueAmount, _ := new(big.Float).SetString(value)
		// valueAmount = new(big.Float).Quo(valueAmount, scale)
		// amount, _ := valueAmount.Float64()

		// feeAmount, _ := new(big.Float).SetString(fee)
		// feeAmount = new(big.Float).Quo(feeAmount, scale)
		// fe, _ := feeAmount.Float64()

		// if err := model.GetDB().UpdateAccountCoins(
		// 	txType.AUID,
		// 	txType.CoinID,
		// 	model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		// 	meCoin.WalletAddress,
		// 	meCoin.Quantity,
		// 	-amount, // amount
		// 	meCoin.Quantity-amount,
		// 	context.LogID_external_wallet,
		// 	context.EventID_sub,
		// 	txHash); err != nil {
		// 	log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v amount:%v", err, resp.Return, resp.Message, txHash, -amount)
		// 	return resp
		// }
		// meCoin.Quantity = meCoin.Quantity - amount // 남은 수량

		// if err := model.GetDB().UpdateAccountCoins(
		// 	txType.AUID,
		// 	meBaseCoinID,
		// 	model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		// 	meBaseCoin.WalletAddress,
		// 	meBaseCoin.Quantity,
		// 	-fe, // amount
		// 	meBaseCoin.Quantity-fe,
		// 	context.LogID_external_wallet,
		// 	context.EventID_sub,
		// 	txHash); err != nil {
		// 	log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v fee:%v", err, resp.Return, resp.Message, txHash, fe)
		// 	return resp
		// }
		// meBaseCoin.Quantity = meCoin.Quantity - fe // 남은 수량

		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
	}

	return resp
}

func TransferResultDepositWallet(fromAddr, toAddr, value, symbol, txHash string, decimal int) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	// 외부 지갑에서 입금된것은 무조건 입금 처리 해준다.
	// swap통한 부모 입금시 출금쪽에서 누적하기 때문에 넘긴다.
	parentWallet, ok := config.GetInstance().ParentWalletsMap[fromAddr]
	if ok && strings.EqualFold(fromAddr, parentWallet.ParentWalletAddr) {
		return resp
	}

	// 입금 주소로 db 검색해서 AUID추출
	meCoin, err := model.GetDB().GetAccountCoinsByWalletAddress(toAddr, symbol)
	if err != nil {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			fromAddr, toAddr, symbol, value)
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	Lockkey := model.MakeCoinTransferFromUserWalletLockKey(meCoin.AUID)
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

	// lock이 풀렸다면 코인 수량이 변화하였을수도 있어서 다시 한번 더 불러온다.
	meCoin, err = model.GetDB().GetAccountCoinsByWalletAddress(toAddr, symbol)
	if err != nil {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			fromAddr, toAddr, symbol, value)
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	if meCoin.CoinID == 0 {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			fromAddr, toAddr, symbol, value)
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	// USPAU_Mod_AccountCoins 호출 하여 코인량 갱신
	valueAmount, _ := new(big.Float).SetString(value)
	scale := new(big.Float).SetFloat64(1)
	scale.SetString("1e" + fmt.Sprintf("%d", decimal))
	valueAmount = new(big.Float).Quo(valueAmount, scale)
	adjustQuantity, _ := valueAmount.Float64()

	if err := model.GetDB().UpdateAccountCoins(
		meCoin.AUID,
		meCoin.CoinID,
		model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		toAddr,
		meCoin.Quantity,
		adjustQuantity,
		meCoin.Quantity+adjustQuantity,
		context.LogID_external_wallet,
		context.EventID_add,
		txHash); err != nil {
		log.Errorf("UpdateAccountCoins error : %v", err)
	}
	return resp
}
