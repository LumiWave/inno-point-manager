package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// app 포인트 조회
func (o *InternalAPI) GetPointApp(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqGetPointApp()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointApp(params, ctx)
}

// app 포인트 업데이트 처리 모니터링
func (o *InternalAPI) GetPointAppMonitoring(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqPointAppMonitoring()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointAppMonitoring(params, ctx)
}
