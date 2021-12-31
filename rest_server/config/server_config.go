package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
)

var once sync.Once
var currentConfig *ServerConfig

type PointManager struct {
	ApplicationName        string `json:"application_name" yaml:"application_name"`
	APIDocs                bool   `json:"api_docs" yaml:"api_docs"`
	CachePointExpiryPeriod int64  `json:"cache_point_expiry_period" yaml:"cache_point_expiry_period"`
}

type ApiAuth struct {
	AuthEnable    bool   `yaml:"auth_enable"`
	ApiAuthDomain string `json:"api_auth_domain" yaml:"api_auth_domain"`
	ApiAuthVerify string `json:"api_auth_verify" yaml:"api_auth_verify"`
}

type MssqlPoint struct {
	Port     string `json:"port" yaml:"port"`
	ID       string `json:"id" yaml:"id"`
	Password string `json:"password" yaml:"password"`
}

type ServerConfig struct {
	baseconf.Config `yaml:",inline"`

	PManager       PointManager    `yaml:"point_manager"`
	MssqlDBAccount baseconf.DBAuth `yaml:"mssql_db_account"`
	MssqlDBPoint   MssqlPoint      `yaml:"mssql_db_point"`
	Auth           ApiAuth         `yaml:"api_auth"`
}

func GetInstance(filepath ...string) *ServerConfig {
	once.Do(func() {
		if len(filepath) <= 0 {
			panic(baseconf.ErrInitConfigFailed)
		}
		currentConfig = &ServerConfig{}
		if err := baseconf.Load(filepath[0], currentConfig); err != nil {
			currentConfig = nil
		} else {
			if os.Getenv("ASPNETCORE_PORT") != "" {
				port, _ := strconv.ParseInt(os.Getenv("ASPNETCORE_PORT"), 10, 32)
				currentConfig.APIServers[0].Port = int(port)
				currentConfig.APIServers[1].Port = int(port)
				fmt.Println(port)
			}
		}
	})

	return currentConfig
}
