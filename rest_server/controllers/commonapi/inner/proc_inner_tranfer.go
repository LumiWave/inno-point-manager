package inner

import (
	"strconv"
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
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

	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	req := &token_manager_server.ReqSendFromParentWallet{
		Symbol:    params.CoinSymbol,
		ToAddress: params.ToAddress,
		Amount:    strconv.FormatFloat(params.Quantity, 'f', -1, 64),
		Memo:      strconv.FormatInt(params.AUID, 10),
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

		params.ReqId = res.Value.ReqId
		params.TransactionId = res.Value.TransactionId
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
			resp.Message = res.Message
			return resp
		}

		if len(res.Value.TransactionHash) == 0 {
			log.Errorf("PostSendFromUserWallet txid null")
		}
		//params.ReqId = res.Value.ReqId
		params.TransactionId = res.Value.TransactionHash
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

	// 입금 주소로 db 검색해서 AUID추출
	meCoin, err := model.GetDB().GetAccountCoinsByWalletAddress(params.ToAddress, params.CoinSymbol)
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
		context.EventID_add); err != nil {
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
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Invalid_transfer_txid]+" txid:%v", params.Txid)
		resp.SetReturn(resultcode.Result_Invalid_transfer_txid)
		return resp
	}

	// 실패한 tx면 관련 redis 찾아서 삭제 한다.
	if !strings.EqualFold(params.Status, "success") {
		log.Errorf("transfer txid fail txid:%v", params.Txid)
		if txType.Target == context.From_user_to_fee_wallet { // point -> coin swap 요청을 위한 1단계 수수료 출금 정보
			// txid redis 삭제
			model.GetDB().DelCacheCoinTransfer(tKey)
			// from user redis 삭제
			userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
			model.GetDB().DelCacheCoinTransferFromUserWallet(userKey)
			// swap redis 삭제
			swapKey := model.MakeSwapKey(txType.AUID)
			model.GetDB().DelCacheSwapInfo(swapKey)
		} else if txType.Target == context.From_parent_to_other_wallet {
			// 장애 처리를 위해 삭제 하지 않는다.
		}
		return resp
	}

	// 분기처리
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
			context.EventID_sub); err != nil {
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
	} else if txType.Target == context.From_parent_to_other_wallet { // 부모지갑에서 자식지갑으로 코인 전송 : swap point->coin
		// 부모지갑에서 자식지갑으로 입금시 swap으로 간주하고 swap 처리 진행한다.
		swapKey := model.MakeSwapKey(txType.AUID)
		if reqSwapInfo, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
			log.Errorf("GetCacheSwapInfo err => audi:%v txid:%v", err, txType.AUID, params.Txid)
			return resp
		} else {
			if err := model.GetDB().PostPointCoinSwap(reqSwapInfo); err != nil {
				resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
			}
		}

		// parent가 보낸 redis 정보는 삭제 해준다.
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromParentWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromParentWallet(userKey) // from user 삭제
		model.GetDB().DelCacheSwapInfo(swapKey)
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
			context.EventID_sub); err != nil {
			log.Errorf("UpdateAccountCoins error : %v Rceived a fee but failed to send coins => return:%v message:%v txid:%v", err, resp.Return, resp.Message, params.Txid)
			return resp
		}
		meCoin.Quantity = meCoin.Quantity - fee // 남은 수량

		// redis swap 정보 찾아서 스왑 처리
		swapKey := model.MakeSwapKey(txType.AUID)
		if reqSwapInfo, err := model.GetDB().GetCacheSwapInfo(swapKey); err != nil {
			log.Errorf("GetCacheSwapInfo err => audi:%v txid:%v", err, txType.AUID, params.Txid)
			return resp
		} else {
			if err := model.GetDB().PostPointCoinSwap(reqSwapInfo); err != nil {
				resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
			}
		}
		// user가 보낸 redis 정보는 삭제 해준다. 이때는 이미 from parent reids 정보가 생성되어 있음
		model.GetDB().DelCacheCoinTransfer(tKey) //tx 삭제
		userKey := model.MakeCoinTransferFromUserWalletKey(txType.AUID)
		model.GetDB().DelCacheCoinTransferFromUserWallet(userKey) // from user 삭제
	} else if txType.Target == context.From_user_to_other_wallet { // 자식지갑에서 다른 지갑으로 코인 전송

	}

	return resp

	//============================================================================

	// // tx로 redis 검색해서 load
	// tKey := model.MakeCoinTransferKeyByTxID(params.Txid)

	// // 부모지갑에서 출금된건지 특정 지갑에서 출금된건지 모르기 때문에 둘다 체크
	// var AUID, CoinID int64
	// var transferQuantity float64
	// transferInfoFromUser, err := model.GetDB().GetCacheCoinTransferFromUserWallet(tKey)
	// if err != nil {
	// 	// 존재 하지 않는 출금 정보 콜백을 받았다.
	// 	log.Errorf(resultcode.ResultCodeText[resultcode.Result_Invalid_transfer_txid]+" txid:%v", params.Txid)
	// 	resp.SetReturn(resultcode.Result_Invalid_transfer_txid)
	// 	return resp
	// }

	// AUID = transferInfoFromUser.AUID
	// CoinID = transferInfoFromUser.CoinID
	// transferQuantity = transferInfoFromUser.Quantity

	// // 응답 status 성공 여부 체크
	// if strings.EqualFold(params.Status, "success") {
	// 	// from address, 이전 코인량 검색, 전송후 남은량 계산
	// 	_, coinsMap, err := model.GetDB().GetAccountCoins(AUID)
	// 	if err != nil {
	// 		log.Errorf("GetAccountCoins error : %v", err)
	// 		model.MakeDbError(resp, resultcode.Result_DBError, err)
	// 	}

	// 	meCoin, ok := coinsMap[CoinID]
	// 	if !ok {
	// 		log.Errorf("Not file my coinid : %v", CoinID)
	// 		return resp
	// 	}

	// 	// USPAU_Mod_AccountCoins 호출 하여 코인 량 갱신
	// 	if err := model.GetDB().UpdateAccountCoins(
	// 		AUID,
	// 		CoinID,
	// 		model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
	// 		meCoin.WalletAddress,
	// 		meCoin.Quantity,
	// 		-transferQuantity,
	// 		meCoin.Quantity-transferQuantity,
	// 		context.LogID_external_wallet,
	// 		context.EventID_sub); err != nil {
	// 		log.Errorf("UpdateAccountCoins error : %v", err)
	// 	}
	// 	meCoin.Quantity = meCoin.Quantity - transferQuantity // 정보 갱신

	// 	// 자기 지갑에서 수수료 빼기
	// 	gasFeeCoinID := int64(0)
	// 	for _, coin := range model.GetDB().Coins {
	// 		if coin.CoinSymbol == model.GetDB().BaseCoinMapByCoinID[meCoin.BaseCoinID].BaseCoinSymbol {
	// 			gasFeeCoinID = coin.CoinId
	// 			break
	// 		}
	// 	}
	// 	gasCoin := coinsMap[gasFeeCoinID]
	// 	// 가스비 차감
	// 	gasFee, _ := strconv.ParseFloat(params.ActualFee, 64)
	// 	if err := model.GetDB().UpdateAccountCoins(
	// 		AUID,
	// 		gasFeeCoinID,
	// 		model.GetDB().Coins[meCoin.CoinID].BaseCoinID,
	// 		gasCoin.WalletAddress,
	// 		gasCoin.Quantity,
	// 		-gasFee,
	// 		gasCoin.Quantity-gasFee,
	// 		context.LogID_external_wallet,
	// 		context.EventID_sub); err != nil {
	// 		log.Errorf("UpdateAccountCoins gasfee error : %v", err)
	// 	}
	// } else if strings.EqualFold(params.Status, "failure") {
	// 	// 실패 한 경우 두가지 redis 삭제만 유도한다.
	// 	log.Warnf("coin withdrawal callback failure : %v", params.Txid)
	// }

	// // redis 두가지 삭제
	// model.GetDB().DelCacheCoinTransfer(tKey) // txid key redis delete

	// keyFromUser := model.MakeCoinTransferFromUserWalletKey(AUID)

	// Lockkey := model.MakeMemberPointListLockKey(AUID)
	// mutex := model.GetDB().RedSync.NewMutex(Lockkey)
	// if err := mutex.Lock(); err != nil {
	// 	log.Error("redis lock err:%v", err)
	// 	resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
	// 	return resp
	// }

	// defer func() {
	// 	// 1-1. redis unlock
	// 	if ok, err := mutex.Unlock(); !ok || err != nil {
	// 		if err != nil {
	// 			log.Errorf("unlock err : %v", err)
	// 		}
	// 	}
	// }()
	// model.GetDB().DelCacheCoinTransfer(keyFromUser) // audi key redis delete

	// return resp
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
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_NotExistInprogress])
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
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_NotExistInprogress])
		resp.SetReturn(resultcode.Result_Error_Transfer_NotExistInprogress)
	}

	return resp
}
