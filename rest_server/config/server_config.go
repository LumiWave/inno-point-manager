package config

import (
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
	DBID     int64  `json:"db_id" yaml:"db_id"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	ID       string `json:"id" yaml:"id"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	PoolSize string `json:"poolsize" yaml:"poolsize"`
	IdleSize string `json:"idlesize" yaml:"idlesize"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	ReadPref string `json:"readpref" yaml:"readpref"`
}

type ServerConfig struct {
	baseconf.Config `yaml:",inline"`

	PManager       PointManager    `yaml:"point_manager"`
	MssqlDBAccount baseconf.DBAuth `yaml:"mssql_db_account"`
	MssqlDBPoint   []MssqlPoint    `yaml:"mssql_db_points"`
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
		}
	})

	return currentConfig
}
