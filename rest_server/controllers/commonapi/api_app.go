package commonapi

import (
	"net/http"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

// 맴버 포인트 정보 조회
func GetPointApp(req *context.ReqGetPointApp, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 포인트 정보 조회
	if pointInfo, err := inner.LoadPoint(req.CUID, req.AppID, req.DatabaseID); err != nil {
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

	if pointInfo, err := inner.UpdateAppPoint(req); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	} else {
		pointInfos := context.ResPointAppUpdate{
			CUID:         req.CUID,
			AppID:        req.AppID,
			PointID:      pointInfo.PointID,
			LastQuantity: pointInfo.Quantity,
		}
		resp.Value = pointInfos
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointAppHistory(params *context.PointMemberHistoryReq, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if value, totalSize, err := model.GetDB().SelectPointAppHistory(params); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		PageInfo := context.PageInfoResponse{
			PageInfo: context.PageInfo{
				PageOffset: params.PageOffset,
				PageSize:   strconv.FormatInt(int64(len(*value)), 10),
			},
			TotalSize: strconv.FormatInt(int64(totalSize), 10),
		}
		res := context.PointMemberHistoryResponse{
			PageInfo: PageInfo,
			Historys: *value,
		}

		resp.Value = res
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostPointAppExchange(params *context.PostPointAppExchange, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// // 1. app point 잔량 check
	// if value, err := model.GetDB().SelectPointMember(params.CpMemberIdx); err != nil {
	// 	resp.SetReturn(resultcode.Result_DBError)
	// } else {
	// 	if value.CpMemberIdx <= 0 {
	// 		// 존재하지 않는 member
	// 		resp.SetReturn(resultcode.Result_Error_NotExistMember)
	// 	} else {
	// 		// 보유 포인트 수량보다 많은 전환 량을 요청시 에러
	// 		if util.CompareString(value.PointAmount, params.PointAmount) < 1 {
	// 			resp.SetReturn(resultcode.Result_Require_ValidPointAmount)
	// 		} else {
	// 			// 2. private token 전환 요청 = 임시로 성공 처리  todo
	// 			{
	// 				// block chain 네트워크에 전환 요청
	// 				// 전환 성공여부 확인
	// 				// private token 보유량 확인
	// 			}

	// 			// 3. point_exchange_history 테이블에 정보 추가
	// 			history := context.PointMemberExchangeHistory{
	// 				ContextKey: context.ContextKey{
	// 					Idx:         value.Idx,
	// 					CpMemberIdx: value.CpMemberIdx,
	// 				},
	// 				LatestPointAmount:          value.PointAmount,
	// 				ExchangePointAmount:        params.PointAmount,
	// 				LatestPrivateTokenAmount:   value.PrivateTokenAmount,
	// 				ExchangePrivateTokenAmount: params.PointAmount,
	// 				TxnHash:                    "0xtest_hash",
	// 				ExchangeState:              context.Exchange_State_type_complete,
	// 			}
	// 			context.MakeAt(&history.CreateAt)

	// 			if err := model.GetDB().InsertPointAppExchangeHistory(&history); err != nil {
	// 				resp.SetReturn(resultcode.Result_DBError)
	// 			} else {
	// 				// 4. point_member 테이블에 최종 정보 갱신
	// 				value.PointAmount = strconv.FormatInt(util.ParseInt(value.PointAmount)-util.ParseInt(params.PointAmount), 10)
	// 				value.PrivateTokenAmount = strconv.FormatInt(util.ParseInt(value.PrivateTokenAmount)+util.ParseInt(params.PointAmount), 10)
	// 				if err := model.GetDB().UpdatePointMember(value); err != nil {
	// 					resp.SetReturn(resultcode.Result_DBError)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointAppExchangeHistory(params *context.PointMemberExchangeHistory, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if value, totalSize, err := model.GetDB().SelectPointAppExchangeHistory(params); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		PageInfo := context.PageInfoResponse{
			PageInfo: context.PageInfo{
				PageOffset: params.PageOffset,
				PageSize:   strconv.FormatInt(int64(len(*value)), 10),
			},
			TotalSize: strconv.FormatInt(int64(totalSize), 10),
		}
		res := context.PointMemberExchangeHistoryResponse{
			PageInfo: PageInfo,
			Historys: *value,
		}

		resp.Value = res
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
