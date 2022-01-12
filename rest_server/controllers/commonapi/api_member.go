package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func PostPointMemberRegister(req *context.ReqPointMemberRegister, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := model.GetDB().InsertPointMember(req); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		// 강제로 0 point 업데이트
		for _, pointInfo := range model.GetDB().AppPointsMap[req.AppID].Points {
			model.GetDB().UpdateAppPoint(req.MUID, pointInfo.PointId, 0, req.DatabaseID)
		}
		// 포인트 정보 조회
		if pointInfo, err := inner.LoadPointList(req.MUID, req.DatabaseID); err != nil {
			model.MakeDbError(resp, resultcode.Result_DBError, err)
		} else {
			pointInfos := context.ResPointMemberRegister{
				PointInfo: *pointInfo,
			}
			resp.Value = pointInfos
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointMemberWallet(req *context.ReqPointMemberWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if wallets, err := model.GetDB().GetPointMemberWallet(req, ctx.VerifyValue.AppID); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		resp.Value = wallets
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
