package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPPO_Rgstr_Members        = "[dbo].[USPPO_Rgstr_Members]"
	USPAU_GetList_AccountCoins = "[dbo].[USPAU_GetList_AccountCoins]"
)

// 포인트 맴버 등록
func (o *DB) InsertPointMember(params *context.ReqPointMemberRegister) error {
	mssql, ok := o.MssqlPoints[params.DatabaseID]
	if !ok {
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}
	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Rgstr_Members,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", params.MUID),
		sql.Named("AppID", params.AppID),
		&rs); err != nil {
		log.Error("QueryContext err : ", err)
		return err
	}

	if rs == resultcode.Result_Error_duplicate_auid {
		log.Error("returnStatus Result_Error_duplicate_auid : ", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return nil
}

// 지갑 정보 조회
func (o *DB) GetPointMemberWallet(params *context.ReqPointMemberWallet, appID int64) (*context.ResPointMemberWallet, error) {

	coinIds := ""
	for _, coinId := range o.AppCoins[appID] {
		coinIds += "/" + strconv.FormatInt(coinId.CoinID, 10)
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
			WalletInfo.CoinSymbol = o.Coins[WalletInfo.CoinID].CoinName
			walletInfos.WalletInfo = append(walletInfos.WalletInfo, WalletInfo)
		}
	}
	return walletInfos, nil
}
