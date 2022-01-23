package token_manager_server

var gTokenManagerServerInfo *TokenManagerServerInfo

type HostInfo struct {
	IntHostUri string
	ExtHostUri string
	IntVer     string // m1.0
	ExtVer     string // v1.0
}

type AuthInfo struct {
	ApiKey string
}

type TokenManagerServerInfo struct {
	HostInfo

	AuthInfo
}

func GetInstance() *TokenManagerServerInfo {
	return gTokenManagerServerInfo
}

func NewTokenManagerServerInfo(apiKey string, hostInfo HostInfo) *TokenManagerServerInfo {
	if gTokenManagerServerInfo == nil {
		gTokenManagerServerInfo = &TokenManagerServerInfo{
			HostInfo: hostInfo,
			AuthInfo: AuthInfo{
				ApiKey: apiKey,
			},
		}
	}

	return gTokenManagerServerInfo
}

func (o *TokenManagerServerInfo) SetApiKey(key string) {
	gTokenManagerServerInfo.ApiKey = key
}
