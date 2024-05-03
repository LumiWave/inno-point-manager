package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
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
		model.GetDB().PointDocMtx.Lock()
		resp.Value = model.GetDB().PointDoc
		model.GetDB().PointDocMtx.Unlock()
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
