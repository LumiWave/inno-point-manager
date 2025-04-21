package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPPR_Scan_PreSalesExchangeCoinToPoints      = "[dbo].[USPPR_Scan_PreSalesExchangeCoinToPoints]"
	USPPR_Strt_PreSalesExchanges                 = "[dbo].[USPPR_Strt_PreSalesExchanges]"
	USPPR_Mod_TransactPreSalesExchanges_Coin     = "[dbo].[USPPR_Mod_TransactPreSalesExchanges_Coin]"
	USPPR_Mod_TransactPreSalesExchanges_TxStatus = "[dbo].[USPPR_Mod_TransactPreSalesExchanges_TxStatus]"
	USPPR_Cmplt_PreSalesExchanges                = "[dbo].[USPPR_Cmplt_PreSalesExchanges]"
)

func (o *DB) USPPR_Scan_PreSalesExchangeCoinToPoints() ([]*context.PreSalesExchange, error) {
	ProcName := USPPR_Scan_PreSalesExchangeCoinToPoints
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlPreSales.GetDB().QueryContext(originCtx.Background(), ProcName,
		&rs)
	if err != nil {
		log.Errorf(ProcName+" QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	o.SwapAblePreSalesMap = make(map[int64]map[int64]*context.PreSalesExchange)

	o.SwapAblePreSales = make([]*context.PreSalesExchange, 0)
	for rows.Next() {
		item := &context.PreSalesExchange{}
		if err := rows.Scan(&item.SalesID, &item.BaseCoinID, &item.CoinID, &item.PointID, &item.MinimumExchangeQuantity, &item.ExchangeRatio); err != nil {
			log.Errorf(ProcName+" Get error : %v", err)
			return nil, err
		} else {
			o.SwapAblePreSales = append(o.SwapAblePreSales, item)

			if o.SwapAblePreSalesMap[item.CoinID] == nil {
				o.SwapAblePreSalesMap[item.CoinID] = make(map[int64]*context.PreSalesExchange)
			}
			o.SwapAblePreSalesMap[item.CoinID][item.PointID] = item
		}
	}

	if rs != 1 {
		log.Errorf(ProcName+" returnvalue error : %v", rs)
		return nil, errors.New(ProcName + " returnvalue error " + strconv.Itoa(int(rs)))
	}
	return o.SwapAblePreSales, nil
}

func (o *DB) USPPR_Strt_PreSalesExchanges(salesID int64, auid, muid, baseCoinID, walletTypeID, walletID int64, walletAddress string,
	coinID int64, adjCoinQuantity string, pointID int64, adjPointQuantity string) (int64, error) {
	ProcName := USPPR_Strt_PreSalesExchanges
	var rs orginMssql.ReturnStatus
	txID := int64(0)
	rows, err := o.MssqlPreSales.GetDB().QueryContext(originCtx.Background(), ProcName,
		sql.Named("SalesID", salesID),
		sql.Named("AUID", auid),
		sql.Named("MUID", muid),
		sql.Named("BaseCoinID", baseCoinID),
		sql.Named("WalletTypeID", walletTypeID),
		sql.Named("WalletID", walletID),
		sql.Named("WalletAddress", walletAddress),
		sql.Named("CoinID", coinID),
		sql.Named("AdjCoinQuantity", adjCoinQuantity), // 프리세일은 c to p 이니 무조건 음수 부호가 있어야됨
		sql.Named("PointID", pointID),
		sql.Named("AdjPointQuantity", adjPointQuantity),
		sql.Named("TxID", sql.Out{Dest: &txID}),
		&rs)
	if err != nil {
		log.Errorf(ProcName+" QueryContext err : %v", err)
		return txID, err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf(ProcName+" returnvalue error : %v", rs)
		return txID, errors.New(ProcName + " returnvalue error " + strconv.Itoa(int(rs)))
	}

	return txID, nil
}

func (o *DB) USPPR_Mod_TransactPreSalesExchanges_Coin(txID int64, txStatus int64, txHash string, baseCoinID int64, gasFee string) error {
	ProcName := USPPR_Mod_TransactPreSalesExchanges_Coin
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlPreSales.GetDB().QueryContext(originCtx.Background(), ProcName,
		sql.Named("TxID", txID),
		sql.Named("TxStatus", txStatus),
		sql.Named("TxHash", txHash),
		sql.Named("TransactedDT", time.Now().Format("2006-01-02 15:04:05.000")),
		sql.Named("BaseCoinID", baseCoinID),
		sql.Named("GasFee", gasFee),
		&rs)
	if err != nil {
		log.Errorf(ProcName+" QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf(ProcName+" returnvalue error : %v", rs)
		return errors.New(ProcName + " returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

func (o *DB) USPPR_Mod_TransactPreSalesExchanges_TxStatus(txID int64, txStatus int64, baseCoinID int64, gasFee string) error {
	ProcName := USPPR_Mod_TransactPreSalesExchanges_TxStatus
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlPreSales.GetDB().QueryContext(originCtx.Background(), ProcName,
		sql.Named("TxID", txID),
		sql.Named("TxStatus", txStatus),
		sql.Named("BaseCoinID", baseCoinID),
		sql.Named("GasFee", gasFee),
		&rs)
	if err != nil {
		log.Errorf(ProcName+" QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf(ProcName+" returnvalue error : %v", rs)
		return errors.New(ProcName + " returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

func (o *DB) USPPR_Cmplt_PreSalesExchanges(txID int64, pointID int64, prePointQuantity int64, adjPointQuantity int64, pointQuantity int64) error {
	ProcName := USPPR_Cmplt_PreSalesExchanges
	var rs orginMssql.ReturnStatus

	rows, err := o.MssqlPreSales.GetDB().QueryContext(originCtx.Background(), ProcName,
		sql.Named("TxID", txID),
		sql.Named("PointID", pointID),
		sql.Named("PrePointQuantity", prePointQuantity),
		sql.Named("AdjPointQuantity", adjPointQuantity),
		sql.Named("PointQuantity", pointQuantity),
		sql.Named("CompletedDT", time.Now().Format("2006-01-02 15:04:05.000")),
		&rs)
	if err != nil {
		log.Errorf(ProcName+" QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf(ProcName+" returnvalue error : %v", rs)
		return errors.New(ProcName + " returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}
