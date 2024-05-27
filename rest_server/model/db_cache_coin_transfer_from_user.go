package model

import (
	"strconv"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
)

// redis coin transfer from user walllet lock key generate
func MakeCoinTransferFromUserWalletLockKey(AUID int64) string {
	return config.GetInstance().DBPrefix + "-COIN-TRANSFER-FROM-USER-" + strconv.FormatInt(AUID, 10) + "-lock"
}

// redis coin transfer key generate
func MakeCoinTransferFromUserWalletKey(AUID int64) string {
	return config.GetInstance().DBPrefix + ":COIN-TRANSFER-FROM-USER:" + strconv.FormatInt(AUID, 10)
}

// func MakeCoinTransferFromUserWalletKeyByTxID(transactionID string) string {
// 	return config.GetInstance().DBPrefix + ":COIN-TX:" + transactionID
// }

func (o *DB) GetCacheCoinTransferFromUserWallet(key string) (*context.ReqCoinTransferFromUserWallet, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	reqCoinTransfer := new(context.ReqCoinTransferFromUserWallet)

	err := o.Cache.Get(key, reqCoinTransfer)

	return reqCoinTransfer, err
}

func (o *DB) SetCacheCoinTransferFromUserWallet(key string, reqCoinTransfer *context.ReqCoinTransferFromUserWallet) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, reqCoinTransfer, -1)
}

func (o *DB) DelCacheCoinTransferFromUserWallet(key string) error {
	return o.Cache.Del(key)
}
