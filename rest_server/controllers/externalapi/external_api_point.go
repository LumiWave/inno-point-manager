package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// app 포인트 업데이트
func (o *ExternalAPI) PutPointAppUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqPointMemberAppUpdate()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PutPointAppUpdate(params, ctx)
}

// app 포인트 수신 이력
func (o *ExternalAPI) GetPointAppHistory(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPointMemberHistoryReq()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointAppHistory(params, ctx)
}

// app 포인트 전환 요청
func (o *ExternalAPI) PostPointAppExchange(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPostPointAppExchange()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostPointAppExchange(params, ctx)
}

// app 포인트 전환 이력
func (o *ExternalAPI) GetPointAppExchangeHistory(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPointMemberExchangeHistory()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(true); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointAppExchangeHistory(params, ctx)
}
