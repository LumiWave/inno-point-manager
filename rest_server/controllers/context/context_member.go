package context

import (
	"time"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
)

// /////// 회원 추가
type ReqPointMemberRegister struct {
	AUID       int64 `json:"au_id"`
	MUID       int64 `json:"mu_id"`
	AppID      int64 `json:"app_id"`
	DatabaseID int64 `json:"database_id"`
}

func NewReqPointMemberRegister() *ReqPointMemberRegister {
	return new(ReqPointMemberRegister)
}

func (o *ReqPointMemberRegister) CheckValidate() *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
	if o.MUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MUID)
	}
	if o.AppID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AppID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}

	return nil
}

type ResPointMemberRegister struct {
	PointInfo
}

type InsertPointMemberInfo struct {
	PointID  int64
	Quantity int64
}

////////////////////////////////////////

// /////// 지갑 정보 조회
type ReqPointMemberWallet struct {
	AUID int64 `query:"au_id"`
}

func NewPointMemberWallet() *ReqPointMemberWallet {
	return new(ReqPointMemberWallet)
}

func (o *ReqPointMemberWallet) CheckValidate() *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
	return nil
}

type WalletInfo struct {
	CoinID        int64  `json:"coin_id"`
	BaseCoinID    int64  `json:"base_coin_id"`
	CoinSymbol    string `json:"coin_symbol"`
	WalletAddress string `json:"wallet_address"`
	CoinQuantity  string `json:"coin_quantity"`
}

type ResPointMemberWallet struct {
	AUID       int64        `json:"au_id"`
	WalletInfo []WalletInfo `json:"wallet_info"`
}

////////////////////////////////////////

// /////// 코인 정보 조회
type AccountCoin struct {
	CoinID        int64  `json:"coin_id"`
	BaseCoinID    int64  `json:"base_coin_id"`
	WalletAddress string `json:"wallet_address"`
	//Quantity                  float64   `json:"quantity"`
	TodayAcqQuantity          float64   `json:"today_acq_quantity" query:"today_acq_quantity"`
	TodayCnsmQuantity         float64   `json:"today_cnsm_quantity" query:"today_cnsm_quantity"`
	TodayExchangeAcqQuantity  float64   `json:"today_exchange_acq_quantity" query:"today_exchange_acq_quantity"`
	TodayExchangeCnsmQuantity float64   `json:"today_exchange_cnsm_quantity" query:"today_exchange_cnsm_quantity"`
	ResetDate                 time.Time `json:"reset_date" query:"reset_date"`
}

////////////////////////////////////////

// /////// 코인 정보 조회 by 지갑 주소
type AccountCoinByWalletAddress struct {
	AUID     int64   `json:"au_id"`
	CoinID   int64   `json:"coin_id"`
	Quantity float64 `json:"quantity"`
}

////////////////////////////////////////

// 내 지갑 정보
type AccountWallet struct {
	WalletID         int64  `json:"walllet_id"`
	BaseCoinID       int64  `json:"base_coin_id"`
	WalletAddress    string `json:"wallet_address"`
	WalletTypeID     int64  `json:"wallet_type_id"`
	ConnectionStatus int    `json:"connection_status"`
	ModifiedDT       string `json:"modified_dt"`
}

////////////////////////////////////////

// /////// Me App Point List
type ReqMeAppPoint struct {
	AUID              int64  `json:"au_id" query:"au_id"`
	MUID              int64  `json:"mu_id" query:"mu_id"`
	AppID             int64  `json:"app_id"`
	PointID           int64  `json:"point_id" query:"point_id"`
	TodayAcqQuantity  int64  `json:"today_acq_quantity" query:"today_acq_quantity"`
	TodayCnsmQuantity int64  `json:"today_cnsm_quantity" query:"today_cnsm_quantity"`
	ResetDate         string `json:"reset_date" query:"reset_date"`
}

////////////////////////////////////////
