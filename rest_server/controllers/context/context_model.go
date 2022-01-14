package context

// point db info
type PointDB struct {
	DatabaseID   int64
	DatabaseName string
	ServerName   string
}

// me app/point info
type AccountPoint struct {
	AppId         int64  `json:"app_id"`
	PointId       int64  `json:"point_id"`
	DailyQuantity int64  `json:"daily_quantity"`
	ResetDate     string `json:"reset_date"`
}
