package context

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
)

const (
	SWAP_status_init                 = int64(1)
	SWAP_status_fee_transfer_start   = int64(2)
	SWAP_status_fee_transfer_success = int64(3)
	SWAP_status_fee_transfer_fail    = int64(4)

	SWAP_status_token_transfer_withdrawal_start   = int64(5) // 법인 지갑에서 토큰 출금 시작
	SWAP_status_token_transfer_withdrawal_success = int64(6) // 법인 지갑에서 토큰 출금 성공
	SWAP_status_token_transfer_withdrawal_fail    = int64(7) // 법인 지갑에서 토큰 출금 실패

	SWAP_status_token_transfer_deposit_start   = int64(8) // 법인 지갑으로 토큰 입금 시작
	SWAP_status_token_transfer_deposit_success = int64(9) // 법인 지갑으로 토큰 입금 성공
	SWAP_status_token_transfer_deposit_fail    = int64(0) // 법인 지갑으로 토큰 입금 실패
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
	ToWalletAddress    string  `json:"to_wallet"`
	WalletTypeID       int64   `json:"wallet_type_id"`
	WalletID           int64   `json:"wallet_id"`
	AdjustCoinQuantity float64 `json:"adjust_coin_quantity"` // 요청 인자
	TokenTxHash        string  `json:"token_tx_hash"`        // swap 코인 전송 txhash
	IsComplete         bool    `json:"is_complete"`          // 전송 완료 여부
}

type ReqSwapInfo struct {
	AUID int64 `json:"au_id"`

	//SwapPoint `json:"point"`
	//SwapCoin `json:"coin"`

	SwapFromPoint SwapPoint `json:"from_point"`
	SwapToPoint   SwapPoint `json:"to_point"`

	SwapFromCoin SwapCoin `json:"from_coin"`
	SwapToCoin   SwapCoin `json:"to_coin"`

	TxType int64 `json:"tx_type"` // 3: point->coin,  4: coin->point, 26: coin->coin

	SwapFeeCoinID     int64   `json:"swap_fee_coin_id"` // 코인 수수료 전송용 코인 아이디
	SwapFeeCoinSymbol string  `json:"swap_fee_coin_symbol"`
	SwapFee           float64 `json:"swap_fee"` // point->coin 시 전환시 부모지갑에 전송될 코인량 coin->point는 0 고정
	SwapFeeT          string  `json:"swap_fee_string"`
	SwapFeeD          string  `json:"swap_fee_string_d"`
	ToWalletAddress   string  `json:"to_wallet"`
	InnoUID           string  `json:"inno_uid"`
	TxID              int64   `json:"tx_id"`
	CreateAt          int64   `json:"create_at"`
	TxHash            string  `json:"tx_hash"`
	IsFeeComplete     bool    `json:"is_fee_complete"`
	TxStatus          int64   `json:"tx_status"`
	TxGasFee          float64 `json:"tx_gas_fee"`
}

func NewReqSwapInfo() *ReqSwapInfo {
	return new(ReqSwapInfo)
}

// 대시보드가 프리세일용 코인->포인트 스왑에 쓰는 내부 전용 타입
const EventID_Server_C2P = 100000

func (o *ReqSwapInfo) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}

	// 내주는 쪽은 줄고(음수) 받는 쪽은 늘어야(양수) 한다.
	// 부호를 확인하지 않으면 포인트를 늘리면서 코인까지 받아 갈 수 있다.
	switch o.TxType {
	case EventID_P2C:
		if !(o.SwapFromPoint.AdjustPointQuantity < 0 && o.SwapToCoin.AdjustCoinQuantity > 0) {
			return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
		}
	case EventID_C2P, EventID_Server_C2P:
		if !(o.SwapFromCoin.AdjustCoinQuantity < 0 && o.SwapToPoint.AdjustPointQuantity > 0) {
			return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
		}
	case EventID_C2C:
		if !(o.SwapFromCoin.AdjustCoinQuantity != 0 && o.SwapToCoin.AdjustCoinQuantity > 0) {
			return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
		}
	case EventID_P2P:
		if !(o.SwapFromPoint.AdjustPointQuantity != 0 && o.SwapToPoint.AdjustPointQuantity > 0) {
			return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
		}
	default:
		return base.MakeBaseResponse(resultcode.Result_Error_Invalid_data)
	}

	return nil
}

