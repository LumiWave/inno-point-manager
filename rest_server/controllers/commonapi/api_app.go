package commonapi

import (
	"net/http"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

// 맴버 포인트 리스트 정보 조회
func GetPointAppList(req *context.ReqGetPointApp, ctx *context.PointManagerContext) error {
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

	if pointInfo, err := inner.UpdateAppPoint(req, ctx.GetValue().AppID); err != nil {
		if strings.EqualFold(resultcode.ResultCodeText[resultcode.Result_Error_NotEqual_PreviousQuantity], err.Error()) {
			resp.SetReturn(resultcode.Result_Error_NotEqual_PreviousQuantity)
		} else if strings.EqualFold(resultcode.ResultCodeText[resultcode.Result_Error_Exceeded_DailyPoints_earned], err.Error()) {
			resp.SetReturn(resultcode.Result_Error_Exceeded_DailyPoints_earned)
		} else {
			model.MakeDbError(resp, resultcode.Result_DBError, err)
		}
	} else {
		pointInfos := context.ResPointAppUpdate{
			MUID:          req.MUID,
			PointID:       pointInfo.PointID,
			PreQuantity:   pointInfo.Quantity,
			DailyQuantity: pointInfo.DailyQuantity,
		}
		resp.Value = pointInfos
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
