package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Strt_Exchanges                     = "[dbo].[USPAU_Strt_Exchanges]"
	USPAU_Mod_TransactExchanges_ExchangeFees = "[dbo].[USPAU_Mod_TransactExchanges_ExchangeFees]"
	USPAU_Mod_TransactExchanges_TxStatus     = "[dbo].[USPAU_Mod_TransactExchanges_TxStatus]"
	USPAU_Mod_TransactExchanges_Coin         = "[dbo].[USPAU_Mod_TransactExchanges_Coin]"
	USPAU_Cmplt_Exchanges                    = "[dbo].[USPAU_Cmplt_Exchanges]"

	USPAU_Scan_ExchangeCoinToCoins  = "[dbo].[USPAU_Scan_ExchangeCoinToCoins]"
	USPAU_Scan_ExchangePointToCoins = "[dbo].[USPAU_Scan_ExchangePointToCoins]"
	USPAU_Scan_ExchangeCoinToPoints = "[dbo].[USPAU_Scan_ExchangeCoinToPoints]"
)

// 스왑 시작 : 코인 <-> 포인트
func (o *DB) USPAU_Strt_Exchanges(params *context.ReqSwapInfo) (*int64, error) {
	txID := int64(0)
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Strt_Exchanges,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", func() int64 {
			if params.TxType == context.EventID_P2C {
				return params.SwapFromPoint.MUID
			} else if params.TxType == context.EventID_C2P {
				return params.SwapToPoint.MUID
			}
			return 0
		}()),
		sql.Named("DatabaseID", func() int64 {
			if params.TxType == context.EventID_P2C {
				return params.SwapFromPoint.DatabaseID
			} else if params.TxType == context.EventID_C2P {
				return params.SwapToPoint.DatabaseID
			}
			return 0
		}()),
		sql.Named("TxType", params.TxType),
		sql.Named("FromBaseCoinID", params.SwapFromCoin.BaseCoinID),
		sql.Named("FromWalletTypeID", params.SwapFromCoin.WalletTypeID),
		sql.Named("FromWalletID", params.SwapFromCoin.WalletID),
		sql.Named("FromWalletAddress", params.SwapFromCoin.WalletAddress),
		sql.Named("FromID", func() int64 {
			if params.TxType == context.EventID_P2C {
				return params.SwapFromPoint.PointID
			} else if params.TxType == context.EventID_C2P || params.TxType == context.EventID_C2C {
				return params.SwapFromCoin.CoinID
			}
			return 0
		}()),
		sql.Named("FromAdjQuantity", func() string {
			if params.TxType == context.EventID_P2C {
				return strconv.FormatInt(params.SwapFromPoint.AdjustPointQuantity, 10)
			} else if params.TxType == context.EventID_C2P || params.TxType == context.EventID_C2C {
				return strconv.FormatFloat(params.SwapFromCoin.AdjustCoinQuantity, 'f', -1, 64)
			}
			return ""
		}()),
		sql.Named("ToBaseCoinID", params.SwapToCoin.BaseCoinID),
		sql.Named("ToWalletTypeID", params.SwapToCoin.WalletTypeID),
		sql.Named("ToWalletID", params.SwapToCoin.WalletID),
		sql.Named("ToWalletAddress", params.SwapToCoin.WalletAddress),
		sql.Named("ToID", func() int64 {
			if params.TxType == context.EventID_C2P {
				return params.SwapToPoint.PointID
			} else if params.TxType == context.EventID_P2C || params.TxType == context.EventID_C2C {
				return params.SwapToCoin.CoinID
			}
			return 0
		}()),
		sql.Named("ToAdjQuantity", func() string {
			if params.TxType == context.EventID_C2P {
				return strconv.FormatInt(params.SwapToPoint.AdjustPointQuantity, 10)
			} else if params.TxType == context.EventID_P2C || params.TxType == context.EventID_C2C {
				return strconv.FormatFloat(params.SwapToCoin.AdjustCoinQuantity, 'f', -1, 64)
			}
			return ""
		}()),
		sql.Named("PrePointQuantity", func() int64 {
			if params.TxType == context.EventID_C2P {
				return params.SwapToPoint.PreviousPointQuantity
			} else if params.TxType == context.EventID_P2C {
				return params.SwapFromPoint.PreviousPointQuantity
			}
			return 0
		}()),
		sql.Named("PointQuantity", func() int64 {
			if params.TxType == context.EventID_C2P {
				return params.SwapToPoint.PointQuantity
			} else if params.TxType == context.EventID_P2C {
				return params.SwapFromPoint.PointQuantity
			}
			return 0
		}()),
		sql.Named("TxID", sql.Out{Dest: &txID}),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Strt_Exchanges QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	if rs != 1 {
		log.Errorf("USPAU_Strt_Exchanges returnvalue error : %v", rs)
		return &txID, errors.New("USPAU_Strt_Exchanges returnvalue error " + strconv.Itoa(int(rs)))
	}

	return &txID, nil
}

