package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

// type TokenInfo struct {
// 	PrivateTokenAmount string `json:"private_token_amount"`
// 	PrivateWalletAddr  string `json:"private_wallet_address"`

// 	PublicTokenAmount string `json:"public_token_amount"`
// 	PublicWalletAddr  string `json:"public_wallet_address"`
// }

// ///////// member 포인트 기본 정보
// type PointMemberInfo struct {
// 	ContextKey
// 	PointAmount string `json:"point_amount"`

// 	TokenInfo

// 	CreateAt int64 `json:"create_at"`
// }

// func NewPointMemberInfo() *PointMemberInfo {
// 	return new(PointMemberInfo)
// }

// func (o *PointMemberInfo) CheckValidate(bPost bool) *base.BaseResponse {
// 	if o.CpMemberIdx == 0 {
// 		return base.MakeBaseResponse(resultcode.Result_Require_MemberIdx)
// 	}
// 	if bPost {
// 		if len(o.PointAmount) == 0 {
// 			o.PointAmount = "0"
// 		}
// 		if len(o.PrivateTokenAmount) == 0 {
// 			o.PrivateTokenAmount = "0"
// 		}
// 		if len(o.PublicTokenAmount) == 0 {
// 			o.PublicTokenAmount = "0"
// 		}
// 	}
// 	if len(o.PrivateWalletAddr) != 0 && !ethcontroller.CheckValidAddress(o.PrivateWalletAddr) {
// 		return base.MakeBaseResponse(resultcode.Result_Require_ValidPrivateWalletAddr)
// 	}
// 	if len(o.PublicWalletAddr) != 0 && !ethcontroller.CheckValidAddress(o.PublicWalletAddr) {
// 		return base.MakeBaseResponse(resultcode.Result_Require_ValidPublicWalletAddr)
// 	}

// 	return nil
// }

////////////////////////////////////////

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
