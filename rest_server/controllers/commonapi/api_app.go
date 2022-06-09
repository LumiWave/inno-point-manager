package commonapi

import (
	"net/http"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	"github.com/labstack/echo"
)

// 맴버 포인트 리스트 정보 조회
func GetPointAppList(req *context.ReqGetPointApp, ctx *context.PointManagerContext) error {
	log.Debugf("GetPointAppList [AppID:%v][MUID:%v][PointID:%v]", req.AppId, req.MUID, req.PointID)

	resp := new(base.BaseResponse)
	resp.Success()

	// 포인트 정보 조회
	if pointInfo, err := inner.LoadPointList(req.MUID, req.DatabaseID, ctx.GetValue().AppID); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		pointInfos := context.ResPointMemberRegister{
			PointInfo: *pointInfo,
		}
		resp.Value = pointInfos
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 맴버 포인트 리스트 정보 조회
func GetPointApp(req *context.ReqGetPointApp, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 포인트 정보 조회
	if pointInfo, err := inner.LoadPoint(req.MUID, req.PointID, req.DatabaseID, ctx.GetValue().AppID); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		pointInfos := context.ResPointMemberRegister{
			PointInfo: *pointInfo,
		}
		resp.Value = pointInfos
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// app point 업데이트
func PutPointAppUpdate(req *context.ReqPointAppUpdate, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if !model.GetPointUpdateEnable() {
		resp.SetReturn(resultcode.Result_Error_IsSwapMaintenance)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	if pointInfo, err := inner.UpdateAppPoint(req, ctx.GetValue().AppID); err != nil {
		if strings.EqualFold(resultcode.ResultCodeText[resultcode.Result_Error_NotEqual_PreviousQuantity], err.Error()) {
			resp.SetReturn(resultcode.Result_Error_NotEqual_PreviousQuantity)
		} else if strings.EqualFold(resultcode.ResultCodeText[resultcode.Result_Error_Exceeded_TodayPoints_earned], err.Error()) {
			resp.SetReturn(resultcode.Result_Error_Exceeded_TodayPoints_earned)
		} else {
			model.MakeDbError(resp, resultcode.Result_DBError, err)
		}
	} else {
		pointInfos := context.ResPointAppUpdate{
			MUID:          req.MUID,
			PointID:       pointInfo.PointID,
			PreQuantity:   pointInfo.Quantity,
			TodayQuantity: pointInfo.TodayQuantity,
		}
		resp.Value = pointInfos
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// member 코인 리스트 조회
func GetMeCoinList(c echo.Context, reqMeCoin *context.ReqMeCoin) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if coinList, _, err := model.GetDB().GetAccountCoins(reqMeCoin.AUID); err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_Get_Me_CoinList_Error)
	} else {
		resp.Value = []*context.MeCoin{}
		if coinList != nil {
			resp.Value = coinList
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// member 코인 업데이트
func PutMeCoin(c echo.Context, reqMeCoin *context.ReqUpdateMeCoin) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// Lockkey := model.MakeCoinTransferFromUserWalletLockKey(reqMeCoin.AUID)
	// mutex := model.GetDB().RedSync.NewMutex(Lockkey)
	// isValid, _ := mutex.Valid()
	// if isValid {
	// 	log.Errorf("auid:%v %v", reqMeCoin.AUID, resultcode.ResultCodeText[resultcode.Result_RedisError_WaitForProcessing])
	// 	resp.SetReturn(resultcode.Result_RedisError_WaitForProcessing)
	// 	return c.JSON(http.StatusOK, resp)
	// }
	// if err := mutex.Lock(); err != nil {
	// 	log.Error("redis lock err:%v", err)
	// 	resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
	// 	return c.JSON(http.StatusOK, resp)
	// }

	// defer func() {
	// 	// 1-1. redis unlock
	// 	if ok, err := mutex.Unlock(); !ok || err != nil {
	// 		if err != nil {
	// 			log.Errorf("unlock err : %v", err)
	// 		}
	// 	}
	// }()

	if err := model.GetDB().UpdateAccountCoins(
		reqMeCoin.AUID,
		reqMeCoin.CoinID,
		reqMeCoin.BaseCoinID,
		reqMeCoin.WalletAddress,
		reqMeCoin.PreQuantity,
		reqMeCoin.AdjQuantity, // 전송 수수료 + amount
		reqMeCoin.Quantity,
		context.LogID_type(reqMeCoin.LogID),
		context.EventID_type(reqMeCoin.EventID),
		""); err != nil {
		log.Errorf("UpdateAccountCoins error : %v, auid:%v, pre:%v, adj:%v, quantity:%v",
			err, reqMeCoin.AUID, reqMeCoin.PreQuantity, reqMeCoin.AdjQuantity, reqMeCoin.Quantity)
		resp.SetReturn(resultcode.Result_Error_DB_Update_Me_Coin_Error)
	}

	return c.JSON(http.StatusOK, resp)
}