// CanCancelByUser 는 사용자 쪽 전송 실패 신고로 스왑을 취소(선차감 환불)해도 되는 상태인지 알려준다.
// 수수료나 코인 입금이 확인되면 지급이 시작되므로, 그 뒤에는 취소할 수 없다.
func (o *ReqSwapInfo) CanCancelByUser() bool {
	switch o.TxType {
	case EventID_P2C:
		if o.IsFeeComplete || len(o.SwapToCoin.TokenTxHash) != 0 {
			return false
		}
		return o.TxStatus == SWAP_status_init || o.TxStatus == SWAP_status_fee_transfer_start
	case EventID_C2P:
		return o.TxStatus == SWAP_status_init || o.TxStatus == SWAP_status_token_transfer_deposit_start
	case EventID_C2C:
		if (o.IsFeeComplete && o.SwapFromCoin.IsComplete) || len(o.SwapToCoin.TokenTxHash) != 0 {
			return false
		}
		return o.TxStatus == SWAP_status_init ||
			o.TxStatus == SWAP_status_fee_transfer_start ||
			o.TxStatus == SWAP_status_fee_transfer_success ||
			o.TxStatus == SWAP_status_token_transfer_deposit_start
	}
	return false
}

// CanReportDeposit 은 사용자가 "코인 전송을 시작했다"고 알려도 되는 상태인지 알려준다.
// 코인을 받는 스왑(코인->포인트, 코인->코인)에서 입금이 확인되기 전에만 받는다.
func (o *ReqSwapInfo) CanReportDeposit() bool {
	if o.TxType != EventID_C2P && o.TxType != EventID_C2C {
		return false
	}
	return o.CanCancelByUser()
}

// CanPayout 은 수수료 입금 확인 뒤 법인 지갑에서 코인을 내보내도 되는 상태인지 알려준다.
// 이미 출금을 시작했거나 끝난 스왑에 같은 입금 콜백이 다시 오면 두 번 지급하지 않는다.
func (o *ReqSwapInfo) CanPayout() bool {
	if o.TxType != EventID_P2C && o.TxType != EventID_C2C {
		return false
	}
	if len(o.SwapToCoin.TokenTxHash) != 0 {
		return false
	}
	return o.TxStatus == SWAP_status_init ||
		o.TxStatus == SWAP_status_fee_transfer_start ||
		o.TxStatus == SWAP_status_fee_transfer_success ||
		o.TxStatus == SWAP_status_token_transfer_deposit_start ||
		o.TxStatus == SWAP_status_token_transfer_deposit_success
}

////////////////////////////////////////

// swap 상태 변경 요청 : (수료 전송 후 tx정보 저장)
type ReqSwapStatus struct {
	TxID              int64  `json:"tx_id"`
	TxStatus          int64  `json:"tx_status"`
	TxHash            string `json:"tx_hash"`
	FromWalletAddress string `json:"from_wallet_address"`
	// 요청한 회원의 계정. 대시보드가 토큰의 계정을 넣어 보내면 스왑의 주인인지 확인한다.
	// 값이 없으면(0) 예전처럼 확인하지 않는다.
	AUID int64 `json:"au_id"`
}

func NewReqSwapStatus() *ReqSwapStatus {
	return new(ReqSwapStatus)
}

func (o *ReqSwapStatus) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	// 사용자 쪽에서 알릴 수 있는 상태만 받는다. 출금 관련 상태(5~7, 9)는 서버만 바꾼다.
	switch o.TxStatus {
	case SWAP_status_fee_transfer_start, SWAP_status_fee_transfer_success, SWAP_status_fee_transfer_fail,
		SWAP_status_token_transfer_deposit_start, SWAP_status_token_transfer_deposit_fail:
	default:
		return base.MakeBaseResponse(resultcode.Result_Invalid_TxStatus)
	}
	if len(o.FromWalletAddress) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Invalid_WalletAddress_Error)
	}

	return nil
}

////////////////////////////////////////

// swap 진행 중인지 체크
type ReqSwapInprogress struct {
	AUID int64 `query:"au_id"`
}

func NewReqSwapIniprogress() *ReqSwapInprogress {
	return new(ReqSwapInprogress)
}

func (o *ReqSwapInprogress) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
	return nil
}

////////////////////////////////////////

// swap 정보 삭제
type DeleteDeleteSwapInfo struct {
	WalletAddress string `query:"wallet_address"`
}

// / 스왑 가능 메타 데이터 정보
// coin to coin
type SwapC2C struct {
	// FromBaseCoinID는 전환할 재료 코인의 계열 ID입니다.
	FromBaseCoinID int64 `json:"from_base_coin_id"`

	// FromID는 전환할 재료의 ID입니다.
	FromID int64 `json:"from_id"`

	// ToBaseCoinID는 받을 코인의 계열 ID입니다.
	ToBaseCoinID int64 `json:"to_base_coin_id"`

	// ToID는 받을 코인의 ID입니다.
	ToID int64 `json:"to_id"`

	// IsEnabled는 해당 전환이 활성화 되어있는지 여부를 나타냅니다.
	IsEnabled bool `json:"is_enabled"`

	// 화면에 표출 여부
	IsVisible bool `json:"is_visible"`

	// MinimumExchangeQuantity는 최소 전환량을 나타냅니다.
	MinimumExchangeQuantity string `json:"minimum_exchange_quantity"`

	// ExchangeRatio는 받을 전환 비율을 나타냅니다.
	ExchangeRatio float64 `json:"exchange_ratio"`
}

