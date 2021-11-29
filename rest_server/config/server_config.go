package config

import (
	"sync"

	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
)

var once sync.Once
var currentConfig *ServerConfig

type PointManager struct {
	ApplicationName string `json:"application_name" yaml:"application_name"`
	APIDocs         bool   `json:"api_docs" yaml:"api_docs"`
}

type TokenInfo struct {
	MainnetHost      string   `yaml:"mainnet_host"`
	ServerWalletAddr string   `yaml:"server_wallet_address"`
	ServerPrivateKey string   `yaml:"server_private_key"`
	TokenAddrs       []string `yaml:"token_address"`
	NftUriDomain     string   `yaml:"nft_uri_domain"`
}

type ApiAuth struct {
	AuthEnable        bool   `yaml:"auth_enable"`
	JwtSecretKey      string `yaml:"jwt_secret_key"`
	TokenExpiryPeriod int64  `yaml:"token_expiry_period"`
	SignExpiryPeriod  int64  `yaml:"sign_expiry_period"`
	AesKey            string `yaml:"aes_key"`
	InternalAuth      bool   `json:"internal_auth" yaml:"internal_auth"`
	ApiAuthDomain     string `json:"api_auth_domain" yaml:"api_auth_domain"`
	ApiAuthVerify     string `json:"api_auth_verify" yaml:"api_auth_verify"`
}

type Azure struct {
	AzureStorageAccount   string `yaml:"azure_storage_account"`
	AzureStorageAccessKey string `yaml:"azure_storage_access_key"`
	Domain                string `yaml:"azure_storage_domain"`
	ContainerNft          string `yaml:"azure_container_nft_folder"`
	ContainerProduct      string `yaml:"azure_container_product_folder"`
}

type ServerConfig struct {
	baseconf.Config `yaml:",inline"`

	PManager    PointManager    `yaml:"point_manager"`
	MssqlDBAuth baseconf.DBAuth `yaml:"mssql_db_auth"`
	Token       TokenInfo       `yaml:"token_info"`
	Auth        ApiAuth         `yaml:"api_auth"`
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
