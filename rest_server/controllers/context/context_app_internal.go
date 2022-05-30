package context

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

///////// app 포인트 처리 모니터링
type ReqPointAppMonitoring struct {
	MUID int64 `query:"mu_id"`
}

func NewReqPointAppMonitoring() *ReqPointAppMonitoring {
	return new(ReqPointAppMonitoring)
}

func (o *ReqPointAppMonitoring) CheckValidate() *base.BaseResponse {
	return nil
}

///////// Me Coin List
type ReqMeCoin struct {
	AUID int64 `json:"au_id" query:"au_id"`
}

func NewReqMeCoin() *ReqMeCoin {
	return new(ReqMeCoin)
}

func (o *ReqMeCoin) CheckValidate() *base.BaseResponse {
	return nil
}

type MeCoin struct {
	CoinID                    int64     `json:"coin_id" query:"coin_id"`
	BaseCoinID                int64     `json:"base_coin_id" query:"base_coin_id"`
	CoinSymbol                string    `json:"coin_symbol" query:"coin_symbol"`
	WalletAddress             string    `json:"wallet_address" query:"wallet_address"`
	Quantity                  float64   `json:"quantity" query:"quantity"`
	TodayAcqQuantity          float64   `json:"today_acq_quantity" query:"today_acq_quantity"`
	TodayCnsmQuantity         float64   `json:"today_cnsm_quantity" query:"today_cnsm_quantity"`
	TodayAcqExchangeQuantity  float64   `json:"today_acq_exchange_quantity" query:"today_acq_exchange_quantity"`
	TodayCnsmExchangeQuantity float64   `json:"today_cnsm_exchange_quantity" query:"today_cnsm_exchange_quantity"`
	ResetDate                 time.Time `json:"reset_date" query:"reset_date"`
}

////////////////////////////////////////

///////// Me Coin update
type ReqUpdateMeCoin struct {
	AUID          int64   `json:"au_id"`
	CoinID        int64   `json:"coin_id"`
	BaseCoinID    int64   `json:"base_coin_id"`
	WalletAddress string  `json:"wallet_address"`
	PreQuantity   float64 `json:"previous_quantity"`
	AdjQuantity   float64 `json:"adjust_quantity"`
	Quantity      float64 `json:"quantity"`

	//sql.Named("LogID", logID),
	//sql.Named("EventID", eventID),
}

func NewReqUpdateMeCoin() *ReqUpdateMeCoin {
	return new(ReqUpdateMeCoin)
}

func (o *ReqUpdateMeCoin) CheckValidate() *base.BaseResponse {
	return nil
}

////////////////////////////////////////
