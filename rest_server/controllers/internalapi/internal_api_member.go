package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 새로운 포인트 맴버 등록
func (o *InternalAPI) PostPointMemberRegister(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPointMemberInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(true); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostPointMemberRegister(params, ctx)
}

// 포인트 맴버 정보 업데이트
func (o *InternalAPI) PutPointMemberUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPointMemberInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(false); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PutPointMemberUpdate(params, ctx)
}
