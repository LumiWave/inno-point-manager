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

func PutSwapStatus(params *context.ReqSwapStatus, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.PutSwapStatus(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetSwapInprogressNotExist(params *context.ReqSwapInprogress, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 내 지갑 정보를 가져와서 모든 지갑을 뒤져버 진행 중에 있는지 체크
	swapInfos := []*context.ReqSwapInfo{}
	mapWallet := make(map[string]string)
	if wallets, err := model.GetDB().USPAU_GetList_AccountWallets(params.AUID); err == nil {
		for _, wallet := range wallets {
			if _, ok := mapWallet[wallet.WalletAddress]; !ok {
				if swapInfo, err := model.GetDB().CacheGetSwapWallet(wallet.WalletAddress); err == nil {
					swapInfos = append(swapInfos, swapInfo)
					mapWallet[swapInfo.WalletAddress] = swapInfo.WalletAddress
				}
			}
		}
	}
	if len(swapInfos) > 0 {
		resp.Value = swapInfos
		resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
