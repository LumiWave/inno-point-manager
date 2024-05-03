package context

import "github.com/LumiWave/baseutil/datetime"

type LogID_type int

const (
	LogID_cp              = 1 // 고객사
	LogID_exchange        = 2 // 전환
	LogID_external_wallet = 3 // 외부지갑

	LogID_wallet_sync = 6
)

type EventID_type int

const (
	EventID_add     = 1 // 재화 증가
	EventID_sub     = 2 // 재화 감소
	EventID_toCoin  = 3 // 포인트->코인
	EventID_toPoint = 4 // 코인->포인트

	EventID_normal_purchase        = 9  // 일반 상품 구매
	EventID_auction_deposit        = 10 // 경매 입찰 보증금 납부
	EventID_auction_deposit_refund = 11 // 경매 입찰 보증금 환급
	EventID_auction_purchase       = 12 // 경매 상품 구매

	EventID_add_fee = 24 // 수수료 재화 증가 : swap시 수수료 받았을때
	EventID_sub_fee = 25 // 수수료 재화 감소 : 수수료 명목으로 서비스 제공할때
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
