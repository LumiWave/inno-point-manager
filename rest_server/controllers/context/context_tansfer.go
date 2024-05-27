package context

import (
	"time"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
)

const (
	From_user_to_fee_wallet    = int64(0) // 유저 지갑에서 수수료 지갑으로 전송
	From_user_to_parent_wallet = int64(1) // 유저 지갑에서 부모 지갑으로 전송
	From_user_to_other_wallet  = int64(2) // 유저 지갑에서 다른 지갑으로 전송

	From_parent_to_other_wallet = int64(3) // 부모 지갑에서 다른 지갑으로 전송
)

type TxType struct {
	Target int64 `json:"target"`
	AUID   int64 `json:"au_id"`
	CoinID int64 `json:"coin_id"`
}

// /////// 코인 부모 지갑 전송 요청 : to 부모지갑
type ReqCoinTransferToParentWallet struct {
	AUID       int64   `json:"au_id" url:"au_id"` // 계정의 UID (Access Token에서 가져옴)
	CoinID     int64   `json:"coin_id" url:"coin_id"`
	CoinSymbol string  `json:"coin_symbol" url:"coin_symbol"` // 코인 심볼
	ToAddress  string  `json:"to_address" url:"to_address"`   // 보낼 지갑 주소
	Quantity   float64 `json:"quantity" url:"quantity"`       // 보낼 코인량

	// internal used
	ReqId         string `json:"reqid"`
	TransactionId string `json:"transaction_id"`

	ActionDate time.Time `json:"action_date"`
}

func NewReqCoinTransferToParentWallet() *ReqCoinTransferToParentWallet {
	return new(ReqCoinTransferToParentWallet)
}

func (o *ReqCoinTransferToParentWallet) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

// /////// 코인 외부 지갑 전송 요청 : 부모지갑
type ReqCoinTransferFromParentWallet struct {
	AUID             int64   `json:"au_id" url:"au_id"` // 계정의 UID (Access Token에서 가져옴)
	CoinID           int64   `json:"coin_id" url:"coin_id"`
	CoinSymbol       string  `json:"coin_symbol" url:"coin_symbol"` // 코인 심볼
	ToAddress        string  `json:"to_address" url:"to_address"`   // 보낼 지갑 주소
	Quantity         float64 `json:"quantity" url:"quantity"`       // 보낼 코인량
	IsNormalTransfer bool    `json:"is_normal_transfer"`            // 단순전송

	// internal used
	ReqId         string `json:"reqid"`
	TransactionId string `json:"transaction_id"`

	ActionDate time.Time `json:"action_date"`
}

func NewReqCoinTransferFromParentWallet() *ReqCoinTransferFromParentWallet {
	return new(ReqCoinTransferFromParentWallet)
}

func (o *ReqCoinTransferFromParentWallet) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

// /////// 코인 외부 지갑 전송 요청 : 특정지갑
type ReqCoinTransferFromUserWallet struct {
	AUID           int64   `json:"au_id" url:"au_id"` // 계정의 UID (Access Token에서 가져옴)
	CoinID         int64   `json:"coin_id" url:"coin_id"`
	CoinSymbol     string  `json:"coin_symbol" url:"coin_symbol"`           // 코인 심볼
	BaseCoinSymbol string  `json:"base_coin_symbol" url:"base_coin_symbol"` // base 코인 심볼
	FromAddress    string  `json:"from_address" url:"from_address"`         // 보내는 지갑 주소
	ToAddress      string  `json:"to_address" url:"to_address"`             // 보낼 지갑 주소
	Quantity       float64 `json:"quantity" url:"quantity"`                 // 보낼 코인량

	// internal used
	Target        int64  `json:"target"` //0 : to fee wallet, 1:external wallet
	TransactionId string `json:"transaction_id"`

	ActionDate time.Time `json:"action_date"`
}

func NewReqCoinTransferFromUserWallet() *ReqCoinTransferFromUserWallet {
	return new(ReqCoinTransferFromUserWallet)
}

func (o *ReqCoinTransferFromUserWallet) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

// /////// transfer 중인 상태 정보 요청
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

// /////// 코인 외부 지갑 전송 콜백 응답
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

// sui coin object id 리스트 조회
type ReqCoinObjects struct {
	WalletAddress   string `query:"wallet_address"`
	ContractAddress string `query:"contract_address"`
}

func NewReqCoinObjects() *ReqCoinObjects {
	return new(ReqCoinObjects)
}

func (o *ReqCoinObjects) CheckValidate() *base.BaseResponse {
	if len(o.WalletAddress) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Invalid_WalletAddress_Error)
	}
	if len(o.ContractAddress) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Invalid_ContractAddress)
	}
	return nil
}

type ResCoinObjects struct {
	ObjectIDs []string `json:"object_ids"`
}

////////////////////////////////////////

// /////// 코인 수수료
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

// /////// 지갑 잔액
type ReqBalance struct {
	Symbol  string `query:"symbol"`
	Address string `query:"address"`
}

func NewReqBalance() *ReqBalance {
	return new(ReqBalance)
}

func (o *ReqBalance) CheckValidate() *base.BaseResponse {
	if len(o.Symbol) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_Symbol)
	}

	return nil
}

type ResReqBalance struct {
	Balance string `json:"balance"`
	Address string `json:"address"`
}

// /////// 전체 지갑 잔액
type ReqBalanceAll struct {
	AUID int64 `query:"au_id"`
}

func NewReqBalanceAll() *ReqBalanceAll {
	return new(ReqBalanceAll)
}

func (o *ReqBalanceAll) CheckValidate() *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}

	return nil
}

type Balance struct {
	CoinID     int64  `json:"coin_id"`
	BaseCoinID int64  `json:"base_coin_id"`
	Symbol     string `json:"symbol"`
	Balance    string `json:"balance"`
	Address    string `json:"address"`
	Decimal    int64  `json:"decimal"`
}
type ResBalanceAll struct {
	Balances []*Balance `json:"balances"`
}

////////////////////////////////////////

// /////// coin mainnet 보정
type CoinReload struct {
	AUID int64 `json:"au_id" query:"au_id"`
}

func NewCoinReload() *CoinReload {
	return new(CoinReload)
}

func (o *CoinReload) CheckValidate() *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
	return nil
}

////////////////////////////////////////
