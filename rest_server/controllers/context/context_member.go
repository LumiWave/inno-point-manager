package context

import (
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
	CoinSymbol    string `json:"coin_symbol"`
	WalletAddress string `json:"wallet_address"`
	CoinQuantity  string `json:"coin_quantity"`
}

type ResPointMemberWallet struct {
	AUID       int64        `json:"au_id"`
	WalletInfo []WalletInfo `json:"wallet_info"`
}

////////////////////////////////////////
