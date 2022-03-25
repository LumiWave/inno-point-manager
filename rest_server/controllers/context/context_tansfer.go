package context

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

///////// 코인 외부 지갑 전송 요청 : 부모지갑
type ReqCoinTransferFromParentWallet struct {
	AUID       int64   `json:"au_id" url:"au_id"` // 계정의 UID (Access Token에서 가져옴)
	CoinID     int64   `json:"coin_id" url:"coin_id"`
	CoinSymbol string  `json:"coin_symbol" url:"coin_symbol"` // 코인 심볼
	ToAddress  string  `json:"to_address" url:"to_address"`   // 보낼 지갑 주소
	Quantity   float64 `json:"quantity" url:"quantity"`       // 보낼 코인량

	// internal used
	TransferFee   float64 `json:"transfer_fee" url:"transfer_fee"`     // 전송 수수료
	TotalQuantity float64 `json:"total_quantity" url:"total_quantity"` // 보낼 코인량 + 전송 수수료
	ReqId         string  `json:"reqid"`
	TransactionId string  `json:"transaction_id"`

	ActionDate time.Time `json:"action_date"`
}

func NewReqCoinTransferFromParentWallet() *ReqCoinTransferFromParentWallet {
	return new(ReqCoinTransferFromParentWallet)
}

func (o *ReqCoinTransferFromParentWallet) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

///////// 코인 외부 지갑 전송 요청 : 특정지갑
type ReqCoinTransferFromUserWallet struct {
	AUID           int64   `json:"au_id" url:"au_id"` // 계정의 UID (Access Token에서 가져옴)
	CoinID         int64   `json:"coin_id" url:"coin_id"`
	CoinSymbol     string  `json:"coin_symbol" url:"coin_symbol"`           // 코인 심볼
	BaseCoinSymbol string  `json:"base_coin_symbol" url:"base_coin_symbol"` // base 코인 심볼
	FromAddress    string  `json:"from_address" url:"from_address"`         // 보내는 지갑 주소
	ToAddress      string  `json:"to_address" url:"to_address"`             // 보낼 지갑 주소
	Quantity       float64 `json:"quantity" url:"quantity"`                 // 보낼 코인량

	// internal used
	TransferFee   float64 `json:"transfer_fee" url:"transfer_fee"`     // 전송 수수료
	TotalQuantity float64 `json:"total_quantity" url:"total_quantity"` // 보낼 코인량 + 전송 수수료
	ReqId         string  `json:"reqid"`
	TransactionId string  `json:"transaction_id"`

	ActionDate time.Time `json:"action_date"`
}

func NewReqCoinTransferFromUserWallet() *ReqCoinTransferFromUserWallet {
	return new(ReqCoinTransferFromUserWallet)
}

func (o *ReqCoinTransferFromUserWallet) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

///////// transfer 중인 상태 정보 요청
type GetCoinTransferExistInProgress struct {
	AUID int64 `json:"au_id" query:"au_id"`
}

func NewGetCoinTransferExistInProgress() *GetCoinTransferExistInProgress {
	return new(GetCoinTransferExistInProgress)
}

func (o *GetCoinTransferExistInProgress) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
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

///////// 코인 수수료
type ReqCoinFee struct {
	Symbol string `query:"symbol"`
}

func NewReqCoinFee() *ReqCoinFee {
	return new(ReqCoinFee)
}

func (o *ReqCoinFee) CheckValidate() *base.BaseResponse {
	if len(o.Symbol) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_Symbol)
	}

	return nil
}

type ResCoinFeeInfo struct {
	Fast    string `json:"fast"`
	Slow    string `json:"slow"`
	Average string `json:"average"`
	BaseFee string `json:"basefee"`
	Fastest string `json:"fastest"`
}

////////////////////////////////////////
