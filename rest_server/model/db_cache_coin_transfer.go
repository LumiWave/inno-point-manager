package model

import (
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

// redis coin transfer lock key generate
func MakeCoinTransferLockKey(AUID int64) string {
	return config.GetInstance().DBPrefix + "-COIN-TRANSFER-" + strconv.FormatInt(AUID, 10) + "-lock"
}

// redis coin transfer key generate
func MakeCoinTransferKey(AUID int64) string {
	return config.GetInstance().DBPrefix + ":COIN-TRANSFER:" + strconv.FormatInt(AUID, 10)
}

func MakeCoinTransferKeyByTxID(transactionID string) string {
	return config.GetInstance().DBPrefix + ":COIN-TX:" + transactionID
}

func (o *DB) GetCacheCoinTransfer(key string) (*context.ReqCoinTransfer, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	reqCoinTransfer := new(context.ReqCoinTransfer)

	err := o.Cache.Get(key, reqCoinTransfer)

	return reqCoinTransfer, err
}

func (o *DB) SetCacheCoinTransfer(key string, reqCoinTransfer *context.ReqCoinTransfer) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, reqCoinTransfer, -1)
}

func (o *DB) DelCacheCoinTransfer(key string) error {
	return o.Cache.Del(key)
}
