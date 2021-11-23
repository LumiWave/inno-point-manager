package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/util"
)

func (o *DB) InsertPointTokenSwapHistory(params *context.PointMemberTokenSwapHistory) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.point_swap_history(cp_member_idx, "+
		"latest_private_token_amount, swap_private_token_amount, latest_public_token_amount, swap_public_token_amount, txn_hash, swap_state, create_at) output inserted.idx "+
		"VALUES(%v,N'%v',N'%v',N'%v',N'%v',N'%v',N'%v',%v)",
		params.CpMemberIdx,
		params.LatestPrivateTokenAmount, params.SwapPrivateTokenAmount, params.LatestPublicTokenAmount, params.SwapPublicTokenAmount,
		params.TxnHash, params.SwapState, params.CreateAt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertPointTokenSwapHistory idx:", lastInsertId)

	return nil
}

func (o *DB) SelectPointTokenSwapHistory(params *context.PointMemberTokenSwapHistory) (*[]context.PointMemberTokenSwapHistory, int64, error) {
	var totalCount int64
	sqlQuery := fmt.Sprintf("SELECT COUNT(*) FROM onbuff_inno.dbo.point_swap_history WHERE cp_member_idx=%v", params.CpMemberIdx)
	err := o.Mssql.QueryRow(sqlQuery, &totalCount)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	pageSize := util.ParseInt(params.PageSize)
	pageOffset := util.ParseInt(params.PageOffset)

	sqlQuery = fmt.Sprintf("SELECT * from onbuff_inno.dbo.point_swap_history WHERE cp_member_idx=%v ORDER BY idx DESC OFFSET %v ROW FETCH NEXT %v ROW ONLY ",
		params.CpMemberIdx, pageSize*pageOffset, pageSize)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	historys := make([]context.PointMemberTokenSwapHistory, 0)
	for rows.Next() {
		history := &context.PointMemberTokenSwapHistory{}
		if err := rows.Scan(&history.Idx, &history.CpMemberIdx,
			&history.LatestPrivateTokenAmount, &history.SwapPrivateTokenAmount,
			&history.LatestPublicTokenAmount, &history.SwapPublicTokenAmount,
			&history.TxnHash, &history.SwapState, &history.CreateAt); err != nil {
			log.Error("SelectPointTokenSwapHistory::Scan error : ", err)
		} else {
			historys = append(historys, *history)
		}

	}

	return &historys, totalCount, nil
}
