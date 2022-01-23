package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

func PostCoinTransfer(params *context.ReqCoinTransfer, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.Transfer(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinTransferResultDeposit(params *context.ReqCoinTransferResDeposit, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 입금 주소로 db 검색해서 AUID추출
	// USPAU_Mod_AccountCoins 호출 하여 코인 량 갱신

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostCoinTransferResultWithdrawal(params *context.ReqCoinTransferResWithdrawal, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// tx로 redis 검색해서 load 한 후 USPAU_Mod_AccountCoins 호출 하여 코인 량 갱신
	// redis 두가지 삭제

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
