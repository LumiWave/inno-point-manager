package app

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/externalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/internalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/schedule"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/util"
)

type ServerApp struct {
	base.BaseApp
	conf       *config.ServerConfig
	configFile string

	sysMonitor *schedule.SystemMonitor
}

func (o *ServerApp) Init(configFile string) (err error) {
	o.conf = config.GetInstance(configFile)
	base.AppendReturnCodeText(&resultcode.ResultCodeText)
	context.AppendRequestParameter()

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

func (o *ServerApp) NewDB(conf *config.ServerConfig) error {
	account := conf.MssqlDBAccount
	port, err := strconv.ParseInt(account.Port, 10, 32)
	if err != nil {
		log.Errorf("db port error : %v", port)
		return err
	}
	mssqlDB, err := basedb.GetMssql(account.Database, "", account.ID, account.Password, account.Host, int(port))
	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v, %v, %v",
			err, account.Host, account.ID, account.Password, account.Database, account.PoolSize, account.IdleSize)
		return err
	}

	gCache := basedb.GetCache(&conf.Cache)

	// point db create
	pointDBs := make(map[int64]*basedb.Mssql)

	model.SetDB(mssqlDB, gCache, pointDBs)

	if getPointDBs, err := model.GetDB().GetPointDatabases(); err != nil {
		return err
	} else {
		for _, pointDB := range getPointDBs {

			mssqlDBP, err := basedb.NewMssql(pointDB.DatabaseName,
				"pointDB",
				conf.MssqlDBPoint.ID,
				conf.MssqlDBPoint.Password,
				pointDB.ServerName,
				int(util.ParseInt(conf.MssqlDBPoint.Port)))

			if err != nil {
				log.Errorf("err: %v, val: %v, %v, %v, %v",
					err, pointDB.ServerName, conf.MssqlDBPoint.ID, conf.MssqlDBPoint.Password, pointDB.DatabaseName)
				return err
			}

			pointDBs[pointDB.DatabaseID] = mssqlDBP
		}

		model.SetDBPoint(pointDBs)
	}

	return nil
}
