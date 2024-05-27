package api_inno_log

type Common struct {
	Return  int    `json:"return"`
	Message string `json:"message"`
}

type AccountAuthLog struct {
	LogDt      string `json:"log_dt"`
	LogID      int64  `json:"log_id"`
	EventID    int64  `json:"event_id"`
	AUID       int64  `json:"au_id"`
	InnoUID    string `json:"inno_id"`
	SocialID   string `json:"social_id"`
	SocialType int64  `json:"social_type"`
}

type AccountCoinLog struct {
	LogDt         string `json:"log_dt"`
	LogID         int64  `json:"log_id"`
	EventID       int64  `json:"event_id"`
	TxHash        string `json:"tx_hash"`
	AUID          int64  `json:"au_id"`
	CoinID        int64  `json:"coin_id"`
	BaseCoinID    int64  `json:"basecoin_id"`
	WalletAddress string `json:"wallet_address"`
	WalletTypeID  int64  `json:"wallet_type_id"`
	AdjQuantity   string `json:"adjust_quantity"`
	WalletID      int64  `json:"wallet_id"`
}

type ExchangeGoodsLog struct {
	LogDt            string `json:"log_dt"`
	LogID            int64  `json:"log_id"`
	EventID          int64  `json:"event_id"`
	TxHash           string `json:"tx_hash"`
	TxID             int64  `json:"tx_id"`
	AUID             int64  `json:"au_id"`
	InnoUID          string `json:"inno_uid"`
	MUID             int64  `json:"mu_id"`
	AppID            int64  `json:"app_id"`
	CoinID           int64  `json:"coin_id"`
	BaseCoinID       int64  `json:"basecoin_id"`
	WalletAddress    string `json:"wallet_address"`
	WalletTypeID     int64  `json:"wallet_type_id"`
	AdjCoinQuantity  string `json:"adjust_coin_quantity"`
	PointID          int64  `json:"point_id"`
	AdjPointQuantity int64  `json:"adjust_point_quantity"`
	WalletID         int64  `json:"wallet_id"`
}

type MemberAuthLog struct {
	LogDt      string `json:"log_dt"`
	LogID      int64  `json:"log_id"`
	EventID    int64  `json:"event_id"`
	AUID       int64  `json:"au_id"`
	InnoUID    string `json:"inno_id"`
	MUID       int64  `json:"mu_id"`
	AppID      int64  `json:"app_id"`
	DataBaseID int64  `json:"database_id"`
}

type MemberPointsLog struct {
	LogDt       string `json:"log_dt"`
	LogID       int64  `json:"log_id"`
	EventID     int64  `json:"event_id"`
	AUID        int64  `json:"au_id"`
	MUID        int64  `json:"mu_id"`
	AppID       int64  `json:"app_id"`
	PointID     int64  `json:"point_id"`
	AdjQuantity int64  `json:"adjust_quantity"`
}
