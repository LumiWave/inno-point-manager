package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *InternalAPI) PostCoinTransferFromParentWallet(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqCoinTransferFromParentWallet()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCoinTransferFromParentWallet(params, ctx)
}

// 코인 외부 지갑 전송 요청 : 특정지갑
func (o *InternalAPI) PostCoinTransferFromUserWallet(c echo.Context) error {
	ctx := base.GetContext(c).(*context.PointManagerContext)

	params := context.NewReqCoinTransferFromUserWallet()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCoinTransferFromUserWallet(params, ctx)
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

func (o *InternalAPI) GetCoinFee(c echo.Context) error {
	params := context.NewReqCoinFee()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetCoinFee(params, c)
}
