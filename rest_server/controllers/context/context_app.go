package context

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
)

type Point struct {
	PointID       int64  `json:"point_id"`
	Quantity      int64  `json:"quantity"`
	TodayQuantity int64  `json:"today_quantity"`
	ResetDate     string `json:"reset_date"`

	PreQuantity    int64 `json:"previous_quantity,omitempty"`
	AdjustQuantity int64 `json:"adjust_quantity,omitempty"`
}

type PointInfo struct {
	MyUuid     string `json:"my_uuid,omitempty"`
	DatabaseID int64  `json:"database_id"`

	MUID   int64    `json:"mu_id"`
	Points []*Point `json:"points"`
}

// /////// member 포인트 조회
type ReqGetPointApp struct {
	AppId      int64 `query:"app_id"`
	MUID       int64 `query:"mu_id"`
	PointID    int64 `query:"point_id"`
	DatabaseID int64 `query:"database_id"`
}

func NewReqGetPointApp() *ReqGetPointApp {
	return new(ReqGetPointApp)
}

func (o *ReqGetPointApp) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if o.MUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MUID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}
	if o.AppId != 0 && ctx.GetValue() != nil {
		ctx.GetValue().AppID = o.AppId
	}
	return nil
}

////////////////////////////////////////

// /////// app 포인트 업데이트
type ReqPointAppUpdate struct {
	AppID      int64 `json:"app_id"`
	MUID       int64 `json:"mu_id"`
	PointID    int64 `json:"point_id"`
	DatabaseID int64 `json:"database_id"`

	PreQuantity    int64 `json:"previous_quantity"`
	AdjustQuantity int64 `json:"adjust_quantity"`
}

func NewReqPointMemberAppUpdate() *ReqPointAppUpdate {
	return new(ReqPointAppUpdate)
}

func (o *ReqPointAppUpdate) CheckValidate(ctx *PointManagerContext) *base.BaseResponse {
	if ctx.GetValue() != nil && ctx.GetValue().AppID != 0 {
		o.AppID = ctx.GetValue().AppID
	}

	if o.MUID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_MUID)
	}
	if o.PointID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_PointID)
	}
	if o.DatabaseID == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_DatabaseID)
	}
	if o.AdjustQuantity == 0 {
		return base.MakeBaseResponse(resultcode.Result_Require_AdjustQuantity)
	}
	return nil
}

type ResPointAppUpdate struct {
	MUID    int64 `json:"mu_id"`
	PointID int64 `json:"point_id"`

	PreQuantity   int64 `json:"previous_quantity"`
	TodayQuantity int64 `json:"today_quantity"`
}

////////////////////////////////////////

// /////// BaseCoinInfo
type BaseCoinInfo struct {
	BaseCoinID         int64  `json:"base_coin_id"`
	BaseCoinName       string `json:"base_coin_name"`
	BaseCoinSymbol     string `json:"base_coin_symbol"`
	IsUsedParentWallet bool   `json:"is_used_parent_wallet"`
	AccessWallet       string `json:"access_wallet"`
}

type BaseCoinList struct {
	Coins []*BaseCoinInfo `json:"base_coins"`
}

////////////////////////////////////////
