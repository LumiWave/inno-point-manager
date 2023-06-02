package inner

import (
	"strconv"
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/token_manager_server"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func TransferFromParentWallet(params *context.ReqCoinTransferFromParentWallet, isLockCheck bool) *base.BaseResponse {

	resp := new(base.BaseResponse)
	resp.Success()

	//1. redis에 외부 전송이 진행 중인지 체크 redis에 정보가 있다면 전송중으로 인지하면됨
	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	//3. redis에 전송 정보 저장
	//4. 콜백(internal api)으로 완료or실패 확인 후 db 프로지저 호출, redis 삭제

	// 0. redis lock
	if isLockCheck {
		Lockkey := model.MakeCoinTransferFromParentWalletLockKey(params.AUID)
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

		// 1. redis에 외부 전송 정보 존재하는지 check
		key := model.MakeCoinTransferFromParentWalletKey(params.AUID)
		_, err := model.GetDB().GetCacheCoinTransferFromParentWallet(key)
		if err == nil {
			// 전송중인 기존 정보가 있다면 에러를 리턴한다.
			log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_Inprogress])
			resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
			return resp
		}
	}

	coinInfo := model.GetDB().CoinsBySymbol[params.CoinSymbol]
	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	req := &token_manager_server.ReqSendFromParentWallet{
		BaseSymbol: model.GetDB().BaseCoinMapByCoinID[coinInfo.BaseCoinID].BaseCoinSymbol,
		Symbol:     params.CoinSymbol,
		ToAddress:  params.ToAddress,
		Amount:     strconv.FormatFloat(params.Quantity, 'f', -1, 64),
		Memo:       strconv.FormatInt(params.AUID, 10),
	}
	if res, err := token_manager_server.GetInstance().PostSendFromParentWallet(req); err != nil {
		resp.SetReturn(resultcode.ResultInternalServerError)
		return resp
	} else {
		if res.Return != 0 { // token manager 전송 에러
			resp.Return = res.Return
			resp.Message = res.Message
			return resp
		}

		if !res.Value.IsSuccess {
			resp.SetReturn(resultcode.ResultInternalServerError)
		}

		//params.ReqId = res.Value.ReqId
		//params.TransactionId = res.Value.TransactionId
		params.TransactionId = res.Value.TxHash
	}

	params.ActionDate = time.Unix(time.Now().Unix(), 0)

	//3. redis에 전송 정보 저장
	key := model.MakeCoinTransferFromParentWalletKey(params.AUID)
	if err := model.GetDB().SetCacheCoinTransferFromParentWallet(key, params); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetTransfer])
		resp.SetReturn(resultcode.Result_RedisError_SetTransfer)
		return resp
	}

	//4. redis에 전송 정보 transaction id key로 다시 한번더 저장 : 추후 콜백 api를 통해 검증하기 위해서
	tKey := model.MakeCoinTransferKeyByTxID(params.TransactionId)
	txType := &context.TxType{
		Target: context.From_parent_to_other_wallet,
		AUID:   params.AUID,
		CoinID: params.CoinID,
	}
	if err := model.GetDB().SetCacheCoinTransferTx(tKey, txType); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetTransfer_Tx])
		resp.SetReturn(resultcode.Result_RedisError_SetTransfer_Tx)
		return resp
	}

	resp.Value = params

	return resp
}

