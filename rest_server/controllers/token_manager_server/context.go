package token_manager_server

type Common struct {
	Return  int    `json:"return"`
	Message string `json:"message"`
}

///////// 부모 지갑에서 전송 요청
type ReqSendFromParentWallet struct {
	Symbol    string `json:"symbol"`
	ToAddress string `json:"to_address"`
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
}

type ResSendFromParentWalletValue struct {
	ReqId         string `json:"reqid"`
	TransactionId string `json:"transaction_id"`
}

type ResSendFromParentWallet struct {
	Common
	Value ResSendFromParentWalletValue `json:"value"`
}

////////////////////////////////////////

///////// 특정 지갑에서 전송 요청
type ReqSendFromUserWallet struct {
	Symbol         string `json:"symbol"`
	BaseCoinSymbol string `json:"base_coin_symbol"`
	FromAddress    string `json:"from_address"`
	ToAddress      string `json:"to_address"`
	Amount         string `json:"amount"`
	Memo           string `json:"memo"`
}

type ResSendFromUserWalletValue struct {
	TransactionHash string `json:"transaction_hash"`
}

type ResSendFromUserWallet struct {
	Common
	Value ResSendFromUserWalletValue `json:"value"`
}

////////////////////////////////////////

///////// 코인 가스비 요청
type ReqCoinFee struct {
	Symbol string `query:"symbol"`
}

type ResCoinFeeInfoValue struct {
	Fast    string `json:"fast"`
	Slow    string `json:"slow"`
	Average string `json:"average"`
	BaseFee string `json:"basefee"`
	Fastest string `json:"fastest"`
}

type ResCoinFeeInfo struct {
	Common
	ResCoinFeeInfoValue `json:"value"`
}

////////////////////////////////////////
