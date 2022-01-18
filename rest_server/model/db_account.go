package model

import (
	originCtx "context"
	"database/sql"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_GetList_AccountPoints = "[dbo].[USPAU_GetList_AccountPoints]"
	USPAU_GetList_AccountCoins  = "[dbo].[USPAU_GetList_AccountCoins_By_CoinString]"
	USPAU_Mod_ApplicationPoints = "[dbo].[USPAU_Mod_ApplicationPoints]"
)

// 계정 일일 포인트량 조회
func (o *DB) GetListAccountPoints(auid, muid int64) (map[int64]*context.AccountPoint, error) {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccount.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_AccountPoints,
		sql.Named("AUID", auid),
		sql.Named("MUID", muid),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	defer rows.Close()

	accountPoints := make(map[int64]*context.AccountPoint)
	accountPoint := context.AccountPoint{}
	for rows.Next() {
		if err := rows.Scan(&accountPoint.AppId, &accountPoint.PointId, &accountPoint.DailyQuantity, &accountPoint.ResetDate); err == nil {
			accountPoints[accountPoint.PointId] = &accountPoint
		}
	}
	return accountPoints, nil
}

// 지갑 정보 조회
func (o *DB) GetPointMemberWallet(params *context.ReqPointMemberWallet, appID int64) (*context.ResPointMemberWallet, error) {
	coinIds := ""
	for _, coinId := range o.AppCoins[appID] {
		coinIds += "/" + strconv.FormatInt(coinId.CoinId, 10)
	}

	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccount.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_AccountCoins,
		sql.Named("AUID", params.AUID),
		sql.Named("CoinString", coinIds),
		sql.Named("RowSeparator", "/"),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	defer rows.Close()

	walletInfos := &context.ResPointMemberWallet{
		AUID: params.AUID,
	}
	WalletInfo := context.WalletInfo{}
	for rows.Next() {
		if err := rows.Scan(&WalletInfo.CoinID, &WalletInfo.WalletAddress, &WalletInfo.CoinQuantity); err == nil {
			WalletInfo.CoinSymbol = o.Coins[WalletInfo.CoinID].CoinSymbol
			walletInfos.WalletInfo = append(walletInfos.WalletInfo, WalletInfo)
		}
	}
	return walletInfos, nil
}

func (o *DB) UpdateApplicationPoints(appId, pointId, adjustQuantity, adjustExchangeQuantity int64) (int64, int64, string, error) {
	var rs orginMssql.ReturnStatus

	var dailyQuantity, dailyExchangeQuantity int64
	var resetDate string
	_, err := o.MssqlAccount.GetDB().QueryContext(originCtx.Background(), USPAU_Mod_ApplicationPoints,
		sql.Named("AppID", appId),
		sql.Named("PointID", pointId),
		sql.Named("AdjQuantity", adjustQuantity),
		sql.Named("AdjExchangeQuantity", adjustExchangeQuantity),

		sql.Named("DailyQuantity", sql.Out{Dest: &dailyQuantity}),
		sql.Named("DailyExchangeQuantity", sql.Out{Dest: &dailyExchangeQuantity}),
		sql.Named("ResetDate", sql.Out{Dest: &resetDate}),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return 0, 0, "", err
	}

	return dailyQuantity, dailyExchangeQuantity, resetDate, nil
}
