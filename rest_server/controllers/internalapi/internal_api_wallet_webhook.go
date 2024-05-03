package internalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *InternalAPI) PostWalletWebHookETHDeposit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqPostWalletETHResult()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWalletWebHookETHDeposit(params, ctx)
}

func (o *InternalAPI) PostWalletWebHookETHWithdrawal(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqPostWalletETHResult()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWalletWebHookETHWithdrawal(params, ctx)
}

func (o *InternalAPI) PostWalletWebHookSUIDeposit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewSUI_CB_Balance_Changes()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWalletWebHookSUIDeposit(params, ctx)
}

func (o *InternalAPI) PostWalletWebHookSUIWithdrawal(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewSUI_CB_Balance_Changes()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWalletWebHookSUIWithdrawal(params, ctx)
}
