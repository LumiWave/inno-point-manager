package commonapi

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

// 외부 API 는 요청에 적힌 회원 번호를 그대로 믿지 않고, 토큰의 신원과 맞는지 여기서 확인한다.
// 내부 API(대시보드 등 우리 서버가 호출)는 토큰이 없으므로 이 검사를 거치지 않는다.
//
// 회원(MUID)이 어느 앱 소속인지 바로 알려주는 프로시저는 없다. 대신 포인트는 앱마다 따로 있고
// 회원은 자기 앱의 포인트만 가지므로, "포인트가 토큰의 앱 것인가"로 앱 소속을 확인한다.

func authEnabled() bool {
	return config.GetInstance().Auth.AuthEnable
}

// appHasPoint 는 포인트가 그 앱의 포인트인지 메타 정보로 확인한다.
func appHasPoint(appID, pointID int64) bool {
	return appHasPointIn(model.GetDB().AppPointsMap, appID, pointID)
}

func appHasPointIn(appPoints map[int64]*model.AppPointInfo, appID, pointID int64) bool {
	app, ok := appPoints[appID]
	if !ok || app == nil {
		return false
	}
	_, ok = app.PointsMap[pointID]
	return ok
}

// accountHasMember 는 계정에 그 회원이 속하는지 확인한다.
func accountHasMember(auID, muID, databaseID int64) bool {
	_, members, err := model.GetDB().USPAU_GetList_Members(auID)
	if err != nil {
		log.Errorf("accountHasMember USPAU_GetList_Members fail [au_id:%v] : %v", auID, err)
		return false
	}
	for _, member := range members {
		if member.MUID == muID && member.DatabaseID == databaseID {
			return true
		}
	}
	return false
}

// pointsBelongToApp 는 회원의 포인트 중 그 앱의 포인트가 있는지 확인한다.
func pointsBelongToApp(appPoints map[int64]*model.AppPointInfo, points []*context.Point, appID int64) bool {
	for _, point := range points {
		if point != nil && appHasPointIn(appPoints, appID, point.PointID) {
			return true
		}
	}
	return false
}

// CheckPointUpdate 는 포인트 조정 요청이 앱(게임 서버) 토큰으로 왔는지, 조정할 포인트가 그 앱의 것인지 확인한다.
// 회원 토큰으로 본문의 app_id, mu_id 를 바꿔 포인트를 조정하거나, 다른 앱의 포인트를 건드리지 못하게 한다.
// 통과하면 app_id 는 본문 값 대신 토큰의 앱으로 고정한다.
func CheckPointUpdate(ctx *context.PointManagerContext, req *context.ReqPointAppUpdate) *base.BaseResponse {
	if !authEnabled() {
		return nil
	}
	if !ctx.IsAppToken() {
		log.Errorf("CheckPointUpdate not app token [logintype:%v][app:%v][au_id:%v]", ctx.GetValue().LoginType, ctx.GetValue().AppID, ctx.GetValue().AUID)
		return base.MakeBaseResponse(resultcode.Result_Auth_InvalidLoginInfo)
	}
	req.AppID = ctx.GetValue().AppID
	if !appHasPoint(req.AppID, req.PointID) {
		log.Errorf("CheckPointUpdate other app point [token app:%v][point_id:%v][mu_id:%v]", req.AppID, req.PointID, req.MUID)
		return base.MakeBaseResponse(resultcode.Result_Require_PointID)
	}
	return nil
}

// CheckPointAccess 는 회원 포인트를 조회해도 되는 토큰인지 확인한다.
// 회원 토큰은 본인 것만 조회할 수 있다. 앱 토큰은 조회한 뒤 그 앱의 포인트를 가진 회원인지 확인한다(CheckPointOwnerApp).
func CheckPointAccess(ctx *context.PointManagerContext, muID, databaseID int64) *base.BaseResponse {
	if !authEnabled() || ctx.IsAppToken() {
		return nil
	}
	if !ctx.IsMemberToken() {
		return base.MakeBaseResponse(resultcode.Result_Auth_InvalidLoginInfo)
	}
	if !accountHasMember(ctx.GetValue().AUID, muID, databaseID) {
		log.Errorf("CheckPointAccess not owner [token au_id:%v][mu_id:%v][database_id:%v]", ctx.GetValue().AUID, muID, databaseID)
		return base.MakeBaseResponse(resultcode.Result_Error_NotExistMember)
	}
	return nil
}

