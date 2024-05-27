package context

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/shirou/gopsutil/disk"
)

type SystemRedisRemove struct {
	AuctionList string `query:"auction_list"`
	ProductList string `query:"product_list"`
	BidList     string `query:"bid_list"`
	AuctionId   string `query:"auc_id"`
}

func NewSystemRedisRemove() *SystemRedisRemove {
	return new(SystemRedisRemove)
}

func (o *SystemRedisRemove) CheckValidate() *base.BaseResponse {

	return nil
}

type DiskUsage struct {
	Disk disk.UsageStat
}

type NodeMetric struct {
	Host string `json:"host"`

	Version       string      `json:"version"`
	IsRunning     bool        `json:"is_running"`
	UpTime        string      `json:"up_time"`
	CpuTime       string      `json:"cpu_time"`
	MemTotalBytes uint64      `json:"mem_total_bytes"`
	MemAllocBytes uint64      `json:"mem_alloc_bytes"`
	MemPercent    float32     `json:"mem_usage_percent"`
	CpuUsage      int32       `json:"cpu_usage"`
	DiskUsage     []DiskUsage `json:"disk_usage"`
}

type PSMaintenance struct {
	Enable bool `json:"enable"`
}

func NewPSMaintenance() *PSMaintenance {
	return new(PSMaintenance)
}

type PSSwap struct {
	Enable bool `json:"enable"`
}

func NewPSSwap() *PSSwap {
	return new(PSSwap)
}

type PSCoinTransferExternal struct {
	Enable bool `json:"enable"`
}

func NewPSCoinTransferExternal() *PSCoinTransferExternal {
	return new(PSCoinTransferExternal)
}

type PSPointUpdate struct {
	Enable bool `json:"enable"`
}

func NewPSPointUpdate() *PSPointUpdate {
	return new(PSPointUpdate)
}

type Meta struct {
	PointList interface{} `json:"point_list"`
	AppCoins  interface{} `json:"app_coins"`
	Coins     interface{} `json:"coins"`
}
