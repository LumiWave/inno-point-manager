package model

import (
	originCtx "context"
	"database/sql"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/api_inno_log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Exchange_Goods = "[dbo].[USPAU_Exchange_Goods]"
)

// 지갑 정보 조회
func (o *DB) PostPointCoinSwap(params *context.ReqSwapInfo, txHash string) error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Exchange_Goods,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", params.MUID),
		sql.Named("AppID", params.AppID),
		sql.Named("DatabaseID", params.DatabaseID),
		sql.Named("CoinID", params.CoinID),
		sql.Named("BaseCoinID", params.BaseCoinID),
		sql.Named("WalletAddress", params.WalletAddress),
		sql.Named("PreCoinQuantity", params.PreviousCoinQuantity),
		sql.Named("AdjCoinQuantity", params.AdjustCoinQuantity),
		sql.Named("CoinQuantity", params.CoinQuantity),
		sql.Named("PointID", params.PointID),
		sql.Named("PrePointQuantity", params.PreviousPointQuantity),
		sql.Named("AdjPointQuantity", params.AdjustPointQuantity),
		sql.Named("PointQuantity", params.PointQuantity),
		sql.Named("LogID", params.LogID),
		sql.Named("EventID", params.EventID),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Exchange_Goods QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	apiParams := &api_inno_log.ExchangeGoodsLog{
		LogDt:            time.Now().Format("2006-01-02 15:04:05.000"),
		LogID:            int64(params.LogID),
		EventID:          int64(params.EventID),
		TxHash:           txHash,
		AUID:             params.AUID,
		InnoUID:          params.InnoUID,
		MUID:             params.MUID,
		AppID:            params.AppID,
		CoinID:           params.CoinID,
		BaseCoinID:       params.BaseCoinID,
		WalletAddress:    params.WalletAddress,
		PreCoinQuantity:  params.PreviousCoinQuantity,
		AdjCoinQuantity:  params.AdjustCoinQuantity,
		CoinQuantity:     params.CoinQuantity,
		PointID:          params.PointID,
		PrePointQuantity: params.PreviousPointQuantity,
		AdjPointQuantity: params.AdjustPointQuantity,
		PointQuantity:    params.PointQuantity,
	}
	go api_inno_log.GetInstance().PostExchangeGoods(apiParams)

	return nil
}
