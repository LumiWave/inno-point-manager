package context

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

///////// 회원 추가
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

///////// 지갑 정보 조회
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

///////// 코인 정보 조회
type AccountCoin struct {
	CoinID                    int64     `json:"coin_id"`
	BaseCoinID                int64     `json:"base_coin_id"`
	WalletAddress             string    `json:"wallet_address"`
	Quantity                  float64   `json:"quantity"`
	TodayAcqQuantity          float64   `json:"today_acq_quantity" query:"today_acq_quantity"`
	TodayCnsmQuantity         float64   `json:"today_cnsm_quantity" query:"today_cnsm_quantity"`
	TodayAcqExchangeQuantity  float64   `json:"today_acq_exchange_quantity" query:"today_acq_exchange_quantity"`
	TodayCnsmExchangeQuantity float64   `json:"today_cnsm_exchange_quantity" query:"today_cnsm_exchange_quantity"`
	ResetDate                 time.Time `json:"reset_date" query:"reset_date"`
}

////////////////////////////////////////

///////// 코인 정보 조회 by 지갑 주소
type AccountCoinByWalletAddress struct {
	AUID     int64   `json:"au_id"`
	CoinID   int64   `json:"coin_id"`
	Quantity float64 `json:"quantity"`
}

////////////////////////////////////////
