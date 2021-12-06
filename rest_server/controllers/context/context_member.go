package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
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

type Point struct {
	PointID  int64 `json:"point_id"`
	Quantity int64 `json:"quantity"`
}

// 회원 추가
type ReqPointMemberRegister struct {
	AUID       int64  `json:"au_id"`
	CUID       string `json:"cu_id"`
	AppID      int64  `json:"app_id"`
	DatabaseID int64  `json:"database_id"`
}

func NewReqPointMemberRegister() *ReqPointMemberRegister {
	return new(ReqPointMemberRegister)
}

func (o *ReqPointMemberRegister) CheckValidate() *base.BaseResponse {
	if o.AUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AUID)
	}
	if len(o.CUID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_CUID)
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
	CUID   string   `json:"cu_id"`
	AppID  int64    `json:"app_id"`
	Points *[]Point `json:"points"`
}