func TransferFromUserWallet(params *context.ReqCoinTransferFromUserWallet, isLockCheck bool) *base.BaseResponse {

	resp := new(base.BaseResponse)
	resp.Success()

	//1. redis에 외부 전송이 진행 중인지 체크 redis에 정보가 있다면 전송중으로 인지하면됨
	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	//3. redis에 전송 정보 저장
	//4. 콜백(internal api)으로 완료or실패 확인 후 db 프로지저 호출, redis 삭제

	// 0. redis lock
	if isLockCheck {
		Lockkey := model.MakeCoinTransferFromUserWalletLockKey(params.AUID)
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

		// 1. redis에 외부 전송 정보 존재하는지 check
		key := model.MakeCoinTransferFromUserWalletKey(params.AUID)
		_, err := model.GetDB().GetCacheCoinTransferFromUserWallet(key)
		if err == nil {
			// 전송중인 기존 정보가 있다면 에러를 리턴한다.
			log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_Inprogress])
			resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
			return resp
		}
	}

	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	req := &token_manager_server.ReqSendFromUserWallet{
		BaseCoinSymbol: params.BaseCoinSymbol,
		Symbol:         params.CoinSymbol,
		FromAddress:    params.FromAddress,
		ToAddress:      params.ToAddress,
		Amount:         strconv.FormatFloat(params.Quantity, 'f', -1, 64),
		Memo:           strconv.FormatInt(params.AUID, 10),
	}
	if res, err := token_manager_server.GetInstance().PostSendFromUserWallet(req); err != nil {
		resp.SetReturn(resultcode.ResultInternalServerError)
		return resp
	} else {
		if res.Return != 0 { // token manager 전송 에러
			resp.Return = res.Return
			if errMsg, ok := token_manager_server.ResultCodeText[res.Message]; ok {
				resp.Message = errMsg
			} else {
				resp.Message = res.Message
			}
			return resp
		}

		// if len(res.Value.TransactionHash) == 0 {
		// 	log.Errorf("PostSendFromUserWallet txid null")
		// }
		//params.ReqId = res.Value.ReqId
		params.TransactionId = res.Value.TxHash
	}

	params.ActionDate = time.Unix(time.Now().Unix(), 0)

	//3. redis에 전송 정보 저장
	key := model.MakeCoinTransferFromUserWalletKey(params.AUID)
	if err := model.GetDB().SetCacheCoinTransferFromUserWallet(key, params); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetTransfer])
		resp.SetReturn(resultcode.Result_RedisError_SetTransfer)
		return resp
	}

	//4. redis에 전송 정보 transaction id key로 다시 한번더 저장 : 추후 콜백 api를 통해 검증하기 위해서
	tKey := model.MakeCoinTransferKeyByTxID(params.TransactionId)
	txType := &context.TxType{
		Target: params.Target,
		AUID:   params.AUID,
		CoinID: params.CoinID,
	}
	if err := model.GetDB().SetCacheCoinTransferTx(tKey, txType); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetTransfer_Tx])
		resp.SetReturn(resultcode.Result_RedisError_SetTransfer_Tx)
		return resp
	}

	resp.Value = params

	return resp
}

func TransferResultDeposit(params *context.ReqCoinTransferResDeposit) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	// 외부 지갑에서 입금된것은 무조건 입금 처리 해준다.
	// swap통한 부모 입금시 출금쪽에서 누적하기 때문에 넘긴다.
	parentWallet, ok := config.GetInstance().ParentWalletsMap[params.FromAddress]
	if ok && strings.EqualFold(params.FromAddress, parentWallet.ParentWalletAddr) {
		return resp
	}

	// 입금 주소로 db 검색해서 AUID추출
	meCoin, err := model.GetDB().GetAccountCoinsByWalletAddress(params.ToAddress, params.CoinSymbol)
	if err != nil {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			params.FromAddress, params.ToAddress, params.CoinSymbol, params.Amount)
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
	meCoin, err = model.GetDB().GetAccountCoinsByWalletAddress(params.ToAddress, params.CoinSymbol)
	if err != nil {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			params.FromAddress, params.ToAddress, params.CoinSymbol, params.Amount)
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	if meCoin.CoinID == 0 {
		log.Errorf("not exist deposit info fromAddr:%v, toAddr:%v, symbol:%v, amount:%v",
			params.FromAddress, params.ToAddress, params.CoinSymbol, params.Amount)
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	// USPAU_Mod_AccountCoins 호출 하여 코인량 갱신
	adjustQuantity, _ := strconv.ParseFloat(params.Amount, 64)

	if err := model.GetDB().UpdateAccountCoins(
		meCoin.AUID,
		meCoin.CoinID,
		model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
		params.ToAddress,
		meCoin.Quantity,
		adjustQuantity,
		meCoin.Quantity+adjustQuantity,
		context.LogID_external_wallet,
		context.EventID_add,
		params.TxId); err != nil {
		log.Errorf("UpdateAccountCoins error : %v", err)
	}
	return resp
}

