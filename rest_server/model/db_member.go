package model

import (
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
)

func (o *DB) InsertPointMember(params *context.PointMemberInfo) error {
	// sqlQuery := fmt.Sprintf("EXEC @return_value=onbuff_inno.dbo.proc_point_member_insert @cp_member_idx=?,@point_amount=?," +
	// 	"@private_token_amount=?,@private_wallet_address=?,@public_token_amount=?,@public_wallet_address=?," +
	// 	"@create_at=? SELECT Return Value=@return_value")

	// result, err := o.Mssql.PrepareAndExec(sqlQuery, params.CpMemberIdx, params.PointAmount,
	// 	params.PrivateTokenAmount, params.PrivateWalletAddr, params.PublicTokenAmount, params.PublicWalletAddr, params.CreateAt)

	// sqlQuery := fmt.Sprintf("DECLARE @return_value int; EXEC @return_value=onbuff_inno.dbo.proc_point_member_insert @cp_member_idx=%v,@point_amount=N'%v',"+
	// 	"@private_token_amount=N'%v',@private_wallet_address=N'%v',@public_token_amount=N'%v',@public_wallet_address=N'%v',"+
	// 	"@create_at=%v;  SELECT 'Return Value'=@return_value", params.CpMemberIdx, params.PointAmount,
	// 	params.PrivateTokenAmount, params.PrivateWalletAddr, params.PublicTokenAmount, params.PublicWalletAddr, params.CreateAt)

	// rows, err := o.Mssql.Query(sqlQuery)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// ret := 0
	// for rows.Next() {

	// 	err := rows.Scan(&ret)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	} else {
	// 		fmt.Println("ret : ", ret)
	// 	}
	// }

	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.point_member(cp_member_idx, point_amount, "+
		"private_token_amount, private_wallet_address, public_token_amount, public_wallet_address, create_at) output inserted.idx "+
		"VALUES(%v,N'%v',N'%v',N'%v',N'%v',N'%v',%v)",
		params.CpMemberIdx, params.PointAmount, params.PrivateTokenAmount, params.PrivateWalletAddr, params.PublicTokenAmount, params.PublicWalletAddr, params.CreateAt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertPointMember idx:", lastInsertId)

	return nil
}

func (o *DB) UpdatePointMember(params *context.PointMemberInfo) error {
	sqlQuery := makeUpdateString(params)
	result, err := o.Mssql.PrepareAndExec(sqlQuery)

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
	rows, err := o.Mssql.Query(sqlQuery)
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
