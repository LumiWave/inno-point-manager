package token_manager_server

type Common struct {
	Return  int    `json:"return"`
	Message string `json:"message"`
}

// /////// 부모 지갑에서 전송 요청
type ReqSendFromParentWallet struct {
	BaseSymbol string `json:"base_symbol"`
	Symbol     string `json:"symbol"`
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
	Symbol         string `json:"symbol"`
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
	Symbol  string `query:"symbol"`
	Address string `query:"address"`
}

type ResReqBalanceValue struct {
	Balance string `json:"balance"`
	Address string `json:"address"`
}

type ResBalanc struct {
	Common
	ResReqBalanceValue `json:"value"`
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
