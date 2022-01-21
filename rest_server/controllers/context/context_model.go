package context

// point db info
type PointDB struct {
	DatabaseID   int64
	DatabaseName string
	ServerName   string
}

// me app/point info
type AccountPoint struct {
	AppId                int64  `json:"app_id"`
	PointId              int64  `json:"point_id"`
	TodayLimitedQuantity int64  `json:"today_limited_quantity"`
	TodayAcqQuantity     int64  `json:"today_acq_quantity"`
	TodayCnsmQuantity    int64  `json:"today_cnsm_quantity"`
	ResetDate            string `json:"reset_date"`
}
