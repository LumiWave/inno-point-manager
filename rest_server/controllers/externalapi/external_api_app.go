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
	// 회원 토큰은 본인 것만 조회할 수 있다. 앱 토큰은 조회 결과에서 자기 앱 회원인지 확인한다.
	if err := commonapi.CheckPointAccess(ctx, params.MUID, params.DatabaseID); err != nil {
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
	// 회원 토큰은 본인 것만 조회할 수 있다. 앱 토큰은 조회 결과에서 자기 앱 회원인지 확인한다.
	if err := commonapi.CheckPointAccess(ctx, params.MUID, params.DatabaseID); err != nil {
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
	// 포인트 조정은 앱(게임 서버) 토큰만 할 수 있고, 그 앱의 포인트만 조정할 수 있다.
	// 회원 토큰으로 본문의 app_id, mu_id 를 바꿔 남의 포인트를 조정하지 못하게 한다.
	if err := commonapi.CheckPointUpdate(ctx, params); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PutPointAppUpdate(params, ctx)
}
