package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

///////// member private 토큰 swap 요청
type SwapPoint struct {
	MUID                  int64 `json:"mu_id"`
	AppID                 int64 `json:"app_id"` // 요청 인자
	DatabaseID            int64 `json:"database_id"`
	PointID               int64 `json:"point_id"` // 요청 인자
	PreviousPointQuantity int64 `json:"previous_point_quantity"`
	AdjustPointQuantity   int64 `json:"adjust_point_quantity"` // 요청 이니자
	PointQuantity         int64 `json:"point_quantity"`
}

type SwapCoin struct {
	CoinID               int64   `json:"coin_id"` // 요청 인자
	WalletAddress        string  `json:"walletaddress"`
	PreviousCoinQuantity float64 `json:"previous_coin_quantity"`
	AdjustCoinQuantity   float64 `json:"adjust_coin_quantity"` // 요청 인자
	CoinQuantity         float64 `json:"coin_quantity"`
}

type ReqSwapInfo struct {
	AUID int64 `json:"au_id"`

	SwapPoint `json:"point"`
	SwapCoin  `json:"coin"`

	LogID   int64 `json:"log_id"`   // 2: 전환
	EventID int64 `json:"event_id"` // 3: point->coin,  4: coin->point
}

func NewReqSwapInfo() *ReqSwapInfo {
	return new(ReqSwapInfo)
}

func (o *ReqSwapInfo) CheckValidate() *base.BaseResponse {

	return nil
}

////////////////////////////////////////