// CheckPointOwnerApp 은 앱 토큰으로 조회한 회원이 그 앱의 포인트를 가진 회원인지 확인한다.
// 다른 앱 회원의 포인트를 번호만 바꿔 조회하지 못하게 한다. 앱 토큰이 아니면 검사하지 않는다.
func CheckPointOwnerApp(ctx *context.PointManagerContext, pointInfo *context.PointInfo) *base.BaseResponse {
	if !authEnabled() || !ctx.IsAppToken() {
		return nil
	}
	// 돌려줄 포인트가 없으면 새는 정보도 없으므로 그대로 둔다.
	if pointInfo == nil || len(pointInfo.Points) == 0 {
		return nil
	}
	if !pointsBelongToApp(model.GetDB().AppPointsMap, pointInfo.Points, ctx.GetValue().AppID) {
		log.Errorf("CheckPointOwnerApp other app member [token app:%v][mu_id:%v]", ctx.GetValue().AppID, pointInfo.MUID)
		return base.MakeBaseResponse(resultcode.Result_Error_NotExistMember)
	}
	return nil
}

// CheckWalletAccess 는 계정 지갑 정보를 조회해도 되는 토큰인지 확인한다.
// 회원 토큰은 요청의 au_id 대신 토큰의 계정으로 조회하고, 앱 토큰은 그 앱에 가입한 계정만 조회할 수 있다.
func CheckWalletAccess(ctx *context.PointManagerContext, req *context.ReqPointMemberWallet) *base.BaseResponse {
	if !authEnabled() {
		return nil
	}
	if ctx.IsMemberToken() {
		req.AUID = ctx.GetValue().AUID
		return nil
	}
	if !ctx.IsAppToken() {
		return base.MakeBaseResponse(resultcode.Result_Auth_InvalidLoginInfo)
	}
	_, members, err := model.GetDB().USPAU_GetList_Members(req.AUID)
	if err != nil {
		log.Errorf("CheckWalletAccess USPAU_GetList_Members fail [au_id:%v] : %v", req.AUID, err)
		return base.MakeBaseResponse(resultcode.Result_Error_NotExistMember)
	}
	if _, ok := members[ctx.GetValue().AppID]; !ok {
		log.Errorf("CheckWalletAccess not app member [token app:%v][au_id:%v]", ctx.GetValue().AppID, req.AUID)
		return base.MakeBaseResponse(resultcode.Result_Error_NotExistMember)
	}
	return nil
}

// BindSwapOwner 는 외부에서 온 스왑 요청을 토큰의 계정으로 고정하고, 포인트 회원이 그 계정 소속인지 확인한다.
// 본문의 au_id, mu_id 를 믿으면 남의 계정 이름으로 스왑을 걸 수 있다.
func BindSwapOwner(ctx *context.PointManagerContext, params *context.ReqSwapInfo) *base.BaseResponse {
	if !authEnabled() {
		return nil
	}
	if !ctx.IsMemberToken() {
		return base.MakeBaseResponse(resultcode.Result_Auth_InvalidLoginInfo)
	}
	params.AUID = ctx.GetValue().AUID

	points := []*context.SwapPoint{}
	switch params.TxType {
	case context.EventID_P2C:
		points = append(points, &params.SwapFromPoint)
	case context.EventID_C2P, context.EventID_Server_C2P:
		points = append(points, &params.SwapToPoint)
	case context.EventID_P2P:
		points = append(points, &params.SwapFromPoint, &params.SwapToPoint)
	}
	for _, point := range points {
		if !accountHasMember(params.AUID, point.MUID, point.DatabaseID) {
			log.Errorf("BindSwapOwner not owner [token au_id:%v][mu_id:%v][database_id:%v]", params.AUID, point.MUID, point.DatabaseID)
			return base.MakeBaseResponse(resultcode.Result_Error_NotExistMember)
		}
	}
	return nil
}
