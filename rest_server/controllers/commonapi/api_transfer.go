package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func PostCoinTransferFromParentWallet(params *context.ReqCoinTransferFromParentWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.TransferFromParentWallet(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinTransferFromUserWallet(params *context.ReqCoinTransferFromUserWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.TransferFromUserWallet(params); err != nil {
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
	if res := inner.IsExistInprogressTransferFromUserWallet(params); res != nil {
		resp = res
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinTransferResultDeposit(params *context.ReqCoinTransferResDeposit, c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.TransferResultDeposit(params); err != nil {
		resp = err
	}

	return c.JSON(http.StatusOK, resp)
}

func PostCoinTransferResultWithdrawal(params *context.ReqCoinTransferResWithdrawal, c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.TransferResultWithdrawal(params); err != nil {
		resp = err
	}

	return c.JSON(http.StatusOK, resp)
}
