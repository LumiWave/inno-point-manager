package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

func PostPointCoinSwap(params *context.ReqSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 보유수량이 전환에 적합한지 확인
	// 전환 비율 계산 후 타당성 확인

	// 보유 포인트 수량 가져오기
	// if memberInfo, err := point_manager_server.GetInstance().GetPointAppList(member.MUID, member.DatabaseID); err == nil {
	// 	for _, point := range memberInfo.Points {
	// 		if point.PointID == swapInfo.PointID {
	// 			swapInfo.PreviousPointQuantity = point.Quantity
	// 			swapInfo.PointQuantity = swapInfo.PreviousPointQuantity + swapInfo.AdjustPointQuantity
	// 			break
	// 		}
	// 	}

	// 	// 코인으로 전환시 0이 아닌 수량을 가지고 있는지 확인
	// 	if swapInfo.PreviousCoinQuantity == 0 {
	// 		// 내 포인트를 찾지 못함
	// 		log.Errorf("not find me point id [point_id:%v]", swapInfo.PointID)
	// 		resp.SetReturn(resultcode.Result_Invalid_PointID_Error)
	// 		return ctx.EchoContext.JSON(http.StatusOK, resp)
	// 	}
	// } else {
	// 	// point manager server 호출 에러
	// 	log.Errorf("point_manager_server GetPointAppList error : %v", err)
	// 	return ctx.EchoContext.JSON(http.StatusOK, base.BaseResponseInternalServerError())
	// }

	// // 코인으로 전환은 point의 최소 수량만 확인
	// if appPointInfo, ok := model.GetDB().AppPointsMap[swapInfo.AppID]; ok {
	// 	for _, point := range appPointInfo.Points {
	// 		if point.PointId == swapInfo.PointID {
	// 			if swapInfo.PreviousPointQuantity < point.MinExchangeQuantity || swapInfo.PreviousPointQuantity < swapInfo.AdjustPointQuantity{
	// 				resp.SetReturn(resultcode.Result_Invalid_PointQuantity_Error)
	// 				return ctx.EchoContext.JSON(http.StatusOK, resp)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	// 존재하지 않는 appid를 전환 시도 에러
	// 	resp.SetReturn(resultcode.Result_Invalid_AppID_Error)
	// 	return ctx.EchoContext.JSON(http.StatusOK, resp)
	// }

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
