package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/labstack/echo"
)

type ExternalAPI struct {
	base.BaseController

	conf    *config.ServerConfig
	apiConf *baseconf.APIServer
	echo    *echo.Echo
}

func PreCheck(c echo.Context) base.PreCheckResponse {
	conf := config.GetInstance()
	if err := base.SetContext(c, &conf.Config, context.NewPointManagerContext); err != nil {
		log.Error(err)
		return base.PreCheckResponse{
			IsSucceed: false,
		}
	}

	// auth token 검증

	if conf.Auth.AuthEnable {
		author, ok := c.Request().Header["Authorization"]
		if !ok {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		// if !conf.Auth.InternalAuth {
		// 	walletAddr, isValid := auth.GetIAuth().IsValidAuthToken(author[0][7:])
		// 	if !isValid {
		// 		// auth token 오류 리턴
		// 		res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

		// 		return base.PreCheckResponse{
		// 			IsSucceed: false,
		// 			Response:  res,
		// 		}
		// 	}
		// 	log.Info("from : ", *walletAddr)
		// 	base.GetContext(c).(*context.PointManagerContext).SetWalletAddr(*walletAddr)
		// } else {
		// membership server 인증 진행
		// walletAddr, _, isValid := auth.GetIAuth().GetAuthInfo(author[0][7:])
		// if !isValid {
		// 	// auth token 오류 리턴
		// 	res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

		// 	return base.PreCheckResponse{
		// 		IsSucceed: false,
		// 		Response:  res,
		// 	}
		// }

		if ret, err := auth.CheckAuthToken(author[0][7:]); err != nil || !ret {
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)
			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		//log.Info("from : ", walletAddr)
		//base.GetContext(c).(*context.PointManagerContext).SetWalletAddr(walletAddr)
		//}
	}

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func (o *ExternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	// serving documents for RESTful APIs
	if o.conf.PManager.APIDocs {
		e.Static("/docs", "docs/ext")
	}

	return nil
}

func (o *ExternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[1]
	return o.apiConf
}

func NewAPI() *ExternalAPI {
	return &ExternalAPI{}
}

func (o *ExternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *ExternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}
