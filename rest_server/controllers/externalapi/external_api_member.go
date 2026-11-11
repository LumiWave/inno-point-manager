package externalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 내 지갑 정보 조회
func (o *ExternalAPI) GetPointMemberWallet(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewPointMemberWallet()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// 회원 토큰은 요청의 au_id 와 상관없이 본인 계정으로만 조회한다.
	if ctx.IsMemberToken() {
		params.AUID = ctx.GetValue().AUID
	}
	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	// 앱 토큰은 그 앱에 가입한 계정의 지갑만 조회할 수 있다.
	if err := commonapi.CheckWalletAccess(ctx, params); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetPointMemberWallet(params, ctx)
}
