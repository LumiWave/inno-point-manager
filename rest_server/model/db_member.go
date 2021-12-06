package model

// func (o *DB) UpdatePointMember(params *context.PointMemberInfo) error {
// 	sqlQuery := makeUpdateString(params)
// 	result, err := o.MssqlAccount.PrepareAndExec(sqlQuery)

// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}

// 	cnt, err := result.RowsAffected()
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	log.Debug("UpdatePointMember Affected Count: ", cnt)
// 	return nil
// }

// func makeUpdateString(params *context.PointMemberInfo) string {
// 	sqlQuery := "UPDATE onbuff_inno.dbo.point_member set"

// 	bValid := false
// 	if len(params.PointAmount) != 0 {
// 		sqlQuery += " point_amount=" + params.PointAmount
// 		bValid = true
// 	}
// 	if len(params.PrivateTokenAmount) != 0 {
// 		getString(&sqlQuery, &bValid)
// 		sqlQuery += fmt.Sprintf("private_token_amount=N'%v'", params.PrivateTokenAmount)
// 	}
// 	if len(params.PrivateWalletAddr) != 0 {
// 		getString(&sqlQuery, &bValid)
// 		sqlQuery += fmt.Sprintf("private_wallet_address=N'%v'", params.PrivateWalletAddr)
// 	}
// 	if len(params.PublicTokenAmount) != 0 {
// 		getString(&sqlQuery, &bValid)
// 		sqlQuery += fmt.Sprintf("public_token_amount=N'%v'", params.PublicTokenAmount)
// 	}
// 	if len(params.PublicWalletAddr) != 0 {
// 		getString(&sqlQuery, &bValid)
// 		sqlQuery += fmt.Sprintf("public_wallet_address=N'%v'", params.PublicWalletAddr)
// 	}
// 	sqlQuery += fmt.Sprintf(" WHERE cp_member_idx=%v", params.CpMemberIdx)
// 	return sqlQuery
// }

// func getString(sqlQuery *string, existValid *bool) {
// 	if *existValid {
// 		*sqlQuery += ","
// 	} else {
// 		*existValid = true
// 		*sqlQuery += " "
// 	}
// }
