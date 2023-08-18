package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_XchgStrt_Goods                         = "[dbo].[USPAU_XchgStrt_Goods]"
	USPWA_Mod_TransactExchangeGoods_Gasfee       = "[dbo].[USPWA_Mod_TransactExchangeGoods_Gasfee]"
	USPWA_Mod_TransactExchangeGoods_TxStatus     = "[dbo].[USPWA_Mod_TransactExchangeGoods_TxStatus]"
	USPAU_Mod_TransactExchangeGoods_TransactedDT = "[dbo].[USPAU_Mod_TransactExchangeGoods_TransactedDT]"
	USPAU_XchgCmplt_Goods                        = "[dbo].[USPAU_XchgCmplt_Goods]"
)

// 스왑 시작 : 코인 <-> 포인트
func (o *DB) USPAU_XchgStrt_Goods(params *context.ReqSwapInfo) (*int64, error) {
	txID := int64(0)
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_XchgStrt_Goods,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", params.MUID),
		sql.Named("AppID", params.AppID),
		sql.Named("DatabaseID", params.DatabaseID),
		sql.Named("BaseCoinID", params.BaseCoinID),
		sql.Named("WalletAddress", params.WalletAddress),
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
		log.Errorf("USPAU_XchgStrt_Goods QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPAU_XchgStrt_Goods returnvalue error : %v", rs)
		return &txID, errors.New("USPAU_XchgStrt_Goods returnvalue error " + strconv.Itoa(int(rs)))
	}

	return &txID, nil
}

// 가스비 처리
// txStatus 2:수수료 전송 시작, 3:수수료 전송 성공, 4:수수료 전송 실패
func (o *DB) USPWA_Mod_TransactExchangeGoods_Gasfee(txID int64, txStatus int64, txHash, gasFee string) error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPWA_Mod_TransactExchangeGoods_Gasfee,
		sql.Named("TxID", txID),
		sql.Named("TxStatus", txStatus),
		sql.Named("TxHash", txHash),
		sql.Named("Gasfee", gasFee),
		&rs)
	if err != nil {
		log.Errorf("USPWA_Mod_TransactExchangeGoods_Gasfee QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPWA_Mod_TransactExchangeGoods_Gasfee returnvalue error : %v", rs)
		return errors.New("USPWA_Mod_TransactExchangeGoods_Gasfee returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 거래 상태 갱신
func (o *DB) USPWA_Mod_TransactExchangeGoods_TxStatus(txID, txStatus int64) error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPWA_Mod_TransactExchangeGoods_TxStatus,
		sql.Named("TxID", txID),
		sql.Named("TxStatus", txStatus),
		&rs)
	if err != nil {
		log.Errorf("USPWA_Mod_TransactExchangeGoods_TxStatus QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPWA_Mod_TransactExchangeGoods_TxStatus returnvalue error : %v", rs)
		return errors.New("USPWA_Mod_TransactExchangeGoods_TxStatus returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 토큰 처리 : transactionedDT => time.Now().Format("2006-01-02 15:04:05.000")
func (o *DB) USPAU_Mod_TransactExchangeGoods_TransactedDT(txID, txStatus int64, txHash, transactedDT string) error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchangeGoods_TransactedDT,
		sql.Named("TxID", txID),
		sql.Named("TxStatus", txStatus),
		sql.Named("TxHash", txHash),
		sql.Named("TransactedDT", transactedDT),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Mod_TransactExchangeGoods_TransactedDT QueryContext err : %v", err)
		return err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchangeGoods_TransactedDT returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchangeGoods_TransactedDT returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 종료
func (o *DB) USPAU_XchgCmplt_Goods(params *context.ReqSwapInfo, completedDT string, isSuccess bool) error {
	var rs orginMssql.ReturnStatus
	if params.TxType == context.EventID_toCoin {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_XchgCmplt_Goods,
				sql.Named("TxID", params.TxID),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_XchgCmplt_Goods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_XchgCmplt_Goods,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.PointID),
				sql.Named("PointPreQuantity", params.PreviousPointQuantity),
				sql.Named("PointAdjQuantity", params.AdjustPointQuantity),
				sql.Named("PointQuantity", params.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_XchgCmplt_Goods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	} else if params.TxType == context.EventID_toPoint {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_XchgCmplt_Goods,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.PointID),
				sql.Named("PointPreQuantity", params.PreviousPointQuantity),
				sql.Named("PointAdjQuantity", params.AdjustPointQuantity),
				sql.Named("PointQuantity", params.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_XchgCmplt_Goods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_XchgCmplt_Goods,
				sql.Named("TxID", params.TxID),
				&rs)
			if err != nil {
				log.Errorf("USPAU_XchgCmplt_Goods QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	}
	if rs != 1 {
		log.Errorf("USPAU_XchgCmplt_Goods returnvalue error : %v", rs)
		return errors.New("USPAU_XchgCmplt_Goods returnvalue error " + strconv.Itoa(int(rs)))
	}
	return nil
}
