package token_manager_server

type Common struct {
	Return  int    `json:"return"`
	Message string `json:"message"`
}

// /////// 부모 지갑에서 전송 요청
type ReqSendFromParentWallet struct {
	BaseSymbol string `json:"base_symbol"`
	Contract   string `json:"contract"`
	ToAddress  string `json:"to_address"`
	Amount     string `json:"amount"`
	Memo       string `json:"memo"`
}

type ResSendFromParentWalletValue struct {
	IsSuccess bool   `json:"is_success"`
	TxHash    string `json:"tx_hash"`
}

type ResSendFromParentWallet struct {
	Common
	Value ResSendFromParentWalletValue `json:"value"`
}

////////////////////////////////////////

// /////// 특정 지갑에서 전송 요청
type ReqSendFromUserWallet struct {
	Contract       string `json:"contract"`
	BaseCoinSymbol string `json:"base_symbol"`
	FromAddress    string `json:"from_address"`
	ToAddress      string `json:"to_address"`
	Amount         string `json:"amount"`
	Memo           string `json:"memo"`
}

type ResSendFromUserWalletValue struct {
	IsSuccess bool   `json:"is_success"`
	TxHash    string `json:"tx_hash"`
}

type ResSendFromUserWallet struct {
	Common
	Value ResSendFromUserWalletValue `json:"value"`
}

// //////////////////////////////////////
// 잔액 조회
type ReqBalance struct {
	BaseSymbol string `url:"base_symbol"`
	Contract   string `url:"contract"`
	Address    string `url:"address"`
}

type ResReqBalanceValue struct {
	Balance string `json:"balance"`
	Address string `json:"address"`
	Decimal int64  `json:"decimal"`
}

type ResBalanc struct {
	Common
	ResReqBalanceValue `json:"value"`
}

////////////////////////////////////////

// //////// sui 코인/토큰 보유 object id 리스트 조회
type ReqCoinObjects struct {
	WalletAddress   string `url:"wallet_address"`
	ContractAddress string `url:"contract_address"`
}

type ResCoinObjects struct {
	Common
	Value struct {
		ObjectIDs []string `json:"object_ids"`
	} `json:"value"`
}

////////////////////////////////////////

// /////// 코인 가스비 요청
type ReqCoinFee struct {
	Symbol string `url:"base_symbol"`
}

type ResCoinFeeInfoValue struct {
	GasPrice string `json:"gas_price"`
	Decimal  int64  `json:"decimal"`
}

type ResCoinFeeInfo struct {
	Common
	ResCoinFeeInfoValue `json:"value"`
}

////////////////////////////////////////
