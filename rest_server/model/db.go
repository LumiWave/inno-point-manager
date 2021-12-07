package model

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

type PointDB struct {
	DatabaseID   int64
	DatabaseName string
	ServerName   string
}

type Point struct {
	PointIds []int64
}

type DB struct {
	Mysql        *basedb.Mysql
	MssqlAccount *basedb.Mssql
	Cache        *basedb.Cache

	MssqlPoints map[int64]*basedb.Mssql

	PointDoc map[string]*MemberPointInfo

	PointList map[int64]Point // 전체 포인트 종류
}

var gDB *DB

func SetDB(db *basedb.Mssql, cache *basedb.Cache, pointdbs map[int64]*basedb.Mssql) {
	gDB = &DB{
		MssqlAccount: db,
		Cache:        cache,
		MssqlPoints:  pointdbs,
	}
}

func SetDBPoint(pointdbs map[int64]*basedb.Mssql) {
	gDB.PointDoc = make(map[string]*MemberPointInfo)
	gDB.PointList = make(map[int64]Point)
	gDB.MssqlPoints = pointdbs

	gDB.GetPointList()
}

func GetDB() *DB {
	return gDB
}

func MakeDbError(resp *base.BaseResponse, errCode int, err error) {
	resp.Return = errCode
	resp.Message = resultcode.ResultCodeText[errCode] + " : " + err.Error()
}
