package model

import "fmt"

const (
	PubSub      = "pubsub"
	InternalCmd = "internal_cmd"
)

const (
	PubSub_type_healthcheck          = "HealthCheck_InnoPoint"
	PubSub_type_maintenance          = "Maintenance"
	PubSub_type_Swap                 = "Swap"
	PubSub_type_CoinTransferExternal = "CoinTransferExternal"
	PubSub_type_meta_refresh         = "MetaRefresh"
	PubSub_type_point_update         = "PointUpdate"
)

type PSHeader struct {
	Type string `json:"type"`
}

type PSHealthCheck struct {
	PSHeader
	Value struct {
		Timestamp int64 `json:"ts"`
	} `json:"value"`
}

type PSMaintenance struct {
	PSHeader
	Value struct {
		Enable    bool   `json:"enable"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"value"`
}

type PSSwap struct {
	PSHeader
	Value struct {
		Enable bool `json:"enable"`
	} `json:"value"`
}

type PSCoinTransferExternal struct {
	PSHeader
	Value struct {
		Enable bool `json:"enable"`
	} `json:"value"`
}

type PSMetaRefresh struct {
	PSHeader
	Value struct {
		Enable bool `json:"enable"`
	} `json:"value"`
}

type PSPointUpdate struct {
	PSHeader
	Value struct {
		Enable bool `json:"enable"`
	} `json:"value"`
}

func MakePubSubKey(val string) string {
	return fmt.Sprintf("%s:%s", PubSub, val)
}
