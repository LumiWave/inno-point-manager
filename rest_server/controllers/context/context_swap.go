package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

const (
	SWAP_status_init                 = int64(0)
	SWAP_status_fee_transfer_start   = int64(2)
	SWAP_status_fee_transfer_success = int64(3)
	SWAP_status_fee_transfer_fail    = int64(4)

	SWAP_status_token_transfer_start   = int64(5)
	SWAP_status_token_transfer_success = int64(6)
	SWAP_status_token_transfer_fail    = int64(7)
)

// /////// member app coin swap 요청
type SwapPoint struct {
	MUID                  int64 `json:"mu_id"`
	AppID                 int64 `json:"app_id"` // 요청 인자
	DatabaseID            int64 `json:"database_id"`
	PointID               int64 `json:"point_id"` // 요청 인자
	PreviousPointQuantity int64 `json:"previous_point_quantity"`
	AdjustPointQuantity   int64 `json:"adjust_point_quantity"` // 요청 인자
	PointQuantity         int64 `json:"point_quantity"`
}

type SwapCoin struct {
	CoinID             int64   `json:"coin_id"` // 요청 인자
	CoinSymbol         string  `json:"coin_symbol"`
	BaseCoinID         int64   `json:"base_coin_id"` // 요청 인자
	BaseCoinSymbol     string  `json:"base_coin_symbol"`
	WalletAddress      string  `json:"walletaddress"`
	AdjustCoinQuantity float64 `json:"adjust_coin_quantity"` // 요청 인자
}

type ReqSwapInfo struct {
	AUID int64 `json:"au_id"`

	SwapPoint `json:"point"`
	SwapCoin  `json:"coin"`

	TxType int64 `json:"tx_type"` // 3: point->coin,  4: coin->point

	SwapFee           float64 `json:"swap_fee"` // point->coin 시 전환시 부모지갑에 전송될 코인량 coin->point는 0 고정
	SwapWalletAddress string  `json:"swap_fee_to_wallet"`
	InnoUID           string  `json:"inno_uid"`
	TxID              int64   `json:"tx_id"`
	CreateAt          int64   `json:"create_at"`
	TxHash            string  `json:"tx_hash"`
	TxStatus          int64   `json:"tx_status"`
}

func NewReqSwapInfo() *ReqSwapInfo {
	return new(ReqSwapInfo)
}

func (o *ReqSwapInfo) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////

// swap 상태 변경 요청 : (수료 전송 후 tx정보 저장)
type ReqSwapGasFee struct {
	TxStatus          int64  `json:"tx_status"`
	TxHash            string `json:"tx_hash"`
	FromWalletAddress string `json:"from_wallet_address"`
}

func NewReqSwapStatus() *ReqSwapGasFee {
	return new(ReqSwapGasFee)
}

func (o *ReqSwapGasFee) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.TxStatus < 2 && o.TxStatus > 4 {
		return base.MakeBaseResponse(resultcode.Result_Invalid_TxStatus)
	}
	if len(o.FromWalletAddress) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Invalid_WalletAddress_Error)
	}

	return nil
}

////////////////////////////////////////
