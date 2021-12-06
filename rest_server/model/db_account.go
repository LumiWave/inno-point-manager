package model

import (
	originCtx "context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_DatabaseServers = "[dbo].[USPAU_Scan_DatabaseServers]"
)

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
