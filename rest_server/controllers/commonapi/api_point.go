package commonapi

import (
	"net/http"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func PutPointAppUpdate(params *context.ReqPointMemberAppUpdate, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// context.MakeAt(&params.CreateAt)

	// // 1. 존재하는 member 인지 확인
	// if value, err := model.GetDB().SelectPointMember(params.CpMemberIdx); err != nil {
	// 	resp.SetReturn(resultcode.Result_DBError)
	// } else {
	// 	if value.CpMemberIdx <= 0 {
	// 		// 존재하지 않는 member
	// 		resp.SetReturn(resultcode.Result_Error_NotExistMember)
	// 	} else {
	// 		// 2. latest_point_amount 정보가 point_member 테이블과 동일 한지 확인
	// 		if params.LatestPointAmount != value.PointAmount {
	// 			resp.SetReturn(resultcode.Result_Error_LatestPointAmountIsDiffrent)
	// 		} else {
	// 			// 4. point_member 테이블에 point_amount 정보 update
	// 			pA1, _ := strconv.ParseInt(value.PointAmount, 10, 64)
	// 			pA2, _ := strconv.ParseInt(params.ChangePointAmount, 10, 64)
	// 			value.PointAmount = strconv.FormatInt(pA1+pA2, 10)
	// 			if err := model.GetDB().UpdatePointMember(value); err != nil {
	// 				resp.SetReturn(resultcode.Result_DBError)
	// 			} else {
	// 				// 3. app_point_history 테이블에 정보 insert
	// 				if err := model.GetDB().InsertPointAppHistory(params); err != nil {
	// 					resp.SetReturn(resultcode.Result_DBError)
	// 				}
	// 			}
	// 		}

	// 		resp.Value = value
	// 	}
	// }

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