func TransferResultWithdrawal(params *context.ReqCoinTransferResWithdrawal) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	tKey := model.MakeCoinTransferKeyByTxID(params.Txid)
	txType, err := model.GetDB().GetCacheCoinTransferTx(tKey)
	if err != nil {
		// 존재 하지 않는 출금 정보 콜백을 받았다.
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Invalid_transfer_txid]+" txid:%v from:%v to:%v amount:%v",
			params.Txid, params.FromAddress, params.ToAddress, params.Amount)
		resp.SetReturn(resultcode.Result_Invalid_transfer_txid)
		return resp
	}

	// 실패한 tx면 관련 redis 찾아서 삭제 한다.
	if !strings.EqualFold(params.Status, "success") {
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

		log.Errorf("transfer txid fail txid:%v", params.Txid)

		// 아래 두경우 swap 복구 루틴이 똑같다.
		if txType.Target == context.From_user_to_fee_wallet || txType.Target == context.From_parent_to_other_wallet {
			// swap 정보 복구 하고 redis 삭제
			swapKey := model.MakeSwapKey(txType.AUID)
			if reqSwapInfo, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
				log.Errorf("GetCacheSwapInfo err => auid:%v txid:%v", err, txType.AUID, params.Txid)
				return resp
			} else {
				// 최신 코인 수량 정보 수집 해서 swap 시도 하려 했던 만큼 더해준다.
				_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
				if err != nil {
					log.Errorf("GetAccountCoins error : %v, auid:%v txid:%v", err, txType.AUID, params.Txid)
					model.MakeDbError(resp, resultcode.Result_DBError, err)
					return resp
				}
				meCoin := coinsMap[reqSwapInfo.CoinID]
				reqSwapInfo.PreviousCoinQuantity = meCoin.Quantity
				reqSwapInfo.CoinQuantity = meCoin.Quantity - reqSwapInfo.AdjustCoinQuantity
				reqSwapInfo.AdjustCoinQuantity = -(reqSwapInfo.AdjustCoinQuantity)
				// 최신 포인트 정보 수집
				nowPointQuantity, err := model.GetDB().GetPointApp(reqSwapInfo.MUID, reqSwapInfo.AppID, reqSwapInfo.DatabaseID)
				if err != nil {
					log.Errorf("GetPointApp error : %v, auid:%v muid:%v txid:%v", err, txType.AUID, reqSwapInfo.MUID, params.Txid)
					model.MakeDbError(resp, resultcode.Result_DBError, err)
					return resp
				}
				reqSwapInfo.PreviousPointQuantity = nowPointQuantity
				reqSwapInfo.PointQuantity = nowPointQuantity - reqSwapInfo.AdjustPointQuantity
				reqSwapInfo.AdjustPointQuantity = -(reqSwapInfo.AdjustPointQuantity)
				// event id 역으로 바꿔준다
				if reqSwapInfo.EventID == context.EventID_toCoin {
					reqSwapInfo.EventID = context.EventID_toPoint
				} else if reqSwapInfo.EventID == context.EventID_toPoint {
					reqSwapInfo.EventID = context.EventID_toCoin
				}

				if err := model.GetDB().PostPointCoinSwap(reqSwapInfo, params.Txid); err != nil {
					resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
					return resp
				}

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
		_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		if err != nil {
			log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, params.Txid)
			model.MakeDbError(resp, resultcode.Result_DBError, err)
			return resp
		}
		meCoin, ok := coinsMap[txType.CoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
			return resp
		}

		amount, _ := strconv.ParseFloat(params.Amount, 64)
		fee, _ := strconv.ParseFloat(params.ActualFee, 64)
		if err := model.GetDB().UpdateAccountCoins(
			txType.AUID,
			txType.CoinID,
			model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
			meCoin.WalletAddress,
			meCoin.Quantity,
			-(amount + fee), // 전송 수수료 + amount
			meCoin.Quantity-(amount+fee),
			context.LogID_external_wallet,
			context.EventID_sub,
			params.Txid); err != nil {
			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v", err, resp.Return, resp.Message, params.Txid)
			return resp
		}
		meCoin.Quantity = meCoin.Quantity - (amount + fee) // 남은 수량

		// redis swap 정보 찾아서 부모지갑에서 자식지갑으로 코인 전송
		swapKey := model.MakeSwapKey(txType.AUID)
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
		swapKey := model.MakeSwapKey(txType.AUID)
		parentKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
		_, err1 := model.GetDB().GetCacheSwapInfo(swapKey)
		_, err2 := model.GetDB().GetCacheCoinTransferFromParentWallet(parentKey)
		if err1 != nil && err2 != nil {
			log.Errorf("GetCacheSwapInfo err1 => %v auid:%v txid:%v", err1, txType.AUID, params.Txid)
			log.Errorf("GetCacheCoinTransferFromParentWallet err2 => %v auid:%v txid:%v", err2, txType.AUID, params.Txid)
			return resp
		}

		if err2 == nil {
			// 부모지갑사용 코인이 아니면 패스
			coin := model.GetDB().Coins[txType.CoinID]
			if baseCoin, ok := model.GetDB().BaseCoinMapByCoinID[coin.BaseCoinID]; ok {
				if baseCoin.IsUsedParentWallet {
					// 부모지갑출금을 사용하는 코인 외부 지갑 전송 성공으로 처리하고 그 보낸 사람의 코인을 차감 시킨다.
					// 자식 지갑의 코인 빼고 수수료 코인은 빼지 말고 일단 주석처리한다.
					_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
					if err != nil {
						log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, params.Txid)
						model.MakeDbError(resp, resultcode.Result_DBError, err)
						return resp
					}

					meCoin, ok := coinsMap[txType.CoinID]
					if !ok {
						log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
						return resp
					}

					amount, _ := strconv.ParseFloat(params.Amount, 64)
					if err := model.GetDB().UpdateAccountCoins(
						txType.AUID,
						txType.CoinID,
						model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
						meCoin.WalletAddress,
						meCoin.Quantity,
						-amount, // amount
						meCoin.Quantity-amount,
						context.LogID_external_wallet,
						context.EventID_sub,
						params.Txid); err != nil {
						log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v amount:%v", err, resp.Return, resp.Message, params.Txid, -amount)
						return resp
					}
					// baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID]
					// meBaseCoinID := int64(0)
					// for _, coin := range model.GetDB().Coins {
					// 	if strings.EqualFold(coin.CoinSymbol, baseCoinInfo.BaseCoinSymbol) {
					// 		meBaseCoinID = coin.CoinId
					// 	}
					// }

					// meCoin.Quantity = meCoin.Quantity - amount // 남은 수량
					// meBaseCoin, ok := coinsMap[meBaseCoinID]
					// if !ok {
					// 	log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
					// 	return resp
					// }

					// fee, _ := strconv.ParseFloat(params.ActualFee, 64)
					// if err := model.GetDB().UpdateAccountCoins(
					// 	txType.AUID,
					// 	meBaseCoinID,
					// 	model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
					// 	meBaseCoin.WalletAddress,
					// 	meBaseCoin.Quantity,
					// 	-fee, // amount
					// 	meBaseCoin.Quantity-fee,
					// 	context.LogID_external_wallet,
					// 	context.EventID_sub,
					// 	params.Txid); err != nil {
					// 	log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v fee:%v", err, resp.Return, resp.Message, params.Txid, fee)
					// 	return resp
					// }
				}
			}

		}
		// parent가 보낸 redis 정보는 삭제 해준다.
		model.GetDB().DelCacheCoinTransfer(tKey)                      //tx 삭제
		model.GetDB().DelCacheCoinTransferFromParentWallet(parentKey) // from user 삭제
		model.GetDB().DelCacheSwapInfo(swapKey)

		// if _, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
		// 	log.Errorf("GetCacheSwapInfo err => audi:%v txid:%v", err, txType.AUID, params.Txid)
		// 	return resp
		// }
		// // parent가 보낸 redis 정보는 삭제 해준다.
		// model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		// model.GetDB().DelCacheCoinTransferFromParentWallet(parentKey) // from user 삭제
		// model.GetDB().DelCacheSwapInfo(swapKey)
	} else if txType.Target == context.From_user_to_parent_wallet { // 자식 지갑에서 부모 지갑으로 전송 : swap coin->point
		// 자식 지갑의 basecoin 수량 acutual fee 만큼 축소
		_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		if err != nil {
			log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, params.Txid)
			model.MakeDbError(resp, resultcode.Result_DBError, err)
			return resp
		}

		meCoin, ok := coinsMap[txType.CoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
			return resp
		}

		baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID]
		meBaseCoinID := int64(0)
		for _, coin := range model.GetDB().Coins {
			if strings.EqualFold(coin.CoinSymbol, baseCoinInfo.BaseCoinSymbol) {
				meBaseCoinID = coin.CoinId
			}
		}

		meBaseCoin, ok := coinsMap[meBaseCoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
			return resp
		}

		fee, _ := strconv.ParseFloat(params.ActualFee, 64)
		if err := model.GetDB().UpdateAccountCoins(
			txType.AUID,
			meBaseCoinID,
			model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
			meBaseCoin.WalletAddress,
			meBaseCoin.Quantity,
			-fee, // amount
			meBaseCoin.Quantity-fee,
			context.LogID_external_wallet,
			context.EventID_sub,
			params.Txid); err != nil {
			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v", err, resp.Return, resp.Message, params.Txid)
			return resp
		}
		meBaseCoin.Quantity = meCoin.Quantity - fee // 남은 수량

		// redis swap 정보 찾아서 스왑 처리
		swapKey := model.MakeSwapKey(txType.AUID)
		if _, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
			log.Errorf("GetCacheSwapInfo err => audi:%v txid:%v", err, txType.AUID, params.Txid)
			return resp
		}
		// user가 보낸 redis 정보는 삭제 해준다. 이때는 이미 from parent reids 정보가 생성되어 있음
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
		model.GetDB().DelCacheSwapInfo(swapKey)
	} else if txType.Target == context.From_user_to_other_wallet { // 자식지갑에서 다른 지갑으로 코인 전송
		// 자식 지갑의 코인 빼고, 수수료 코인 뺀다.
		_, coinsMap, err := model.GetDB().GetAccountCoins(txType.AUID)
		if err != nil {
			log.Errorf("GetAccountCoins error : %v, audi:%v txid:%v", err, txType.AUID, params.Txid)
			model.MakeDbError(resp, resultcode.Result_DBError, err)
			return resp
		}

		meCoin, ok := coinsMap[txType.CoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
			return resp
		}

		baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID]
		meBaseCoinID := int64(0)
		for _, coin := range model.GetDB().Coins {
			if strings.EqualFold(coin.CoinSymbol, baseCoinInfo.BaseCoinSymbol) {
				meBaseCoinID = coin.CoinId
			}
		}

		meBaseCoin, ok := coinsMap[meBaseCoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v txid:%v", txType.CoinID, params.Txid)
			return resp
		}

		amount, _ := strconv.ParseFloat(params.Amount, 64)
		if err := model.GetDB().UpdateAccountCoins(
			txType.AUID,
			txType.CoinID,
			model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
			meCoin.WalletAddress,
			meCoin.Quantity,
			-amount, // amount
			meCoin.Quantity-amount,
			context.LogID_external_wallet,
			context.EventID_sub,
			params.Txid); err != nil {
			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v amount:%v", err, resp.Return, resp.Message, params.Txid, -amount)
			return resp
		}
		meCoin.Quantity = meCoin.Quantity - amount // 남은 수량

		fee, _ := strconv.ParseFloat(params.ActualFee, 64)
		if err := model.GetDB().UpdateAccountCoins(
			txType.AUID,
			meBaseCoinID,
			model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
			meBaseCoin.WalletAddress,
			meBaseCoin.Quantity,
			-fee, // amount
			meBaseCoin.Quantity-fee,
			context.LogID_external_wallet,
			context.EventID_sub,
			params.Txid); err != nil {
			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v fee:%v", err, resp.Return, resp.Message, params.Txid, fee)
			return resp
		}
		meBaseCoin.Quantity = meCoin.Quantity - fee // 남은 수량

		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
	}

	return resp
}

