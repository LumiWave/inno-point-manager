package commonapi

import (
	"net/http"
	"strings"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/token_manager_server"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
	"github.com/labstack/echo"
)

func PostCoinTransferFromParentWallet(params *context.ReqCoinTransferFromParentWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if params.IsNormalTransfer {
		// 부모지갑에서 단순 출금용
		resp = inner.TransferFromParentWalletNormal(params, params.IsNormalTransfer)
	} else {
		if err := inner.TransferFromParentWallet(params, true); err != nil {
			resp = err
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinTransferFromUserWallet(params *context.ReqCoinTransferFromUserWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if !model.GetExternalTransferEnable() {
		resp.SetReturn(resultcode.Result_Error_IsCoinTransferExternalMaintenance)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	params.Target = context.From_user_to_other_wallet
	if err := inner.TransferFromUserWallet(params, true); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 코인 외부 지갑 전송 중인 상태 정보 요청
func GetCoinTransferExistInProgress(params *context.GetCoinTransferExistInProgress, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if res := inner.IsExistInprogressTransferFromParentWallet(params); res != nil {
		resp = res
	}
	if resp.Return != 0 {
		if res := inner.IsExistInprogressTransferFromUserWallet(params); res != nil {
			resp = res
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetCoinObjects(req *context.ReqCoinObjects, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	params := &token_manager_server.ReqCoinObjects{
		WalletAddress:   req.WalletAddress,
		ContractAddress: req.ContractAddress,
	}

	if res, err := token_manager_server.GetInstance().GetCoinObjectIDS(params); err != nil {
		log.Errorf("GetCoinObjectIDS err : %v,  wallet:%v, contract:%v ", err, req.WalletAddress, req.ContractAddress)
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else {
		if res.Common.Return == 0 {
			resValue := new(context.ResCoinObjects)
			resValue.ObjectIDs = res.Value.ObjectIDs
			resp.Value = resValue
		} else {
			log.Errorf("GetCoinObjectIDS error return : %v, %s, wallet:%v, contract:%v", res.Common.Return, res.Common.Message, req.WalletAddress, req.ContractAddress)
			resp.SetReturn(resultcode.ResultInternalServerError)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetCoinFee(params *context.ReqCoinFee, c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	req := &token_manager_server.ReqCoinFee{
		Symbol: params.Symbol,
	}

	if res, err := token_manager_server.GetInstance().GetCoinFee(req); err != nil {
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else {
		if res.Return != 0 { // token manager 전송 에러
			resp.Return = res.Return
			resp.Message = res.Message
		} else {
			resp.Value = res.ResCoinFeeInfoValue
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetBalance(params *context.ReqBalance, c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	coinInfo, ok := model.GetDB().CoinsBySymbol[params.Symbol]
	if !ok {
		resp.SetReturn(resultcode.Result_Require_Symbol)
		return c.JSON(http.StatusOK, resp)
	}
	baseCoinInfo := model.GetDB().BaseCoinMapByCoinID[coinInfo.BaseCoinID]
	req := &token_manager_server.ReqBalance{
		BaseSymbol: baseCoinInfo.BaseCoinSymbol,
		Contract: func() string {
			// 코인 타입이면 contract 정보를 를 보내지 않는다.
			if strings.EqualFold(baseCoinInfo.BaseCoinSymbol, coinInfo.CoinSymbol) {
				return ""
			}
			return coinInfo.ContractAddress
		}(),
		Address: params.Address,
	}

	if res, err := token_manager_server.GetInstance().GetBalance(req); err != nil {
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else {
		if res.Return != 0 { // token manager 전송 에러
			resp.Return = res.Return
			resp.Message = res.Message
		} else {
			resp.Value = res.ResReqBalanceValue
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetBalanceAll(params *context.ReqBalanceAll, ctx *context.PointManagerContext) error {
	resp := inner.GetBalanceAll(params.AUID)
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinReload(params *context.CoinReload, ctx *context.PointManagerContext) error {
	resp := inner.CoinReload(params)
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
