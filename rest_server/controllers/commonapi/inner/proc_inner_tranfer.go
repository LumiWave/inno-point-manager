package inner

import (
	"fmt"
	"math/big"
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

	coinInfo := model.GetDB().CoinsBySymbol[params.CoinSymbol]
	scale := new(big.Float).SetInt64(1)
	scale.SetString("1e" + fmt.Sprintf("%d", coinInfo.Decimal))

	valueAmount := new(big.Float).SetFloat64(params.Quantity)
	valueAmount = new(big.Float).Mul(valueAmount, scale)
	nAmount := new(big.Int)
	valueAmount.Int(nAmount)
	amount := nAmount.String()

	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	req := &token_manager_server.ReqSendFromParentWallet{
		BaseSymbol: model.GetDB().BaseCoinMapByCoinID[coinInfo.BaseCoinID].BaseCoinSymbol,
		Contract: func() string {
			// 코인 타입이면 contract 정보를 를 보내지 않는다.
			if strings.EqualFold(model.GetDB().BaseCoinMapByCoinID[coinInfo.BaseCoinID].BaseCoinSymbol, coinInfo.CoinSymbol) {
				return ""
			}
			return coinInfo.ContractAddress
		}(),
		ToAddress: params.ToAddress,
		Amount:    amount,
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

	decimal := model.GetDB().CoinsBySymbol[params.CoinSymbol].Decimal
	scale := new(big.Float).SetInt64(1)
	scale.SetString("1e" + fmt.Sprintf("%d", decimal))

	valueAmount := new(big.Float).SetFloat64(params.Quantity)
	valueAmount = new(big.Float).Mul(valueAmount, scale)
	nAmount := new(big.Int)
	valueAmount.Int(nAmount)
	amount := nAmount.String()

	//2. tokenmanager에 외부 전송 요청, 전송 transaction 유효한지 확인
	req := &token_manager_server.ReqSendFromUserWallet{
		BaseCoinSymbol: params.BaseCoinSymbol,
		Contract: func() string {
			// 코인 타입이면 contract 정보를 를 보내지 않는다.
			if strings.EqualFold(params.BaseCoinSymbol, model.GetDB().Coins[params.CoinID].CoinSymbol) {
				return ""
			}
			return model.GetDB().Coins[params.CoinID].ContractAddress
		}(),
		FromAddress: params.FromAddress,
		ToAddress:   params.ToAddress,
		Amount:      amount,
		Memo:        strconv.FormatInt(params.AUID, 10),
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
				BaseSymbol: model.GetDB().BaseCoinMapByCoinID[coin.BaseCoinID].BaseCoinSymbol,
				Contract: func() string {
					// 코인 타입이면 contract 정보를 를 보내지 않는다.
					if strings.EqualFold(model.GetDB().BaseCoinMapByCoinID[coin.BaseCoinID].BaseCoinSymbol, model.GetDB().Coins[coin.CoinID].CoinSymbol) {
						return ""
					}
					return model.GetDB().Coins[coin.CoinID].ContractAddress
				}(),
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

					// scale := new(big.Float).SetFloat64(1)
					// scale.SetString("1e" + fmt.Sprintf("%d", res.Decimal))
					// valueAmount, _ := new(big.Float).SetString(res.Balance)
					// valueAmount = new(big.Float).Quo(valueAmount, scale)
					// newQuantity, _ := valueAmount.Float64()

					// //newQuantity, err := strconv.ParseFloat(res.ResReqBalanceValue.Balance, 64)

					// if err != nil {
					// 	log.Errorf("new coin balance parse err : %v", err)
					// } else if coin.Quantity != newQuantity {
					// 	adjustCoinAmount := coin.Quantity - newQuantity
					// 	adjustCoinAmount = toFixed(adjustCoinAmount, 9)
					// 	if adjustCoinAmount == 0 {
					// 		continue
					// 	}

					// 	eventID := context.EventID_add
					// 	if adjustCoinAmount > 0 {
					// 		eventID = context.EventID_sub
					// 	}

					// 	if err := model.GetDB().UpdateAccountCoins(
					// 		params.AUID,
					// 		coin.CoinID,
					// 		model.GetDB().Coins[coin.CoinID].BaseCoinID,
					// 		coin.WalletAddress,
					// 		coin.Quantity,
					// 		-adjustCoinAmount,
					// 		coin.Quantity-adjustCoinAmount,
					// 		context.LogID_wallet_sync,
					// 		context.EventID_type(eventID),
					// 		"coin reload"); err != nil {
					// 		log.Errorf("UpdateAccountCoins error : %v", err)
					// 	}

					// 	coin.Quantity = coin.Quantity - adjustCoinAmount // 응답값 수정
					//}
				}
			}
		}
	}

	return resp
}
