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
	USPAU_GetList_AccountPoints              = "[dbo].[USPAU_GetList_AccountPoints]"
	USPAU_GetList_AccountCoins               = "[dbo].[USPAU_GetList_AccountCoins]"
	USPAU_GetList_AccountCoins_By_CoinString = "[dbo].[USPAU_GetList_AccountCoins_By_CoinString]"
	USPAU_Get_AccountCoins_By_WalletAddress  = "[dbo].[USPAU_Get_AccountCoins_By_WalletAddress]"
	USPAU_Mod_AccountCoins                   = "[dbo].[USPAU_Mod_AccountCoins]"
)

// 계정 일일 포인트량 조회
func (o *DB) GetListAccountPoints(auid, muid int64) (map[int64]*context.AccountPoint, error) {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_AccountPoints,
		sql.Named("AUID", auid),
		sql.Named("MUID", muid),
		&rs)
	if err != nil {
		log.Errorf("USPAU_GetList_AccountPoints QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	accountPoints := make(map[int64]*context.AccountPoint)
	for rows.Next() {
		accountPoint := context.AccountPoint{}
		if err := rows.Scan(&accountPoint.AppId,
			&accountPoint.PointId,
			&accountPoint.TodayAcqQuantity,
			&accountPoint.TodayCnsmQuantity,
			&accountPoint.TodayAcqExchangeQuantity,
			&accountPoint.TodayCnsmExchangeQuantity,
			&accountPoint.ResetDate); err == nil {
			accountPoints[accountPoint.PointId] = &accountPoint
		} else if err != nil {
			log.Errorf("USPAU_GetList_AccountPoints Scan error :%v", err)
		}
	}
	if rs != 1 {
		log.Errorf("USPAU_GetList_AccountPoints returnvalue error : %v", rs)
		return nil, errors.New("USPAU_GetList_AccountPoints returnvalue error " + strconv.Itoa(int(rs)))
	}

	return accountPoints, nil
}

// 코인 정보 조회
func (o *DB) GetAccountCoins(auid int64) ([]*context.AccountCoin, map[int64]*context.AccountCoin, error) {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_AccountCoins,
		sql.Named("AUID", auid),
		&rs)
	if err != nil {
		log.Errorf("USPAU_GetList_AccountCoins QueryContext err : %v", err)
		return nil, nil, err
	}

	defer rows.Close()

	accountCoins := []*context.AccountCoin{}
	accountCoinsMap := make(map[int64]*context.AccountCoin)
	for rows.Next() {
		accountCoin := &context.AccountCoin{}
		if err := rows.Scan(&accountCoin.CoinID,
			&accountCoin.BaseCoinID,
			&accountCoin.WalletAddress,
			&accountCoin.Quantity,
			&accountCoin.TodayAcqQuantity,
			&accountCoin.TodayCnsmQuantity,
			&accountCoin.TodayAcqExchangeQuantity,
			&accountCoin.TodayCnsmExchangeQuantity,
			&accountCoin.ResetDate); err == nil {
			accountCoins = append(accountCoins, accountCoin)
			accountCoinsMap[accountCoin.CoinID] = accountCoin
		} else if err != nil {
			log.Errorf("USPAU_GetList_AccountCoins Scan error :%v", err)
		}
	}
	if rs != 1 {
		log.Errorf("USPAU_GetList_AccountCoins returnvalue error : %v", rs)
		return nil, nil, errors.New("USPAU_GetList_AccountCoins returnvalue error " + strconv.Itoa(int(rs)))
	}

	return accountCoins, accountCoinsMap, nil
}