// 가스비 처리
// txStatus 2:수수료 전송 시작, 3:수수료 전송 성공, 4:수수료 전송 실패
func (o *DB) USPAU_Mod_TransactExchanges_ExchangeFees(txID int64, txStatus int64, txHash, swapFee string, baseCoinID int64, gasFee string) error {
	var rs orginMssql.ReturnStatus
	if txStatus == context.SWAP_status_fee_transfer_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_ExchangeFees,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("ExchangeFees", swapFee),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("GasFees", gasFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_ExchangeFees QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_ExchangeFees,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("ExchangeFee", swapFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_ExchangeFees QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchanges_ExchangeFees returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchanges_ExchangeFees returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 거래 상태 갱신
func (o *DB) USPAU_Mod_TransactExchanges_TxStatus(baseCoinID int64, gasFees string, swapInfo *context.ReqSwapInfo) error {
	var rs orginMssql.ReturnStatus
	if swapInfo.TxStatus == context.SWAP_status_fee_transfer_success ||
		swapInfo.TxStatus == context.SWAP_status_token_transfer_deposit_success ||
		swapInfo.TxStatus == context.SWAP_status_token_transfer_withdrawal_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_TxStatus,
			sql.Named("TxID", swapInfo.TxID),
			sql.Named("TxStatus", swapInfo.TxStatus),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("GasFees", gasFees),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_TxStatus QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_TxStatus,
			sql.Named("TxID", swapInfo.TxID),
			sql.Named("TxStatus", swapInfo.TxStatus),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_TxStatus QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchanges_TxStatus returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchanges_TxStatus returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 토큰 처리 : transactionedDT => time.Now().Format("2006-01-02 15:04:05.000")
func (o *DB) USPAU_Mod_TransactExchanges_Coin(txID, txStatus int64, txHash, transactedDT string, baseCoinID int64, gasFee string) error {
	var rs orginMssql.ReturnStatus
	if txStatus == context.SWAP_status_token_transfer_deposit_success ||
		txStatus == context.SWAP_status_token_transfer_withdrawal_success {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_Coin,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("TransactedDT", transactedDT),
			sql.Named("BaseCoinID", baseCoinID),
			sql.Named("GasFees", gasFee),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_Coin QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	} else {
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_TransactExchanges_Coin,
			sql.Named("TxID", txID),
			sql.Named("TxStatus", txStatus),
			sql.Named("TxHash", txHash),
			sql.Named("TransactedDT", transactedDT),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Mod_TransactExchanges_Coin QueryContext err : %v", err)
			return err
		}

		defer rows.Close()
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_TransactExchanges_Coin returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_TransactExchanges_Coin returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}

// swap 종료
// point -> coin, coin->poin, 성공, 실패에 따라 인자값이 달라진다.
func (o *DB) USPAU_Cmplt_Exchanges(params *context.ReqSwapInfo, completedDT string, isSuccess bool) error {
	var rs orginMssql.ReturnStatus
	if params.TxType == context.EventID_P2C {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_Exchanges,
				sql.Named("TxID", params.TxID),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_Exchanges,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.SwapFromPoint.PointID),
				sql.Named("PrePointQuantity", params.SwapFromPoint.PreviousPointQuantity),
				sql.Named("AdjPointQuantity", params.SwapFromPoint.AdjustPointQuantity),
				sql.Named("PointQuantity", params.SwapFromPoint.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	} else if params.TxType == context.EventID_C2P {
		if isSuccess {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_Exchanges,
				sql.Named("TxID", params.TxID),
				sql.Named("PointID", params.SwapToPoint.PointID),
				sql.Named("PrePointQuantity", params.SwapToPoint.PreviousPointQuantity),
				sql.Named("AdjPointQuantity", params.SwapToPoint.AdjustPointQuantity),
				sql.Named("PointQuantity", params.SwapToPoint.PointQuantity),
				sql.Named("CompletedDT", completedDT),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		} else {
			rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_Exchanges,
				sql.Named("TxID", params.TxID),
				&rs)
			if err != nil {
				log.Errorf("USPAU_Cmplt_Exchanges QueryContext err : %v", err)
				return err
			}
			defer rows.Close()
		}
	} else if params.TxType == context.EventID_C2C {
		// c2c는 무조건 데이터 입력해서 종료되도록 유도한다.
		rows, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Cmplt_Exchanges,
			sql.Named("TxID", params.TxID),
			sql.Named("CompletedDT", completedDT),
			&rs)
		if err != nil {
			log.Errorf("USPAU_Cmplt_Exchanges QueryContext err : %v", err)
			return err
		}
		defer rows.Close()
	}
	if rs != 1 {
		log.Errorf("USPAU_Cmplt_Exchanges returnvalue error : %v", rs)
		return errors.New("USPAU_Cmplt_Exchanges returnvalue error " + strconv.Itoa(int(rs)))
	}
	return nil
}

func (o *DB) USPAU_Scan_ExchangeCoinToCoins() error {
	var returnValue orginMssql.ReturnStatus
	proc := USPAU_Scan_ExchangeCoinToCoins
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), proc,
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%v QueryContext error : %v", proc, err)
		return nil
	}

	o.SwapAbleC2Cs = nil
	o.SwapAbleC2CsMap = make(map[int64]map[int64]*context.SwapC2C)

	for rows.Next() {
		swapAble := &context.SwapC2C{}

		if err := rows.Scan(&swapAble.FromBaseCoinID, &swapAble.FromID, &swapAble.ToBaseCoinID, &swapAble.ToID, &swapAble.IsEnabled, &swapAble.MinimumExchangeQuantity, &swapAble.ExchangeRatio); err != nil {
			log.Errorf("%v Scan error : %v", proc, err)
			return err
		} else {
			o.SwapAbleC2Cs = append(o.SwapAbleC2Cs, swapAble)
			if o.SwapAbleC2CsMap[swapAble.FromID] == nil {
				o.SwapAbleC2CsMap[swapAble.FromID] = make(map[int64]*context.SwapC2C)
			}
			o.SwapAbleC2CsMap[swapAble.FromID][swapAble.ToID] = swapAble
		}
	}

	if returnValue != 1 {
		log.Errorf("%v returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return nil
}

func (o *DB) USPAU_Scan_ExchangePointToCoins() error {
	var returnValue orginMssql.ReturnStatus
	proc := USPAU_Scan_ExchangePointToCoins
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), proc,
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%v QueryContext error : %v", proc, err)
		return nil
	}

	o.SwapAbleP2Cs = nil
	o.SwapAbleP2CsMap = make(map[int64]map[int64]*context.SwapP2C)

	for rows.Next() {
		swapAble := &context.SwapP2C{}

		if err := rows.Scan(&swapAble.FromID, &swapAble.ToBaseCoinID, &swapAble.ToID, &swapAble.IsEnabled, &swapAble.MinimumExchangeQuantity, &swapAble.ExchangeRatio); err != nil {
			log.Errorf("%v Scan error : %v", proc, err)
			return err
		} else {
			o.SwapAbleP2Cs = append(o.SwapAbleP2Cs, swapAble)

			if o.SwapAbleP2CsMap[swapAble.FromID] == nil {
				o.SwapAbleP2CsMap[swapAble.FromID] = make(map[int64]*context.SwapP2C)
			}
			o.SwapAbleP2CsMap[swapAble.FromID][swapAble.ToID] = swapAble
		}
	}

	if returnValue != 1 {
		log.Errorf("%v returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return nil
}

func (o *DB) USPAU_Scan_ExchangeCoinToPoints() error {
	var returnValue orginMssql.ReturnStatus
	proc := USPAU_Scan_ExchangeCoinToPoints
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), proc,
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%v QueryContext error : %v", proc, err)
		return nil
	}

	o.SwapAbleC2Ps = nil
	o.SwapAbleC2PsMap = make(map[int64]map[int64]*context.SwapC2P)

	for rows.Next() {
		swapAble := &context.SwapC2P{}

		if err := rows.Scan(&swapAble.FromBaseCoinID, &swapAble.FromID, &swapAble.ToID, &swapAble.IsEnabled, &swapAble.MinimumExchangeQuantity, &swapAble.ExchangeRatio); err != nil {
			log.Errorf("%v Scan error : %v", proc, err)
			return err
		} else {
			o.SwapAbleC2Ps = append(o.SwapAbleC2Ps, swapAble)

			if o.SwapAbleC2PsMap[swapAble.FromID] == nil {
				o.SwapAbleC2PsMap[swapAble.FromID] = make(map[int64]*context.SwapC2P)
			}
			o.SwapAbleC2PsMap[swapAble.FromID][swapAble.ToID] = swapAble
		}
	}

	if returnValue != 1 {
		log.Errorf("%v returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return nil
}
