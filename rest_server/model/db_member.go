package model

import (
	originCtx "context"
	"database/sql"
	"errors"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPPO_Rgstr_Members        = "[dbo].[USPPO_Rgstr_Members]"
	USPPO_GetList_MemberPoints = "[dbo].[USPPO_GetList_MemberPoints]"
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
		sql.Named("CUID", params.CUID),
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

// 맴버의 포인트 정보 조회
func (o *DB) GetPointMember(CUID string, AppID, DatabaseID int64) (*[]context.Point, error) {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_GetList_MemberPoints,
		sql.Named("CUID", CUID),
		sql.Named("AppID", AppID),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	points := new([]context.Point)

	point := context.Point{}
	for rows.Next() {
		point.PointID = 0
		point.Quantity = 0
		if err := rows.Scan(&point.PointID, &point.Quantity); err != nil {
			return nil, err
		}
		*points = append(*points, point)
	}

	if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return points, nil
}

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
