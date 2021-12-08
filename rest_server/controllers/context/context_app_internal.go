package context

import "github.com/ONBUFF-IP-TOKEN/baseapp/base"

///////// ap 포인트 처리 모니터링
type ReqPointAppMonitoring struct {
	MUID int64 `query:"mu_id"`
}

func NewReqPointAppMonitoring() *ReqPointAppMonitoring {
	return new(ReqPointAppMonitoring)
}

func (o *ReqPointAppMonitoring) CheckValidate() *base.BaseResponse {
	return nil
}
