package model

import (
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertPointMember(params *context.ReqPointMemberRegister) error {

	//	log.Debug("InsertPointMember idx:", lastInsertId)

	return nil
}

func (o *DB) UpdatePointMember(params *context.PointMemberInfo) error {
	sqlQuery := makeUpdateString(params)
	result, err := o.MssqlAccount.PrepareAndExec(sqlQuery)

	if err != nil {
		log.Error(err)
		return err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("UpdatePointMember Affected Count: ", cnt)
	return nil
}

func (o *DB) SelectPointMember(cpMemberIdx int64) (*context.PointMemberInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * from onbuff_inno.dbo.point_member WHERE cp_member_idx=%v", cpMemberIdx)
	rows, err := o.MssqlAccount.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var pointAmount, privateTokenAmount, privateWalletAddr sql.NullString
	var publicTokenAmount, publicWalletAddress sql.NullString

	member := new(context.PointMemberInfo)

	for rows.Next() {
		if err := rows.Scan(&member.Idx, &member.CpMemberIdx, &pointAmount,
			&privateTokenAmount, &privateWalletAddr, &publicTokenAmount, &publicWalletAddress, &member.CreateAt); err != nil {
			log.Error(err)
			return nil, err
		}

		member.PointAmount = pointAmount.String
		member.PrivateTokenAmount = privateTokenAmount.String
		member.PrivateWalletAddr = privateWalletAddr.String
		member.PublicTokenAmount = publicTokenAmount.String
		member.PublicWalletAddr = publicWalletAddress.String
	}

	return member, nil
}

func makeUpdateString(params *context.PointMemberInfo) string {
	sqlQuery := "UPDATE onbuff_inno.dbo.point_member set"

	bValid := false
	if len(params.PointAmount) != 0 {
		sqlQuery += " point_amount=" + params.PointAmount
		bValid = true
	}
	if len(params.PrivateTokenAmount) != 0 {
		getString(&sqlQuery, &bValid)
		sqlQuery += fmt.Sprintf("private_token_amount=N'%v'", params.PrivateTokenAmount)
	}
	if len(params.PrivateWalletAddr) != 0 {
		getString(&sqlQuery, &bValid)
		sqlQuery += fmt.Sprintf("private_wallet_address=N'%v'", params.PrivateWalletAddr)
	}
	if len(params.PublicTokenAmount) != 0 {
		getString(&sqlQuery, &bValid)
		sqlQuery += fmt.Sprintf("public_token_amount=N'%v'", params.PublicTokenAmount)
	}
	if len(params.PublicWalletAddr) != 0 {
		getString(&sqlQuery, &bValid)
		sqlQuery += fmt.Sprintf("public_wallet_address=N'%v'", params.PublicWalletAddr)
	}
	sqlQuery += fmt.Sprintf(" WHERE cp_member_idx=%v", params.CpMemberIdx)
	return sqlQuery
}

func getString(sqlQuery *string, existValid *bool) {
	if *existValid {
		*sqlQuery += ","
	} else {
		*existValid = true
		*sqlQuery += " "
	}
}
