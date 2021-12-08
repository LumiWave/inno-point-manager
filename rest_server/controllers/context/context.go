package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/auth"
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
