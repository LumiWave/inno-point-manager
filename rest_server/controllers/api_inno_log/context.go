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

type ExchangeLogs struct {
	LogDT             string `json:"log_dt"`
	EventID           int64  `json:"event_id"`
	TxID              int64  `json:"tx_id"`
	AUID              int64  `json:"au_id"`
	MUID              int64  `json:"mu_id"`
	InnoUID           string `json:"inno_uid"`
	AppID             int64  `json:"app_id"`
	ExchangeFees      string `json:"exchange_fees"`
	FromBaseCoinID    int64  `json:"from_base_coin_id"`
	FromWalletTypeID  int64  `json:"from_wallet_type_id"`
	FromWalletID      int64  `json:"from_wallet_id"`
	FromWalletAddress string `json:"from_wallet_address"`
	FromID            int64  `json:"from_id"`
	FromAdjQuantity   string `json:"from_adjust_quantity"`
	ToBaseCoinID      int64  `json:"to_base_coin_id"`
	ToWalletTypeID    int64  `json:"to_wallet_type_id"`
	ToWalletID        int64  `json:"to_wallet_id"`
	ToWalletAddress   string `json:"to_wallet_address"`
	ToID              int64  `json:"to_id"`
	ToAdjQuantity     string `json:"to_adjust_quantity"`
	TransactedDT      string `json:"transacted_dt"`
	CompletedDT       string `json:"completed_dt"`
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
