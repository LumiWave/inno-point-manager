package model

import (
	originCtx "context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/util"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPPO_Mod_MemberPoints = "[dbo].[USPPO_Mod_MemberPoints]"
)

// 포인트 업데이트
func (o *DB) UpdateAppPoint(CUID string, AppID, PointID, Quantity, DatabaseID int64) error {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}
	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Mod_MemberPoints,
		sql.Named("CUID", CUID),
		sql.Named("AppID", AppID),
		sql.Named("PointID", PointID),
		sql.Named("Quantity", Quantity),
		&rs); err != nil {
		log.Error("QueryContext err : ", err)
		return err
	}

	if rs == resultcode.Result_Error_Invalid_data {
		log.Error("returnStatus Result_Error_Invalid_data : ", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return nil
}

func (o *DB) InsertPointAppHistory(params *context.ReqPointMemberAppUpdate) error {
	// sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.point_history(cp_member_idx, "+
	// 	"type, latest_point_amount, change_point_amount, create_at) output inserted.idx "+
	// 	"VALUES(%v,N'%v',N'%v',N'%v',%v)",
	// 	params.CpMemberIdx, params.Type, params.LatestPointAmount, params.ChangePointAmount, params.CreateAt)

	// var lastInsertId int64
	// err := o.MssqlAccount.QueryRow(sqlQuery, &lastInsertId)

	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	// log.Debug("InsertPointAppHistory idx:", lastInsertId)

	return nil
}

func (o *DB) SelectPointAppHistory(params *context.PointMemberHistoryReq) (*[]context.PointMemberHistory, int64, error) {

	var totalCount int64
	sqlQuery := fmt.Sprintf("SELECT COUNT(*) FROM onbuff_inno.dbo.point_history WHERE cp_member_idx=%v", params.CpMemberIdx)
	err := o.MssqlAccount.QueryRow(sqlQuery, &totalCount)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	pageSize := util.ParseInt(params.PageSize)
	pageOffset := util.ParseInt(params.PageOffset)

	sqlQuery = fmt.Sprintf("SELECT * from onbuff_inno.dbo.point_history WHERE cp_member_idx=%v ORDER BY idx DESC OFFSET %v ROW FETCH NEXT %v ROW ONLY ",
		params.CpMemberIdx, pageSize*pageOffset, pageSize)
	rows, err := o.MssqlAccount.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	historys := make([]context.PointMemberHistory, 0)
	for rows.Next() {
		history := &context.PointMemberHistory{}
		if err := rows.Scan(&history.Idx, &history.CpMemberIdx, &history.Type, &history.LatestPointAmount, &history.ChangePointAmount, &history.CreateAt); err != nil {
			log.Error("SelectPointAppHistory::Scan error : ", err)
		} else {
			historys = append(historys, *history)
		}

	}

	return &historys, totalCount, nil
}

func (o *DB) InsertPointAppExchangeHistory(params *context.PointMemberExchangeHistory) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.point_exchange_history(cp_member_idx, "+
		"latest_point_amount, exchange_point_amount, latest_private_token_amount, exchange_private_token_amount, txn_hash, exchange_state, create_at) output inserted.idx "+
		"VALUES(%v,N'%v',N'%v',N'%v',N'%v',N'%v',N'%v',%v)",
		params.CpMemberIdx, params.LatestPointAmount, params.ExchangePointAmount, params.LatestPrivateTokenAmount, params.ExchangePrivateTokenAmount,
		params.TxnHash, params.ExchangeState, params.CreateAt)

	var lastInsertId int64
	err := o.MssqlAccount.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertPointAppExchangeHistory idx:", lastInsertId)

	return nil
}

func (o *DB) SelectPointAppExchangeHistory(params *context.PointMemberExchangeHistory) (*[]context.PointMemberExchangeHistory, int64, error) {
	var totalCount int64
	sqlQuery := fmt.Sprintf("SELECT COUNT(*) FROM onbuff_inno.dbo.point_exchange_history WHERE cp_member_idx=%v", params.CpMemberIdx)
	err := o.MssqlAccount.QueryRow(sqlQuery, &totalCount)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	pageSize := util.ParseInt(params.PageSize)
	pageOffset := util.ParseInt(params.PageOffset)

	sqlQuery = fmt.Sprintf("SELECT * from onbuff_inno.dbo.point_exchange_history WHERE cp_member_idx=%v ORDER BY idx DESC OFFSET %v ROW FETCH NEXT %v ROW ONLY ",
		params.CpMemberIdx, pageSize*pageOffset, pageSize)
	rows, err := o.MssqlAccount.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	historys := make([]context.PointMemberExchangeHistory, 0)
	for rows.Next() {
		history := &context.PointMemberExchangeHistory{}
		if err := rows.Scan(&history.Idx, &history.CpMemberIdx, &history.LatestPointAmount, &history.ExchangePointAmount,
			&history.LatestPrivateTokenAmount, &history.ExchangePrivateTokenAmount, &history.TxnHash, &history.ExchangeState, &history.CreateAt); err != nil {
			log.Error("SelectPointAppExchangeHistory::Scan error : ", err)
		} else {
			historys = append(historys, *history)
		}

	}

	return &historys, totalCount, nil
}
