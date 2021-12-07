package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

type Point_type string

const (
	Point_type_add = "0" // 추가
	Point_type_sub = "1" // 감소
)

type Exchange_State_type string

const (
	Exchange_State_type_ready    = "0" // 전환 대기
	Exchange_State_type_ing      = "1" // 전환 중
	Exchange_State_type_complete = "2" // 전환 완료
)

///////// member 포인트 누적 history 정보
type PointMemberHistory struct {
	ContextKey

	Type              Point_type `json:"type"`
	LatestPointAmount string     `json:"latest_point_amount"`
	ChangePointAmount string     `json:"change_point_amount"`

	CreateAt int64 `json:"create_at"`
	//ExchangeAt int64 `json:"exchange_at"`

	PageInfo
}

////////////////////////////////////////

///////// member 포인트 업데이트
type ReqPointMemberAppUpdate struct {
	CUID       string `json:"cu_id"`
	AppID      int64  `json:"app_id"`
	PointID    int64  `json:"point_id"`
	DatabaseID int64  `json:"database_id"`

	LastQuantity   int64 `json:"last_quantity"`
	ChangeQuantity int64 `json:"change_quantity"`
}

func NewReqPointMemberAppUpdate() *ReqPointMemberAppUpdate {
	return new(ReqPointMemberAppUpdate)
}

func (o *ReqPointMemberAppUpdate) CheckValidate() *base.BaseResponse {
	if len(o.CUID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_CUID)
	}
	if o.AppID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AppID)
	}
	if o.PointID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PointID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}
	if o.ChangeQuantity == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_ChangeQuantity)
	}
	return nil
}

type ResPointMemberAppUpdate struct {
	CUID    string `json:"cu_id"`
	AppID   int64  `json:"app_id"`
	PointID int64  `json:"point_id"`

	LastQuantity int64 `json:"last_quantity"`
}

////////////////////////////////////////

///////// member 포인트 누적 history 요청
type PointMemberHistoryReq struct {
	ContextKey
	PageInfo
}

func NewPointMemberHistoryReq() *PointMemberHistoryReq {
	return new(PointMemberHistoryReq)
}

func (o *PointMemberHistoryReq) CheckValidate() *base.BaseResponse {
	if o.CpMemberIdx == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
	}
	if len(o.PageOffset) == 0 || len(o.PageSize) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PageInfo)
	}
	return nil
}

type PointMemberHistoryResponse struct {
	PageInfo PageInfoResponse     `json:"page_info"`
	Historys []PointMemberHistory `json:"historys"`
}

////////////////////////////////////////

///////// member 포인트 private 토큰 전환 요청
type PostPointAppExchange struct {
	ContextKey
	PointAmount string `json:"point_amount"`
}

func NewPostPointAppExchange() *PostPointAppExchange {
	return new(PostPointAppExchange)
}

func (o *PostPointAppExchange) CheckValidate() *base.BaseResponse {
	if o.CpMemberIdx == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
	}
	if len(o.PointAmount) == 0 {
		o.PointAmount = "0"
	}
	return nil
}

////////////////////////////////////////

///////// member 포인트 private 토큰 전환 history 정보
type PointMemberExchangeHistory struct {
	ContextKey

	LatestPointAmount          string              `json:"latest_point_amount"`
	ExchangePointAmount        string              `json:"change_point_amount"`
	LatestPrivateTokenAmount   string              `json:"latest_private_token_amount"`
	ExchangePrivateTokenAmount string              `json:"change_private_token_amount"`
	TxnHash                    string              `json:"txn_hash"`
	ExchangeState              Exchange_State_type `json:"exchange_state"`
	CreateAt                   int64               `json:"create_at"`

	PageInfo
}

func NewPointMemberExchangeHistory() *PointMemberExchangeHistory {
	return new(PointMemberExchangeHistory)
}

func (o *PointMemberExchangeHistory) CheckValidate(bGet bool) *base.BaseResponse {
	if o.CpMemberIdx == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
	}
	if bGet {
		if len(o.PageOffset) == 0 || len(o.PageSize) == 0 {
			return base.MakeBaseResponse(resultcode.Result_Require_PageInfo)
		}
	}
	return nil
}

type PointMemberExchangeHistoryResponse struct {
	PageInfo PageInfoResponse             `json:"page_info"`
	Historys []PointMemberExchangeHistory `json:"historys"`
}

////////////////////////////////////////