func IsExistInprogressTransferFromParentWallet(params *context.GetCoinTransferExistInProgress) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	Lockkey := model.MakeCoinTransferFromParentWalletLockKey(params.AUID)
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

	key := model.MakeCoinTransferFromParentWalletKey(params.AUID)
	reqCoinTransfer, err := model.GetDB().GetCacheCoinTransferFromParentWallet(key)
	if err == nil {
		// 전송중인 기존 정보가 있다면 값을 추가해준다.
		resp.Value = reqCoinTransfer
	} else {
		//log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_NotExistInprogress])
		resp.SetReturn(resultcode.Result_Error_Transfer_NotExistInprogress)
	}

	return resp
}

func IsExistInprogressTransferFromUserWallet(params *context.GetCoinTransferExistInProgress) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	Lockkey := model.MakeCoinTransferFromUserWalletLockKey(params.AUID)
	mutex := model.GetDB().RedSync.NewMutex(Lockkey)
	if err := mutex.Lock(); err != nil {
		log.Error("redis lock err:%v", err)
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	}

	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			if err != nil {
				log.Errorf("unlock err : %v", err)
			}
		}
	}()

	key := model.MakeCoinTransferFromUserWalletKey(params.AUID)
	reqCoinTransfer, err := model.GetDB().GetCacheCoinTransferFromUserWallet(key)
	if err == nil {
		// 전송중인 기존 정보가 있다면 값을 추가해준다.
		resp.Value = reqCoinTransfer
	} else {
		//log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_NotExistInprogress])
		resp.SetReturn(resultcode.Result_Error_Transfer_NotExistInprogress)
	}

	return resp
}

