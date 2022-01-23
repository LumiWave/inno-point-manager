package inner

import (
	"strconv"

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
	Lockkey := model.MakeMemberPointListLockKey(params.AUID)
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
