package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

func PostPointCoinSwap(params *context.ReqSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.Swap(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetSwapExistInProgress(params *context.GetSwapExistInProgress, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if res := inner.IsExistInprogressSwap(params); res != nil {
		resp = res
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
