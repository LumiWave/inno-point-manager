package context

import "github.com/ONBUFF-IP-TOKEN/baseapp/base"

///////// 코인 외부 지갑 전송 요청
type ReqCoinTransfer struct {
	AUID       int64   `json:"au_id" url:"au_id"`             // 계정의 UID (Access Token에서 가져옴)
	CoinSymbol string  `json:"coin_symbol" url:"coin_symbol"` // 코인 심볼
	ToAddress  string  `json:"to_address" url:"to_address"`   // 보낼 지갑 주소
	Quantity   float64 `json:"quantity" url:"quantity"`       // 보낼 코인량

	// internal used
	TransferFee   float64 `json:"transfer_fee" url:"transfer_fee"`     // 전송 수수료
	TotalQuantity float64 `json:"total_quantity" url:"total_quantity"` // 보낼 코인량 + 전송 수수료
	ReqId         string  `json:"reqid"`
	TransactionId string  `json:"transaction_id"`
}

func NewReqCoinTransfer() *ReqCoinTransfer {
	return new(ReqCoinTransfer)
}

func (o *ReqCoinTransfer) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

///////// 코인 외부 지갑 전송 콜백 응답
// 입금 콜백
type ReqCoinTransferResDeposit struct {
	Id             int64  `json:"id"`
	CoinSymbol     string `json:"coinSymbol"`
	FromAddress    string `json:"fromAddress"`
	ToAddress      string `json:"toAddress"`
	Amount         string `json:"amount"`
	TxId           string `json:"txid"`
	OutputIndex    int64  `json:"outputindex"`
	Data           string `json:"data"`
	BlockHeight    int64  `json:"blockHeight"`
	DwDate         string `json:"dwDate"`
	DwModifiedDate string `json:"dwModifiedDate"`
}

func NewReqCoinTransferResDeposit() *ReqCoinTransferResDeposit {
	return new(ReqCoinTransferResDeposit)
}

func (o *ReqCoinTransferResDeposit) CheckValidate() *base.BaseResponse {
	return nil
}

// 출금 콜백
type ReqCoinTransferResWithdrawal struct {
	Idx               int64  `json:"idx"`
	RequestId         string `json:"requestId"`
	Type              string `json:"transfer"`
	Status            string `json:"status"`
	CoinSymbol        string `json:"coinSymbol"`
	Txid              string `json:"txid"`
	FromAddress       string `json:"fromAddress"`
	ToAddress         string `json:"toAddress"`
	Amount            string `json:"amount"`
	ActualFee         string `json:"actualFee"`
	Data              string `json:"data"`
	ExpireBlockHeight int64  `json:"expireBlockHeight"`
	Nonce             int64  `json:"nonce"`
	WebhookStatus     int64  `json:"webhookStatus"`
	CreatedDate       string `json:"createdDate"`
	ModifiedDate      string `json:"modifiedDate"`
	FinishedDate      string `json:"finishedDate"`
}

func NewReqCoinTransferResWithdrawal() *ReqCoinTransferResWithdrawal {
	return new(ReqCoinTransferResWithdrawal)
}

func (o *ReqCoinTransferResWithdrawal) CheckValidate() *base.BaseResponse {
	return nil
}

////////////////////////////////////////