// point to coin
type SwapP2C struct {
	// FromID는 전환할 재료의 포인트 ID입니다.
	FromID int64 `json:"from_id"`

	// ToBaseCoinID는 받을 코인의 계열 ID입니다.
	ToBaseCoinID int64 `json:"to_base_coin_id"`

	// ToID는 받을 코인의 ID입니다.
	ToID int64 `json:"to_id"`

	// IsEnabled는 해당 전환이 활성화 되어있는지 여부를 나타냅니다.
	IsEnabled bool `json:"is_enabled"`

	// 화면에 표출 여부
	IsVisible bool `json:"is_visible"`

	// MinimumExchangeQuantity는 최소 전환량을 나타냅니다.
	MinimumExchangeQuantity string `json:"minimum_exchange_quantity"`

	// ExchangeRatio는 받을 전환 비율을 나타냅니다.
	ExchangeRatio float64 `json:"exchange_ratio"`
}

// coin to point 정보
type SwapC2P struct {
	// FromBaseCoinID는 전환할 재료 코인의 계열 ID입니다.
	FromBaseCoinID int64 `json:"from_base_coin_id"`

	// FromID는 전환할 재료의 ID입니다.
	FromID int64 `json:"from_id"`

	// ToID는 받을 포인트의 ID입니다.
	ToID int64 `json:"to_id"`

	// IsEnabled는 해당 전환이 활성화 되어있는지 여부를 나타냅니다.
	IsEnabled bool `json:"is_enabled"`

	// 화면에 표출 여부
	IsVisible bool `json:"is_visible"`

	// MinimumExchangeQuantity는 최소 전환량을 나타냅니다.
	MinimumExchangeQuantity string `json:"minimum_exchange_quantity"`

	// ExchangeRatio는 받을 전환 비율을 나타냅니다.
	ExchangeRatio float64 `json:"exchange_ratio"`
}

// point to point 스왑 정보
type SwapP2P struct {
	// FromID는 전환할 재료의 ID입니다.
	FromID int64 `json:"from_id"`

	// ToID는 받을 포인트의 ID입니다.
	ToID int64 `json:"to_id"`

	// IsEnabled는 해당 전환이 활성화 되어있는지 여부를 나타냅니다.
	IsEnabled bool `json:"is_enabled"`

	// 화면에 표출 여부
	IsVisible bool `json:"is_visible"`

	// MinimumExchangeQuantity는 최소 전환량을 나타냅니다.
	MinimumExchangeQuantity string `json:"minimum_exchange_quantity"`

	// ExchangeRatio는 받을 전환 비율을 나타냅니다.
	ExchangeRatio float64 `json:"exchange_ratio"`
}

////////////////////////////////////////

// //////////////////////////////////////
// ssr point <-> ssrm chip swap 정보 조회
type ReqGameChipSwapInfo struct {
}

func NewReqGameChipSwapInfo() *ReqGameChipSwapInfo {
	return new(ReqGameChipSwapInfo)
}

func (o *ReqGameChipSwapInfo) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

type ResGameChipSwapInfo struct {
	SSRQuantity            int64 `json:"ssr_quantity"`
	SSRMQuantity           int64 `json:"ssrm_quantity"`
	DailyExchangeAcqChips  int64 `json:"daily_exchange_acq_chips"`
	DailyExchangeCnsmChips int64 `json:"daily_exchange_cnsm_chips"`

	ExchangeRatio float64 `json:"exchange_ratio"`
	SSRToSSRMID   int64   `json:"ssr_to_ssrm_id"`
	SSRMToSSRID   int64   `json:"ssrm_to_ssr_id"`
	DailyMaxAcq   *int64  `json:"daily_max_acq,omitempty"`
	DailyMaxCnsm  *int64  `json:"daily_max_cnsm,omitempty"`
}

////////////////////////////////////////

// //////////////////////////////////////
// ssr point <-> ssrm chip swap 요청
type ReqGameChipSwap struct {
	SwapType       int64 `json:"swap_type"`
	SSRAdjustPoint int64 `json:"ssr_adjust_point"`
	SSRMAdjustChip int64 `json:"ssrm_adjust_chip"`
}

func NewReqGameChipSwap() *ReqGameChipSwap {
	return new(ReqGameChipSwap)
}

func (o *ReqGameChipSwap) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

type ResReqGameChipSwap struct {
}

////////////////////////////////////////

// //////////////////////////////////////
// ssr point <-> ssrm chip swap 비율 변경
type ReqGameChipSwapRatio struct {
	ExchangeRatio float64 `json:"exchange_ratio"`
	DailyMaxAcq   int64   `json:"daily_max_acq"`
	DailyMaxCnsm  int64   `json:"daily_max_cnsm"`
}

func NewReqGameChipSwapRatio() *ReqGameChipSwapRatio {
	return new(ReqGameChipSwapRatio)
}

func (o *ReqGameChipSwapRatio) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	return nil
}

////////////////////////////////////////