func CoinReload(params *context.CoinReload) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	Lockkey := model.MakeCoinTransferFromUserWalletLockKey(params.AUID)
	mutex := model.GetDB().RedSync.NewMutex(Lockkey)
	isValid, _ := mutex.Valid()
	if isValid {
		log.Errorf("auid:%v %v", params.AUID, resultcode.ResultCodeText[resultcode.Result_RedisError_WaitForProcessing])
		resp.SetReturn(resultcode.Result_RedisError_WaitForProcessing)
		return resp
	}
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

	meCoins, _, err := model.GetDB().GetAccountCoins(params.AUID)
	if err != nil {
		log.Errorf("GetAccountCoins error : %v, auid:%v", err, params.AUID)
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		for _, coin := range meCoins {
			if model.GetDB().BaseCoinMapByCoinID[coin.BaseCoinID].IsUsedParentWallet {
				continue
			}

			req := &token_manager_server.ReqBalance{
				Symbol:  model.GetDB().Coins[coin.CoinID].CoinSymbol,
				Address: coin.WalletAddress,
			}

			if res, err := token_manager_server.GetInstance().GetBalance(req); err != nil {
				resp.SetReturn(resultcode.ResultInternalServerError)
				return resp
			} else {
				if res.Return != 0 { // token manager 전송 에러
					resp.Return = res.Return
					resp.Message = res.Message
				} else {
					resp.Value = meCoins
					// 내 코인 수량과 비교해서 다르면 업데이트
					newQuantity, err := strconv.ParseFloat(res.ResReqBalanceValue.Balance, 64)

					if err != nil {
						log.Errorf("new coin balance parse err : %v", err)
					} else if coin.Quantity != newQuantity {
						adjustCoinAmount := coin.Quantity - newQuantity
						adjustCoinAmount = toFixed(adjustCoinAmount, 9)
						if adjustCoinAmount == 0 {
							continue
						}

						eventID := context.EventID_add
						if adjustCoinAmount > 0 {
							eventID = context.EventID_sub
						}

						if err := model.GetDB().UpdateAccountCoins(
							params.AUID,
							coin.CoinID,
							model.GetDB().Coins[coin.CoinID].BaseCoinID,
							coin.WalletAddress,
							coin.Quantity,
							-adjustCoinAmount,
							coin.Quantity-adjustCoinAmount,
							context.LogID_wallet_sync,
							context.EventID_type(eventID),
							"coin reload"); err != nil {
							log.Errorf("UpdateAccountCoins error : %v", err)
						}

						coin.Quantity = coin.Quantity - adjustCoinAmount // 응답값 수정
					}
				}
			}
		}
	}

	return resp
}
