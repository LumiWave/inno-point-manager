package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

func PostPointCoinSwap(params *context.ReqSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

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

	log.Debugf("GetSwapInprogressNotExist auid : %v", params.AUID)
	// 내 지갑 정보를 가져와서 모든 지갑을 뒤져버 진행 중에 있는지 캐시 정보로 체크
	swapInfos := []*context.ReqSwapInfo{}
	mapWallet := make(map[string]string)
	if wallets, _, err := model.GetDB().USPAU_GetList_AccountWallets(params.AUID); err == nil {
		for _, wallet := range wallets {
			if _, ok := mapWallet[wallet.WalletAddress]; !ok {
				if swapInfo, err := model.GetDB().CacheGetSwapWallet(wallet.WalletAddress); err == nil {
					swapInfos = append(swapInfos, swapInfo)
					if swapInfo.TxType == context.EventID_P2C {
						mapWallet[swapInfo.SwapToCoin.WalletAddress] = swapInfo.SwapToCoin.WalletAddress
					} else if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_C2C {
						mapWallet[swapInfo.SwapFromCoin.WalletAddress] = swapInfo.SwapFromCoin.WalletAddress
					}
				}
			}
		}
		if len(swapInfos) > 0 {
			resp.Value = swapInfos
			resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
		}
	} else {
		log.Errorf("USPAU_GetList_AccountWallets err : %v, auid:%v", err, params.AUID)
		resp.SetReturn(resultcode.Result_Error_Db_GetAccountWallets)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteDeleteSwapInfo(params *context.DeleteDeleteSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if len(params.WalletAddress) == 0 {
		model.GetDB().CacheDelAllSwapWallet()
	} else if err := model.GetDB().CacheDelSwapWallet(params.WalletAddress); err != nil {
		log.Errorf("CacheDelSwapWallet err:%v, wallet:%v", err, params.WalletAddress)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
