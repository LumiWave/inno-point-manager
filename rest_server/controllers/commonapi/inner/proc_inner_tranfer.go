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

func Transfer(params *context.ReqCoinTransfer) *base.BaseResponse {

	resp := new(base.BaseResponse)
	resp.Success()

	//1. redis에 외부 전송이 진행 중인지 체크 redis에 정보가 있다면 전송중으로 인지하면됨
	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	//3. redis에 전송 정보 저장
	//4. 콜백(internal api)으로 완료or실패 확인 후 db 프로지저 호출, redis 삭제

	// 0. redis lock
	Lockkey := model.MakeCoinTransferLockKey(params.AUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	} else {
		// 0-1. redis unlock
		defer unLock()
	}

	// 1. redis에 외부 전송 정보 존재하는지 check
	key := model.MakeCoinTransferKey(params.AUID)
	_, err = model.GetDB().GetCacheCoinTransfer(key)
	if err == nil {
		// 전송중인 기존 정보가 있다면 에러를 리턴한다.
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_Transfer_Inprogress])
		resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
		return resp
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
	if err := model.GetDB().SetCacheCoinTransfer(key, params); err != nil {
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_RedisError_SetTransfer])
		resp.SetReturn(resultcode.Result_RedisError_SetTransfer)
		return resp
	}

	//4. redis에 전송 정보 transaction id key로 다시 한번더 저장 : 추후 콜백 api를 통해 검증하기 위해서
	tKey := model.MakeCoinTransferKeyByTxID(params.TransactionId)
	if err := model.GetDB().SetCacheCoinTransfer(tKey, params); err != nil {
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
	meCoin, err := model.GetDB().GetAccountCoinsByWalletAddress(params.ToAddress)
	if err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_GetAccountCoinByWalletAddress)
		return resp
	}

	// USPAU_Mod_AccountCoins 호출 하여 코인량 갱신
	adjustQuantity, _ := strconv.ParseFloat(params.Amount, 64)

	if err := model.GetDB().UpdateAccountCoins(
		meCoin.AUID,
		meCoin.CoinID,
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

	// tx로 redis 검색해서 load
	tKey := model.MakeCoinTransferKeyByTxID(params.Txid)
	transferInfo, err := model.GetDB().GetCacheCoinTransfer(tKey)
	if err != nil {
		// 존재 하지 않는 출금 정보 콜백을 받았다.
		log.Errorf(resultcode.ResultCodeText[resultcode.Result_Invalid_transfer_txid]+" txid:%v", params.Txid)
		resp.SetReturn(resultcode.Result_Invalid_transfer_txid)
		return resp
	}

	// 응답 status 성공 여부 체크
	if strings.EqualFold(params.Status, "success") {
		// from address, 이전 코인량 검색, 전송후 남은량 계산
		_, coinsMap, err := model.GetDB().GetAccountCoins(transferInfo.AUID)
		if err != nil {
			log.Errorf("GetAccountCoins error : %v", err)
			model.MakeDbError(resp, resultcode.Result_DBError, err)
		}

		meCoin, ok := coinsMap[transferInfo.CoinID]
		if !ok {
			log.Errorf("Not file my coinid : %v", transferInfo.CoinID)
			return resp
		}

		// USPAU_Mod_AccountCoins 호출 하여 코인 량 갱신
		if err := model.GetDB().UpdateAccountCoins(
			transferInfo.AUID,
			transferInfo.CoinID,
			meCoin.WalletAddress,
			meCoin.Quantity,
			-transferInfo.TotalQuantity,
			meCoin.Quantity-transferInfo.TotalQuantity,
			context.LogID_external_wallet,
			context.EventID_sub); err != nil {
			log.Errorf("UpdateAccountCoins error : %v", err)
		}
	} else if strings.EqualFold(params.Status, "failure") {
		// 실패 한 경우 두가지 redis 삭제만 유도한다.
		log.Warnf("coin withdrawal callback failure : %v", params.Txid)
	}

	// redis 두가지 삭제
	model.GetDB().DelCacheCoinTransfer(tKey) // txid key redis delete
	key := model.MakeCoinTransferKey(transferInfo.AUID)
	Lockkey := model.MakeMemberPointListLockKey(transferInfo.AUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	} else {
		// 0-1. redis unlock
		defer unLock()
	}
	model.GetDB().DelCacheCoinTransfer(key) // audi key redis delete

	return resp
}

func IsExistInprogressTransfer(params *context.GetCoinTransferExistInProgress) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	Lockkey := model.MakeCoinTransferLockKey(params.AUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	} else {
		// 0-1. redis unlock
		defer unLock()
	}

	key := model.MakeCoinTransferKey(params.AUID)
	reqCoinTransfer, err := model.GetDB().GetCacheCoinTransfer(key)
	if err == nil {
		// 전송중인 기존 정보가 있다면 값을 추가해준다.
		resp.Value = reqCoinTransfer
	}

	return resp
}
