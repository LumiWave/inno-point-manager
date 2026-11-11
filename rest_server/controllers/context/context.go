package context

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/auth"
)

// PointManagerContext API의 Request Context
type PointManagerContext struct {
	*base.BaseContext
	VerifyValue *auth.VerifyAuthToken
}

func NewPointManagerContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(PointManagerContext)
	ctx.BaseContext = baseCtx
	ctx.VerifyValue = new(auth.VerifyAuthToken)
	return ctx
}

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
}

func (o *PointManagerContext) SetVerifyAuthToken(value *auth.VerifyAuthToken) {
	o.VerifyValue = value
}

func (o *PointManagerContext) GetValue() *auth.VerifyAuthToken {
	return o.VerifyValue
}

// 인증 서버(inno-auth)가 내려주는 로그인 종류
const (
	AuthLoginType_App        = 2 // 앱(게임 서버) 로그인
	AuthLoginType_AppAccount = 3 // 앱 회원 로그인
	AuthLoginType_WebAccount = 4 // 웹 회원 로그인
)

// IsAppToken 은 앱(게임 서버) 로그인 토큰인지 알려준다.
// 앱 회원 토큰에도 AppID 가 들어 있으므로 로그인 종류까지 확인한다.
func (o *PointManagerContext) IsAppToken() bool {
	v := o.VerifyValue
	return v != nil && v.LoginType == AuthLoginType_App && v.AppID != 0
}

// IsMemberToken 은 회원(앱 회원, 웹 회원) 로그인 토큰인지 알려준다.
func (o *PointManagerContext) IsMemberToken() bool {
	v := o.VerifyValue
	if v == nil || v.AUID == 0 {
		return false
	}
	return v.LoginType == AuthLoginType_AppAccount || v.LoginType == AuthLoginType_WebAccount
}
