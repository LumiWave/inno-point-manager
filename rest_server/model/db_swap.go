package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/api_inno_log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Strt_ExchangeGoods                    = "[dbo].[USPAU_Strt_ExchangeGoods]"
	USPAU_Mod_TransactExchangeGoods_Exchangefee = "[dbo].[USPAU_Mod_TransactExchangeGoods_Exchangefee]"
	USPAU_Mod_TransactExchangeGoods_TxStatus    = "[dbo].[USPAU_Mod_TransactExchangeGoods_TxStatus]"
	USPAU_Mod_TransactExchangeGoods_Coin        = "[dbo].[USPAU_Mod_TransactExchangeGoods_Coin]"
	USPAU_Cmplt_ExchangeGoods                   = "[dbo].[USPAU_Cmplt_ExchangeGoods]"
)

// 스왑 시작 : 코인 <-> 포인트
func (o *DB) USPAU_Strt_ExchangeGoods(params *context.ReqSwapInfo) (*int64, error) {
	txID := int64(0)
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Strt_ExchangeGoods,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", params.MUID),
		sql.Named("AppID", params.AppID),
		sql.Named("DatabaseID", params.DatabaseID),
		sql.Named("BaseCoinID", params.BaseCoinID),
		sql.Named("WalletAddress", params.WalletAddress),
		sql.Named("WalletTypeID", params.WalletTypeID),
		sql.Named("CoinID", params.CoinID),
		sql.Named("CoinAdjQuantity", params.AdjustCoinQuantity),
		sql.Named("PointID", params.PointID),
		sql.Named("PointPreQuantity", params.PreviousPointQuantity),
		sql.Named("PointAdjQuantity", params.AdjustPointQuantity),
		sql.Named("PointQuantity", params.PointQuantity),
		sql.Named("TxType", params.TxType),
		sql.Named("TxID", sql.Out{Dest: &txID}),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Strt_ExchangeGoods QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPAU_Strt_ExchangeGoods returnvalue error : %v", rs)
		return &txID, errors.New("USPAU_Strt_ExchangeGoods returnvalue error " + strconv.Itoa(int(rs)))
	}

	return &txID, nil
}

// 가스비 처리
// txStatus 2:수수료 전송 시작, 3:수수료 전송 성공, 4:수수료 전송 실패
func (o *DB) USPAU_Mod_TransactExchangeGoods_Exchangefee(txID int64, txStatus int64, txHash, swapFee string, baseCoinID int64, gasFee string) error {
	var rs orginMssql.ReturnStatus
	if txStatus == context.SWAP_status_fee_transfer_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_Exchangefee,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("ExchangeFee", swapFee),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("Gasfee", gasFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_Exchangefee QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_Exchangefee,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("ExchangeFee", swapFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_Exchangefee QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchangeGoods_Exchangefee returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchangeGoods_Exchangefee returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 거래 상태 갱신
func (o *DB) USPAU_Mod_TransactExchangeGoods_TxStatus(baseCoinID int64, gasFee string, swapInfo *context.ReqSwapInfo) error {
	var rs orginMssql.ReturnStatus
	if swapInfo.TxStatus == context.SWAP_status_fee_transfer_success || swapInfo.TxStatus == context.SWAP_status_token_transfer_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_TxStatus,
			sql.Named("TxID", swapInfo.TxID),
			sql.Named("TxStatus", swapInfo.TxStatus),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("Gasfee", gasFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_TxStatus QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
		if swapInfo.TxStatus == context.SWAP_status_token_transfer_success {

			strAdjCoinQuantity := strconv.FormatFloat(swapInfo.AdjustCoinQuantity, 'f', -1, 64)

			apiParams := &api_inno_log.ExchangeGoodsLog{
				LogDt:            time.Now().Format("2006-01-02 15:04:05.000"),
				LogID:            context.LogID_exchange,
				EventID:          swapInfo.TxType,
				TxHash:           swapInfo.TokenTxHash,
				TxID:             swapInfo.TxID,
				AUID:             swapInfo.AUID,
				InnoUID:          swapInfo.InnoUID,
				MUID:             swapInfo.MUID,
				AppID:            swapInfo.AppID,
				CoinID:           swapInfo.CoinID,
				BaseCoinID:       baseCoinID,
				WalletAddress:    swapInfo.WalletAddress,
				AdjCoinQuantity:  strAdjCoinQuantity,
				PointID:          swapInfo.PointID,
				AdjPointQuantity: swapInfo.AdjustPointQuantity,
				WalletID:         swapInfo.WalletID,
			}
			log.Debugf("log : %v", apiParams)
			go api_inno_log.GetInstance().PostExchangeGoods(apiParams)
		}
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_TxStatus,
			sql.Named("TxID", swapInfo.TxID),
			sql.Named("TxStatus", swapInfo.TxStatus),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_TxStatus QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchangeGoods_TxStatus returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchangeGoods_TxStatus returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 토큰 처리 : transactionedDT => time.Now().Format("2006-01-02 15:04:05.000")
func (o *DB) USPAU_Mod_TransactExchangeGoods_Coin(txID, txStatus int64, txHash, transactedDT string, baseCoinID int64, gasFee string) error {
	var rs orginMssql.ReturnStatus
	if txStatus == context.SWAP_status_token_transfer_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_Coin,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("TransactedDT", transactedDT),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("Gasfee", gasFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_Coin QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_Coin,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("TransactedDT", transactedDT),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchangeGoods_Coin QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchangeGoods_Coin returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchangeGoods_Coin returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 종료
// point -> coin, coin->poin, 성공, 실패에 따라 인자값이 달라진다.
func (o *DB) USPAU_Cmplt_ExchangeGoods(params *context.ReqSwapInfo, completedDT string, isSuccess bool) error {
	var rs orginMssql.ReturnStatus
	if params.TxType == context.EventID_toCoin {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_ExchangeGoods,
				sql.Named("TxID", params.TxID),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_ExchangeGoods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_ExchangeGoods,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.PointID),
				sql.Named("PointPreQuantity", params.PreviousPointQuantity),
				sql.Named("PointAdjQuantity", params.AdjustPointQuantity),
				sql.Named("PointQuantity", params.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_ExchangeGoods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	} else if params.TxType == context.EventID_toPoint {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_ExchangeGoods,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.PointID),
				sql.Named("PointPreQuantity", params.PreviousPointQuantity),
				sql.Named("PointAdjQuantity", params.AdjustPointQuantity),
				sql.Named("PointQuantity", params.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_ExchangeGoods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_ExchangeGoods,
				sql.Named("TxID", params.TxID),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_ExchangeGoods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	}
	if rs != 1 {
		log.Errorf("USPAU_Cmplt_ExchangeGoods returnvalue error : %v", rs)
		return errors.New("USPAU_Cmplt_ExchangeGoods returnvalue error " + strconv.Itoa(int(rs)))
	}
	return nil
}
