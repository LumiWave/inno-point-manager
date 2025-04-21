package context

type PreSalesExchange struct {
	SalesID                 int64   `json:"sales_id"`
	BaseCoinID              int64   `json:"base_coin_id"`
	CoinID                  int64   `json:"coin_id"`
	PointID                 int64   `json:"point_id"`
	MinimumExchangeQuantity string  `json:"minimum_exchange_quantity"`
	ExchangeRatio           float64 `json:"exchange_ratio"`
}
