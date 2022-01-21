package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

type LogID_type int

const (
	LogID_cp              = 1 // 고객사
	LogID_exchange        = 2 // 전환
	LogID_external_wallet = 3 // 외부지갑
)

type EventID_type int

const (
	EventID_add     = 1 // 재화 증가
	EventID_sub     = 2 // 재화 감소
	EventID_toCoin  = 3 // 포인트->코인
	EventID_toPoint = 4 // 코인->포인트
)

type Exchange_State_type string

const (
	Exchange_State_type_ready    = "0" // 전환 대기
	Exchange_State_type_ing      = "1" // 전환 중
	Exchange_State_type_complete = "2" // 전환 완료
)

type Point struct {
	PointID       int64  `json:"point_id"`
	Quantity      int64  `json:"quantity"`
	TodayQuantity int64  `json:"today_quantity"`
	ResetDate     string `json:"reset_date"`

	PreQuantity    int64 `json:"previous_quantity,omitempty"`
	AdjustQuantity int64 `json:"adjust_quantity,omitempty"`
}

type PointInfo struct {
	MyUuid     string `json:"my_uuid,omitempty"`
	DatabaseID int64  `json:"database_id"`

	MUID   int64    `json:"mu_id"`
	Points []*Point `json:"points"`
}

///////// member 포인트 조회
type ReqGetPointApp struct {
	AppId      int64 `query:"app_id"`
	MUID       int64 `query:"mu_id"`
	PointID    int64 `query:"point_id"`
	DatabaseID int64 `query:"database_id"`
}

func NewReqGetPointApp() *ReqGetPointApp {
	return new(ReqGetPointApp)
}

func (o *ReqGetPointApp) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.MUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MUID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}
	if o.AppId != 0 && ctx.GetValue() != nil {
		ctx.GetValue().AppID = o.AppId
	}
	return nil
}

////////////////////////////////////////

///////// app 포인트 업데이트
type ReqPointAppUpdate struct {
	MUID       int64 `json:"mu_id"`
	PointID    int64 `json:"point_id"`
	DatabaseID int64 `json:"database_id"`

	PreQuantity    int64 `json:"previous_quantity"`
	AdjustQuantity int64 `json:"adjust_quantity"`
}

func NewReqPointMemberAppUpdate() *ReqPointAppUpdate {
	return new(ReqPointAppUpdate)
}

func (o *ReqPointAppUpdate) CheckValidate() *base.BaseResponse {
	if o.MUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MUID)
	}
	if o.PointID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PointID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}
	if o.AdjustQuantity == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
	}
	return nil
}

type ResPointAppUpdate struct {
	MUID    int64 `json:"mu_id"`
	PointID int64 `json:"point_id"`

	PreQuantity   int64 `json:"previous_quantity"`
	TodayQuantity int64 `json:"today_quantity"`
}

////////////////////////////////////////
