package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
	"github.com/labstack/echo"
)

func PostPSMaintenance(ctx *context.PointManagerContext, req *context.PSMaintenance) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.SetMaintenance(req.Enable)

	msg := &model.PSMaintenance{
		PSHeader: model.PSHeader{
			Type: model.PubSub_type_maintenance,
		},
	}
	msg.Value.Enable = req.Enable

	if err := model.GetDB().PublishEvent(model.InternalCmd, msg); err != nil {
		log.Errorf("PublishEvent %v, type : %v, error : %v", model.InternalCmd, model.PubSub_type_maintenance, err)
		resp.SetReturn(resultcode.Result_PubSub_InternalErr)
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostPSSwap(ctx *context.PointManagerContext, req *context.PSSwap) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.SetSwapEnable(req.Enable)

	msg := &model.PSSwap{
		PSHeader: model.PSHeader{
			Type: model.PubSub_type_Swap,
		},
	}
	msg.Value.Enable = req.Enable

	if err := model.GetDB().PublishEvent(model.InternalCmd, msg); err != nil {
		log.Errorf("PublishEvent %v, type : %v, error : %v", model.InternalCmd, model.PubSub_type_Swap, err)
		resp.SetReturn(resultcode.Result_PubSub_InternalErr)
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostPSCoinTransferExternal(ctx *context.PointManagerContext, req *context.PSCoinTransferExternal) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.SetExternalTransferEnable(req.Enable)

	msg := &model.PSCoinTransferExternal{
		PSHeader: model.PSHeader{
			Type: model.PubSub_type_CoinTransferExternal,
		},
	}
	msg.Value.Enable = req.Enable

	if err := model.GetDB().PublishEvent(model.InternalCmd, msg); err != nil {
		log.Errorf("PublishEvent %v, type : %v, error : %v", model.InternalCmd, model.PubSub_type_CoinTransferExternal, err)
		resp.SetReturn(resultcode.Result_PubSub_InternalErr)
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostPSMetaRefresh(ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.GetDB().GetPointList()
	model.GetDB().GetAppCoins()
	model.GetDB().GetCoins()
	model.GetDB().GetApps()
	model.GetDB().GetAppPoints()
	model.GetDB().GetBaseCoins()

	msg := &model.PSMetaRefresh{
		PSHeader: model.PSHeader{
			Type: model.PubSub_type_meta_refresh,
		},
	}

	if err := model.GetDB().PublishEvent(model.InternalCmd, msg); err != nil {
		log.Errorf("PublishEvent %v, type : %v, error : %v", model.InternalCmd, model.PubSub_type_meta_refresh, err)
		resp.SetReturn(resultcode.Result_PubSub_InternalErr)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostPSPointUpdate(ctx *context.PointManagerContext, req *context.PSPointUpdate) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.SetPointUpdateEnable(req.Enable)

	msg := &model.PSPointUpdate{
		PSHeader: model.PSHeader{
			Type: model.PubSub_type_point_update,
		},
	}
	msg.Value.Enable = req.Enable

	if err := model.GetDB().PublishEvent(model.InternalCmd, msg); err != nil {
		log.Errorf("PublishEvent %v, type : %v, error : %v", model.InternalCmd, model.PubSub_type_point_update, err)
		resp.SetReturn(resultcode.Result_PubSub_InternalErr)
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetMeta(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	swapList := context.Meta{
		PointList: model.GetDB().ScanPointsMap,
		AppCoins:  model.GetDB().ScanPointsOfApp,
		Coins:     model.GetDB().Coins,
	}

	resp.Value = swapList

	return c.JSON(http.StatusOK, resp)
}
