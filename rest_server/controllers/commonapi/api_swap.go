package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func PostPointCoinSwap(params *context.ReqSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// if !model.GetSwapEnable() {
	// 	resp.SetReturn(resultcode.Result_Error_IsSwapMaintenance)
	// } else if err := inner.Swap(params, ctx.GetValue().InnoUID); err != nil {
	// 	resp = err
	// }
	if !model.GetSwapEnable() {
		resp.SetReturn(resultcode.Result_Error_IsSwapMaintenance)
	} else if err := inner.SwapWallet(params, ctx.GetValue().InnoUID); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PutSwapGasFee(params *context.ReqSwapGasFee, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.SwapGasFee(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
