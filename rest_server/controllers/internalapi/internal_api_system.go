package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *InternalAPI) GetNodeMetric(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	return commonapi.GetNodeMetric(ctx)
}

func (o *InternalAPI) PostPSMaintenance(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	params := context.NewPSMaintenance()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostPSMaintenance(ctx, params)
}

func (o *InternalAPI) PostPSSwap(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	params := context.NewPSSwap()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostPSSwap(ctx, params)
}

func (o *InternalAPI) PostPSCoinTransferExternal(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	params := context.NewPSCoinTransferExternal()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostPSCoinTransferExternal(ctx, params)
}

func (o *InternalAPI) PostPSMetaRefresh(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	return commonapi.PostPSMetaRefresh(ctx)
}

func (o *InternalAPI) PostPSPointUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	params := context.NewPSPointUpdate()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostPSPointUpdate(ctx, params)
}

func (o *InternalAPI) GetMeta(c echo.Context) error {
	return commonapi.GetMeta(c)
}
