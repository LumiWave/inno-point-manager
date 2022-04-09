package model

import (
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

func MakeCoinTransferKeyByTxID(transactionID string) string {
	return config.GetInstance().DBPrefix + ":COIN-TX:" + transactionID
}

func (o *DB) DelCacheCoinTransfer(key string) error {
	return o.Cache.Del(key)
}

func (o *DB) GetCacheCoinTransferTx(key string) (*context.TxType, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	txType := new(context.TxType)

	err := o.Cache.Get(key, txType)

	return txType, err
}

func (o *DB) SetCacheCoinTransferTx(key string, txType *context.TxType) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(key, txType, -1)
}
