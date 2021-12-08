package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func GetPointAppMonitoring(req *context.ReqPointAppMonitoring, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if req.MUID != 0 {
		for _, memberPointInfo := range model.GetDB().PointDoc {
			if req.MUID == memberPointInfo.MUID {
				resp.Value = memberPointInfo
				break
			}
		}
	} else {
		resp.Value = model.GetDB().PointDoc
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
