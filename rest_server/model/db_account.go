package model

import (
	originCtx "context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_DatabaseServers = "[dbo].[USPAU_Scan_DatabaseServers]"
	USPAU_Scan_Points          = "[dbo].[USPAU_Scan_Points]"
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

// point 전체 list
func (o *DB) GetPointList() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccount.GetDB().QueryContext(originCtx.Background(), USPAU_Scan_Points, &rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return err
	}

	var pointId, appId int64
	for rows.Next() {
		if err := rows.Scan(&pointId, &appId); err == nil {
			points := o.PointList[appId]
			points.PointIds = append(points.PointIds, pointId)
			o.PointList[appId] = points
		}
	}

	return nil
}
