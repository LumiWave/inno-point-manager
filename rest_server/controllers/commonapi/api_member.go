package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func PostPointMemberRegister(req *context.ReqPointMemberRegister, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	log.Infof("register member [AUID:%v][MUDI:%v]", req.AUID, req.MUID)
	if err := model.GetDB().InsertPointMember(req); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		// 강제로 0 point 업데이트
		newPointInfo := context.PointInfo{
			DatabaseID: req.DatabaseID,
			MUID:       req.MUID,
		}
		for _, pointInfo := range model.GetDB().AppPointsMap[req.AppID].Points {
			newPointInfo.Points = append(newPointInfo.Points, &context.Point{
				PointID:       pointInfo.PointId,
				Quantity:      0,
				TodayQuantity: 0,
			})
		}

		pointInfos := context.ResPointMemberRegister{
			PointInfo: newPointInfo,
		}
		resp.Value = pointInfos
	}

	// log.Infof("register member [AUID:%v][MUDI:%v]", req.AUID, req.MUID)
	// if err := model.GetDB().InsertPointMember(req); err != nil {
	// 	model.MakeDbError(resp, resultcode.Result_DBError, err)
	// } else {
	// 	// 강제로 0 point 업데이트
	// 	for _, pointInfo := range model.GetDB().AppPointsMap[req.AppID].Points {
	// 		model.GetDB().InsertMemberPoints(req.DatabaseID, req.MUID, pointInfo.PointId, 0)
	// 	}
	// 	// 포인트 정보 조회
	// 	if pointInfo, err := inner.LoadPointList(req.MUID, req.DatabaseID, req.AppID); err != nil {
	// 		model.MakeDbError(resp, resultcode.Result_DBError, err)
	// 	} else {
	// 		pointInfos := context.ResPointMemberRegister{
	// 			PointInfo: *pointInfo,
	// 		}
	// 		resp.Value = pointInfos
	// 	}
	// }

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointMemberWallet(req *context.ReqPointMemberWallet, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if wallets, err := model.GetDB().GetPointMemberWallet(req, ctx.GetValue().AppID); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		resp.Value = wallets
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
