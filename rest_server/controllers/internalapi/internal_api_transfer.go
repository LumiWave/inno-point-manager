package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *InternalAPI) PostCoinTransfer(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqCoinTransfer()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCoinTransfer(params, ctx)
}

func (o *InternalAPI) PostCoinTransferResultDeposit(c echo.Context) error {
	params := context.NewReqCoinTransferResDeposit()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCoinTransferResultDeposit(params, c)
}

func (o *InternalAPI) PostCoinTransferResultWithdrawal(c echo.Context) error {
	params := context.NewReqCoinTransferResWithdrawal()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCoinTransferResultWithdrawal(params, c)
}

// tranfer 중인 상태 정보 요청
func (o *InternalAPI) GetCoinTransferExistInProgress(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)
	params := context.NewGetCoinTransferExistInProgress()

	// Request json 파싱
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// 유효성 체크
	if err := params.CheckValidate(ctx); err != nil {
		log.Errorf("%v", err)
		return c.JSON(http.StatusOK, err)
	}
	return commonapi.GetCoinTransferExistInProgress(params, ctx)
}
