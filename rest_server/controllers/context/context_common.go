package context

import "github.com/ONBUFF-IP-TOKEN/baseutil/datetime"

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

type ContextKey struct {
	Idx         int64 `json:"idx" query:"idx"`
	CpMemberIdx int64 `json:"cp_member_idx" query:"cp_member_idx"`
}

// page info
type PageInfo struct {
	PageOffset string `json:"page_offset,omitempty" query:"page_offset" validate:"required"`
	PageSize   string `json:"page_size,omitempty" query:"page_size" validate:"required"`
}

// page response
type PageInfoResponse struct {
	PageInfo
	TotalSize string `json:"total_size"`
}

func MakeAt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}
