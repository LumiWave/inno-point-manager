package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

///////// member app coin swap 요청
type SwapPoint struct {
	MUID                  int64 `json:"mu_id"`
	AppID                 int64 `json:"app_id"` // 요청 인자
	DatabaseID            int64 `json:"database_id"`
	PointID               int64 `json:"point_id"` // 요청 인자
	PreviousPointQuantity int64 `json:"previous_point_quantity"`
	AdjustPointQuantity   int64 `json:"adjust_point_quantity"` // 요청 인자
	PointQuantity         int64 `json:"point_quantity"`
}

type SwapCoin struct {
	CoinID               int64   `json:"coin_id"` // 요청 인자
	CoinSymbol           string  `json:"coin_symbol"`
	BaseCoinID           int64   `json:"base_coin_id"` // 요청 인자
	BaseCoinSymbol       string  `json:"base_coin_symbol"`
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

	SwapFee float64 `json:"swap_fee"` // point->coin 시 전환시 부모지갑에 전송될 코인량 coin->point는 0 고정
}

func NewReqSwapInfo() *ReqSwapInfo {
	return new(ReqSwapInfo)
}

func (o *ReqSwapInfo) CheckValidate() *base.BaseResponse {

	return nil
}

////////////////////////////////////////