// 지갑 정보 조회
func (o *DB) GetPointMemberWallet(params *context.ReqPointMemberWallet, appID int64) (*context.ResPointMemberWallet, error) {
	coinIds := ""
	for _, coinId := range o.AppCoins[appID] {
		coinIds += "/" + strconv.FormatInt(coinId.CoinId, 10)
	}

	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_AccountCoins_By_CoinString,
		sql.Named("AUID", params.AUID),
		sql.Named("CoinString", coinIds),
		sql.Named("RowSeparator", "/"),
		&rs)
	if err != nil {
		log.Errorf("USPAU_GetList_AccountCoins_By_CoinString QueryContext err : %v", err)
		return nil, err
	}

	defer rows.Close()

	walletInfos := &context.ResPointMemberWallet{
		AUID: params.AUID,
	}
	WalletInfo := context.WalletInfo{}
	for rows.Next() {
		if err := rows.Scan(&WalletInfo.CoinID, &WalletInfo.BaseCoinID, &WalletInfo.WalletAddress, &WalletInfo.CoinQuantity); err == nil {
			WalletInfo.CoinSymbol = o.Coins[WalletInfo.CoinID].CoinSymbol
			walletInfos.WalletInfo = append(walletInfos.WalletInfo, WalletInfo)
		}
	}

	if rs != 1 {
		log.Errorf("USPAU_GetList_AccountCoins_By_CoinString returnvalue error : %v", rs)
		return nil, errors.New("USPAU_GetList_AccountCoins_By_CoinString returnvalue error " + strconv.Itoa(int(rs)))
	}

	return walletInfos, nil
}

// 지갑 주소로 AUID, coinID 검색
func (o *DB) GetAccountCoinsByWalletAddress(walletAddress, coinSymbol string) (*context.AccountCoinByWalletAddress, error) {
	var rs orginMssql.ReturnStatus
	var auID, coinID int64
	var quantity float64
	_, err := o.MssqlAccountRead.GetDB().QueryContext(originCtx.Background(), USPAU_Get_AccountCoins_By_WalletAddress,
		sql.Named("WalletAddress", walletAddress),
		sql.Named("CoinSymbol", coinSymbol),
		sql.Named("AUID", sql.Out{Dest: &auID}),
		sql.Named("CoinID", sql.Out{Dest: &coinID}),
		sql.Named("Quantity", sql.Out{Dest: &quantity}),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Get_AccountCoins_By_WalletAddress QueryContext err : %v", err)
		return nil, err
	}

	meCoin := &context.AccountCoinByWalletAddress{
		AUID:     auID,
		CoinID:   coinID,
		Quantity: quantity,
	}

	if rs != 1 {
		log.Errorf("USPAU_Get_AccountCoins_By_WalletAddress returnvalue error : %v", rs)
		return nil, errors.New("USPAU_Get_AccountCoins_By_WalletAddress returnvalue error " + strconv.Itoa(int(rs)))
	}

	return meCoin, nil
}

// 내 코인 정보 수정
func (o *DB) UpdateAccountCoins(auid, coinid, baseCoinID int64, walletAddress string, previousCoinQuantity, adjustCoinQuantity, coinQuantity float64,
	logID context.LogID_type, eventID context.EventID_type) error {

	var rs orginMssql.ReturnStatus
	_, err := o.MssqlAccountAll.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_AccountCoins,
		sql.Named("AUID", auid),
		sql.Named("CoinID", coinid),
		sql.Named("BaseCoinID", baseCoinID),
		sql.Named("WalletAddress", walletAddress),
		sql.Named("PreQuantity", previousCoinQuantity),
		sql.Named("AdjQuantity", adjustCoinQuantity),
		sql.Named("Quantity", coinQuantity),
		sql.Named("LogID", logID),
		sql.Named("EventID", eventID),
		&rs)
	if err != nil {
		log.Errorf("USPAU_Mod_AccountCoins QueryContext err : %v", err)
		return err
	}

	if rs != 1 {
		log.Errorf("USPAU_Mod_AccountCoins returnvalue error : %v", rs)
		return errors.New("USPAU_Mod_AccountCoins returnvalue error " + strconv.Itoa(int(rs)))
	}

	return nil
}
