package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func PostPointMemberRegister(req *context.ReqPointMemberRegister, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := model.GetDB().InsertPointMember(req); err != nil {
		model.MakeDbError(resp, resultcode.Result_DBError, err)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PutPointMemberUpdate(params *context.PointMemberInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := model.GetDB().UpdatePointMember(params); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetPointMember(params *context.PointMemberInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if value, err := model.GetDB().SelectPointMember(params.CpMemberIdx); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Value = value
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
