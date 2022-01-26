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
	USPPO_Add_MemberPoints     = "[dbo].[USPPO_Add_MemberPoints]"
)

// 포인트 맴버 등록
func (o *DB) InsertPointMember(params *context.ReqPointMemberRegister) error {
	mssql, ok := o.MssqlPointsAll[params.DatabaseID]
	if !ok {
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}
	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Rgstr_Members,
		sql.Named("AUID", params.AUID),
		sql.Named("MUID", params.MUID),
		sql.Named("AppID", params.AppID),
		&rs); err != nil {
		log.Errorf("USPPO_Rgstr_Members QueryContext error : %v", err)
		return err
	}

	if rs == resultcode.Result_Error_duplicate_auid {
		log.Errorf("USPPO_Rgstr_Members returnStatus : %v", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Errorf("USPPO_Rgstr_Members returnStatus : %v", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return nil
}

// 맴버의 포인트 리스트 정보 조회
func (o *DB) GetPointAppList(MUID, DatabaseID int64) ([]*context.Point, error) {
	mssql, ok := o.MssqlPointsRead[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_GetList_MemberPoints,
		sql.Named("MUID", MUID),
		&rs)
	if err != nil {
		log.Errorf("USPPO_GetList_MemberPoints QueryContext error : %v", err)
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
		point.PreQuantity = point.Quantity // load 시 동일하게 초기화
		points = append(points, point)
	}

	if rs != 1 {
		log.Errorf("USPPO_GetList_MemberPoints returnStatus : %v", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return points, nil
}

// 맴버의 포인트 정보 조회
func (o *DB) GetPointApp(MUID, PointID, DatabaseID int64) (*context.Point, error) {
	mssql, ok := o.MssqlPointsRead[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Get_MemberPoints,
		sql.Named("MUID", MUID),
		sql.Named("PointID", PointID),
		&rs)
	if err != nil {
		log.Errorf("USPPO_Get_MemberPoints QueryContext error : %v", err)
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
		log.Errorf("USPPO_Get_MemberPoints returnStatus : %v", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	if rowCnt == 0 {
		return nil, nil
	}

	return point, nil
}

// 포인트 최초 초기화 등록
func (o *DB) InsertMemberPoints(dbID, muID, pointID, quantity int64) error {
	mssql, ok := o.MssqlPointsAll[dbID]
	if !ok {
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Add_MemberPoints,
		sql.Named("MUID", muID),
		sql.Named("PointID", pointID),
		sql.Named("Quantity", quantity),
		&rs); err != nil {
		log.Errorf("USPPO_Add_MemberPoints QueryContext error : %v", err)
		return err
	}

	if rs == resultcode.Result_Error_Invalid_data {
		log.Errorf("USPPO_Add_MemberPoints returnStatus : %v", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Errorf("USPPO_Add_MemberPoints returnStatus : %v", rs)
		return errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return nil
}

// 포인트 업데이트
func (o *DB) UpdateAppPoint(dbID, muID, pointID, preQuantity, adjQuantity, quantity int64, logID context.LogID_type, eventID context.EventID_type) (int64, string, error) {
	mssql, ok := o.MssqlPointsAll[dbID]
	if !ok {
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var todayLimitedQuantity int64
	var resetDate string
	var rs orginMssql.ReturnStatus
	if _, err := mssql.GetDB().QueryContext(originCtx.Background(), USPPO_Mod_MemberPoints,
		sql.Named("MUID", muID),
		sql.Named("PointID", pointID),
		sql.Named("PreQuantity", preQuantity),
		sql.Named("AdjQuantity", adjQuantity),
		sql.Named("Quantity", quantity),
		sql.Named("LogID", logID),
		sql.Named("EventID", eventID),

		sql.Named("TodayLimitedQuantity", sql.Out{Dest: &todayLimitedQuantity}),
		sql.Named("ResetDate", sql.Out{Dest: &resetDate}),
		&rs); err != nil {
		log.Errorf("USPPO_Mod_MemberPoints QueryContext error : %v", err)
		return 0, "", err
	}

	if rs == resultcode.Result_Error_Invalid_data {
		log.Errorf("USPPO_Mod_MemberPoints returnStatus Result_Error_Invalid_data : %v", rs)
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_Error_duplicate_auid])
	} else if rs != 1 {
		log.Errorf("USPPO_Mod_MemberPoints returnStatus Result_DBError_Unknown : %v", rs)
		return 0, "", errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
	}

	return todayLimitedQuantity, resetDate, nil
}
