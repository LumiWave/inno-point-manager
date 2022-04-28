package model

import (
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

// redis coin transfer lock key generate
func MakeCoinTransferToParentWalletLockKey(AUID int64) string {
	return config.GetInstance().DBPrefix + "-COIN-TRANSFER-TO-PARENT-" + strconv.FormatInt(AUID, 10) + "-lock"
}

// redis coin transfer key generate
func MakeCoinTransferToParentWalletKey(AUID int64) string {
	return config.GetInstance().DBPrefix + ":COIN-TRANSFER-TO-PARENT:" + strconv.FormatInt(AUID, 10)
}

func (o *DB) GetCacheCoinTransferToParentWallet(key string) (*context.ReqCoinTransferToParentWallet, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	reqCoinTransfer := new(context.ReqCoinTransferToParentWallet)

	err := o.Cache.Get(key, reqCoinTransfer)

	return reqCoinTransfer, err
}

func (o *DB) SetCacheCoinTransferToParentWallet(key string, reqCoinTransfer *context.ReqCoinTransferToParentWallet) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, reqCoinTransfer, -1)
}
