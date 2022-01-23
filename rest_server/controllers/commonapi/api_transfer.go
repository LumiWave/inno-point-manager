package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func PostCoinTransfer(params *context.ReqCoinTransfer, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.Transfer(params); err != nil {
		resp = err
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
