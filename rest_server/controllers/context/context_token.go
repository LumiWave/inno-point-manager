package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

type Swap_State_type string

const (
	Swap_State_type_ready    = "0" // swap 대기
	Swap_State_type_ing      = "1" // swap 중
	Swap_State_type_complete = "2" // swap 완료
)

///////// member private 토큰 swap 요청
type PostPointTokenSwap struct {
	ContextKey
	PrivateTokenAmount string `json:"private_token_amount"`
}

func NewPostPointTokenSwap() *PostPointTokenSwap {
	return new(PostPointTokenSwap)
}

func (o *PostPointTokenSwap) CheckValidate() *base.BaseResponse {
	if o.CpMemberIdx == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
	}
	if len(o.PrivateTokenAmount) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PrivateTokenAmount)
	}
	return nil
}

////////////////////////////////////////

///////// member private 토큰 swap history 정보
type PointMemberTokenSwapHistory struct {
	ContextKey

	LatestPrivateTokenAmount string          `json:"latest_private_token_amount"`
	SwapPrivateTokenAmount   string          `json:"change_private_token_amount"`
	LatestPublicTokenAmount  string          `json:"latest_public_token_amount"`
	SwapPublicTokenAmount    string          `json:"change_public_token_amount"`
	TxnHash                  string          `json:"txn_hash"`
	SwapState                Swap_State_type `json:"swap_state"`
	CreateAt                 int64           `json:"create_at"`
	PageInfo
}

func NewPointMemberTokenSwapHistory() *PointMemberTokenSwapHistory {
	return new(PointMemberTokenSwapHistory)
}

func (o *PointMemberTokenSwapHistory) CheckValidate() *base.BaseResponse {
	if o.CpMemberIdx == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
	}
	if len(o.PageOffset) == 0 || len(o.PageSize) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PageInfo)
	}
	return nil
}

type PointMemberTokenSwapHistoryResponse struct {
	PageInfo PageInfoResponse              `json:"page_info"`
	Historys []PointMemberTokenSwapHistory `json:"historys"`
}

////////////////////////////////////////
