package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/util"
)

func (o *DB) InsertPointAppHistory(params *context.PointMemberAppUpdate) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.point_history(cp_member_idx, "+
		"type, latest_point_amount, change_point_amount, create_at) output inserted.idx "+
		"VALUES(%v,N'%v',N'%v',N'%v',%v)",
		params.CpMemberIdx, params.Type, params.LatestPointAmount, params.ChangePointAmount, params.CreateAt)

	var lastInsertId int64
	err := o.MssqlAccount.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertPointAppHistory idx:", lastInsertId)

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
