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
	USPPO_Get_MemberPoints     = "[dbo].[USPPO_Get_MemberPoints]"
	USPPO_Mod_MemberPoints     = "[dbo].[USPPO_Mod_MemberPoints]"
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

// 맴버의 포인트 리스트 정보 조회
func (o *DB) GetPointAppList(MUID, DatabaseID int64) ([]*context.Point, error) {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_GetList_MemberPoints,
		sql.Named("MUID", MUID),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	defer rows.Close()

	points := []*context.Point{}

	point := new(context.Point)
	for rows.Next() {
		point.PointID = 0
		point.Quantity = 0
		if err := rows.Scan(&point.PointID, &point.Quantity); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return points, nil
}

// 맴버의 포인트 정보 조회
func (o *DB) GetPointApp(MUID, PointID, DatabaseID int64) (*context.Point, error) {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Get_MemberPoints,
		sql.Named("MUID", MUID),
		sql.Named("PointID", PointID),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	defer rows.Close()

	rowCnt := 0
	point := new(context.Point)
	for rows.Next() {
		point.PointID = PointID
		point.Quantity = 0
		if err := rows.Scan(&point.Quantity); err != nil {
			return nil, err
		} else {
			rowCnt++
		}
	}

	if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	if rowCnt == 0 {
		return nil, nil
	}

	return point, nil
}

// 맴버의 포인트 정보 조회 by point id
func (o *DB) GetPointAppByPointID(MUID, pointId, DatabaseID int64) (*context.Point, error) {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Get_MemberPoints,
		sql.Named("MUID", MUID),
		sql.Named("PointID", pointId),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	defer rows.Close()

	point := new(context.Point)
	for rows.Next() {
		point.PointID = 0
		point.Quantity = 0
		if err := rows.Scan(&point.PointID, &point.Quantity); err != nil {
			return nil, err
		}
	}

	if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return point, nil
}

// 포인트 업데이트
func (o *DB) UpdateAppPoint(muid, pointId, preQuantity, adjQuantity, quantity, dbId int64) (int64, string, error) {
	mssql, ok := o.MssqlPoints[dbId]
	if !ok {
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var dailyQuantity int64
	var resetDate string
	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Mod_MemberPoints,
		sql.Named("MUID", muid),
		sql.Named("PointID", pointId),
		sql.Named("PreQuantity", preQuantity),
		sql.Named("AdjQuantity", adjQuantity),
		sql.Named("Quantity", quantity),

		sql.Named("DailyQuantity", sql.Out{Dest: &dailyQuantity}),
		sql.Named("ResetDate", sql.Out{Dest: &resetDate}),
		&rs); err != nil {
		log.Error("QueryContext err : ", err)
		return 0, "", err
	}

	if rs == resultcode.Result_Error_Invalid_data {
		log.Error("returnStatus Result_Error_Invalid_data : ", rs)
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return dailyQuantity, resetDate, nil
}
