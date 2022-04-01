package model

import "github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"

func MakeCoinTransferKeyByTxID(transactionID string) string {
	return config.GetInstance().DBPrefix + ":COIN-TX:" + transactionID
}

func (o *DB) DelCacheCoinTransfer(key string) error {
	return o.Cache.Del(key)
}
