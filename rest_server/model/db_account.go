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
	USPAU_Scan_DatabaseServers = "[dbo].[USPAU_Scan_DatabaseServers]"
	USPAU_GetList_MemberPoints = "[dbo].[USPAU_GetList_MemberPoints]"
)

// point database 리스트 요청
func (o *DB) GetPointDatabases() (map[int64]*PointDB, error) {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccount.GetDB().QueryContext(originCtx.Background(), USPAU_Scan_DatabaseServers, &rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	pointdbs := make(map[int64]*PointDB)

	pointdb := new(PointDB)
	for rows.Next() {
		rows.Scan(&pointdb.DatabaseID, &pointdb.DatabaseName, &pointdb.ServerName)
		pointdbs[pointdb.DatabaseID] = pointdb
	}

	return pointdbs, nil
}

// 맴버의 포인트 정보 조회
func (o *DB) GetPointMember(CUID string, AppID, DatabaseID int64) (*[]context.Point, error) {
	mssql, ok := o.MssqlPoints[DatabaseID]
	if !ok {
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Invalid_DBID])
	}

	var rs orginMssql.ReturnStatus
	rows, err := mssql.GetDB().QueryContext(originCtx.Background(), USPAU_GetList_MemberPoints,
		sql.Named("CUID", CUID),
		sql.Named("AppID", AppID),
		&rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return nil, err
	}

	if rs == resultcode.Result_Error_Invalid_data {
		log.Error("returnStatus Result_Error_Invalid_data : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_Error_Invalid_data])
	} else if rs != 1 {
		log.Error("returnStatus Result_DBError_Unknown : ", rs)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_DBError_Unknown])
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

	return points, nil
}
