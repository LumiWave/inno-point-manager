package commonapi

import (
	"net/http"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/util"
)

func PostPointTokenSwap(params *context.PostPointTokenSwap, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. private token 잔여량 체크
	if value, err := model.GetDB().SelectPointMember(params.CpMemberIdx); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if util.CompareString(value.PrivateTokenAmount, params.PrivateTokenAmount) < 0 {
			resp.SetReturn(resultcode.Result_Error_LackOfTokenQuantity)
		} else {
			// 2. public token swap todo

			// 3. swap 기록 저장
			swap := &context.PointMemberTokenSwapHistory{
				ContextKey: context.ContextKey{
					CpMemberIdx: params.CpMemberIdx,
				},
				LatestPrivateTokenAmount: value.PrivateTokenAmount,
				SwapPrivateTokenAmount:   params.PrivateTokenAmount,
				LatestPublicTokenAmount:  util.SumString(value.PublicTokenAmount, params.PrivateTokenAmount),
				SwapPublicTokenAmount:    params.PrivateTokenAmount,
				TxnHash:                  "0xswap_hash",
				SwapState:                context.Swap_State_type_complete,
				CreateAt:                 datetime.GetTS2MilliSec(),
			}

			if err := model.GetDB().InsertPointTokenSwapHistory(swap); err != nil {
				resp.SetReturn(resultcode.Result_DBError)
			} else {
				// 4. point_member 테이블에 최종 정보 갱신
				value.PrivateTokenAmount = util.SubString(value.PrivateTokenAmount, params.PrivateTokenAmount)
				value.PublicTokenAmount = util.SumString(value.PublicTokenAmount, params.PrivateTokenAmount)
				if err := model.GetDB().UpdatePointMember(value); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
				}
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointTokenSwapHistory(params *context.PointMemberTokenSwapHistory, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if value, totalSize, err := model.GetDB().SelectPointTokenSwapHistory(params); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		PageInfo := context.PageInfoResponse{
			PageInfo: context.PageInfo{
				PageOffset: params.PageOffset,
				PageSize:   strconv.FormatInt(int64(len(*value)), 10),
			},
			TotalSize: strconv.FormatInt(int64(totalSize), 10),
		}
		res := context.PointMemberTokenSwapHistoryResponse{
			PageInfo: PageInfo,
			Historys: *value,
		}

		resp.Value = res
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
