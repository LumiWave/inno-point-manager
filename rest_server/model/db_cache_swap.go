package model

import (
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

// redis coin transfer lock key generate
func MakeSwapLockKey(AUID int64) string {
	return config.GetInstance().DBPrefix + "-SWAP" + strconv.FormatInt(AUID, 10) + "-lock"
}

// redis coin transfer key generate
// func MakeSwapKey(AUID int64) string {
// 	return config.GetInstance().DBPrefix + ":SWAP:" + strconv.FormatInt(AUID, 10)
// }

func MakeSwapKey(walletAddr string) string {
	return config.GetInstance().DBPrefix + ":SWAP:" + walletAddr
}

func (o *DB) GetCacheSwapInfo(key string) (*context.ReqSwapInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	swapInfo := new(context.ReqSwapInfo)

	err := o.Cache.Get(key, swapInfo)

	return swapInfo, err
}

func (o *DB) SetCacheSwapInfo(key string, req *context.ReqSwapInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, req, -1)
}

func (o *DB) DelCacheSwapInfo(key string) error {
	return o.Cache.Del(key)
}
