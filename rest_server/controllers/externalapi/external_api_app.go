package externalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// app 포인트 리스트 조회
func (o *ExternalAPI) GetPointAppList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqGetPointApp()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	params.AppId = 0
	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointAppList(params, ctx)
}

// app 포인트 조회
func (o *ExternalAPI) GetPointApp(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqGetPointApp()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	params.AppId = 0
	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return commonapi.GetPointApp(params, ctx)
}

// app 포인트 업데이트
func (o *ExternalAPI) PutPointAppUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqPointMemberAppUpdate()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PutPointAppUpdate(params, ctx)
}
