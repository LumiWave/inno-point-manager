package context

const (
	PointType_EarnPoint     = 0
	PointType_ExchangePoint = 1
)

type DailyAppPoint struct {
	AppId                  int64 `json:"app_id"`
	PointType              int   `json:"point_type"`
	PointId                int64 `json:"point_id"`
	AdjustQuantity         int64 `json:"adjust_quantity"`
	AdjustExchangeQuantity int64 `json:"adjust_exchange_quantity"`
}
