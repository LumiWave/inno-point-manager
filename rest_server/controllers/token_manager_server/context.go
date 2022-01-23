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
