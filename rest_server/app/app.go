package app

import (
	"fmt"
	"sync"

	"github.com/LumiWave/baseapp/base"
	baseconf "github.com/LumiWave/baseapp/config"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/api_inno_log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/auth"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/externalapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/internalapi"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/token_manager_server"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
	"github.com/LumiWave/inno-point-manager/rest_server/schedule"
)

type ServerApp struct {
	base.BaseApp
	conf       *config.ServerConfig
	configFile string

	sysMonitor          *schedule.SystemMonitor
	swapExpireScheduler *schedule.SwapExpireScheduler
}

func (o *ServerApp) Init(configFile string) (err error) {
	o.conf = config.GetInstance(configFile)
	base.AppendReturnCodeText(&resultcode.ResultCodeText)
	context.AppendRequestParameter()

	if err := o.InitScheduler(); err != nil {
		return err
	}
	auth.InitHttpClient()
	o.InitTokenManagerServer(o.conf)
	o.InitLogServer(o.conf)

	if err := o.NewDB(o.conf); err != nil {
		return err
	}

	return err
}

func (o *ServerApp) CleanUp() {
	fmt.Println("CleanUp")
}

func (o *ServerApp) Run(wg *sync.WaitGroup) error {
	return nil
}

func (o *ServerApp) GetConfig() *baseconf.Config {
	return &o.conf.Config
}

func NewApp() (*ServerApp, error) {
	app := &ServerApp{}

	intAPI := internalapi.NewAPI()
	extAPI := externalapi.NewAPI()

	if err := app.BaseApp.Init(app, intAPI, extAPI); err != nil {
		return nil, err
	}

	return app, nil
}

func (o *ServerApp) InitScheduler() error {
	//o.sysMonitor = schedule.GetSystemMonitor()
	o.swapExpireScheduler = schedule.InitSwapExpireScheduler(o.conf)
	return nil
}

func (o *ServerApp) InitTokenManagerServer(conf *config.ServerConfig) {
	pointMgrServer := conf.TokenMgrServer
	hostInfo := token_manager_server.HostInfo{
		IntHostUri: pointMgrServer.InternalpiDomain,
		ExtHostUri: pointMgrServer.ExternalpiDomain,
		IntVer:     pointMgrServer.InternalVer,
		ExtVer:     pointMgrServer.ExternalVer,
	}
	token_manager_server.NewTokenManagerServerInfo("", hostInfo)
}

func (o *ServerApp) NewDB(conf *config.ServerConfig) error {
	return model.InitDB(conf)
}

func (o *ServerApp) InitLogServer(conf *config.ServerConfig) {
	confLog := conf.InnoLog
	hostInfo := api_inno_log.HostInfo{
		IntHostUri: confLog.InternalpiDomain,
		ExtHostUri: confLog.ExternalpiDomain,
		IntVer:     confLog.InternalVer,
		ExtVer:     confLog.ExternalVer,
	}
	api_inno_log.NewServerInfo("", hostInfo)
}
