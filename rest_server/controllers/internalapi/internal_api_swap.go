package internalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// app app 포인트 코인 swap 처리 요청
func (o *InternalAPI) PostPointCoinSwap(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqSwapInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostPointCoinSwap(params, ctx)
}

// swap 상태 정보 처리 요청
func (o *InternalAPI) PutSwapStatus(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqSwapStatus()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PutSwapStatus(params, ctx)
}

// swap 진행 상태 정보 조회
func (o *InternalAPI) GetSwapInprogressNotExist(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqSwapIniprogress()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetSwapInprogressNotExist(params, ctx)
}

// swap 정보 삭제
func (o *InternalAPI) DeleteDeleteSwapInfo(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := &context.DeleteDeleteSwapInfo{}
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.DeleteDeleteSwapInfo(params, ctx)
}
