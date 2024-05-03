package model

import (
	"strconv"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
)

// redis coin transfer lock key generate
func MakeCoinTransferFromParentWalletLockKey(AUID int64) string {
	return config.GetInstance().DBPrefix + "-COIN-TRANSFER-FROM-PARENT-" + strconv.FormatInt(AUID, 10) + "-lock"
}

// redis coin transfer key generate
func MakeCoinTransferFromParentWalletKey(AUID int64) string {
	return config.GetInstance().DBPrefix + ":COIN-TRANSFER-FROM-PARENT:" + strconv.FormatInt(AUID, 10)
}

// func MakeCoinTransferFromParentWalletKeyByTxID(transactionID string) string {
// 	return config.GetInstance().DBPrefix + ":COIN-TX:" + transactionID
// }

func (o *DB) GetCacheCoinTransferFromParentWallet(key string) (*context.ReqCoinTransferFromParentWallet, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	reqCoinTransfer := new(context.ReqCoinTransferFromParentWallet)

	err := o.Cache.Get(key, reqCoinTransfer)

	return reqCoinTransfer, err
}

func (o *DB) SetCacheCoinTransferFromParentWallet(key string, reqCoinTransfer *context.ReqCoinTransferFromParentWallet) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, reqCoinTransfer, -1)
}

func (o *DB) DelCacheCoinTransferFromParentWallet(key string) error {
	return o.Cache.Del(key)
}
