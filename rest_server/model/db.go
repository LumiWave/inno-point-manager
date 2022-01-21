package model

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

type PointInfo struct {
	PointId              int64   `json:"point_id,omitempty"`
	PointName            string  `json:"point_name,omitempty"`
	IconUrl              string  `json:"icon_url,omitempty"`
	MinExchangeQuantity  int64   `json:"minimum_exchange_quantity"`
	ExchangeRatio        float64 `json:"exchange_ratio"`
	DaliyLimitedQuantity int64   `json:"daliy_limited_quantity,omitempty"`
}

type AppPointInfo struct {
	AppId     int64                `json:"app_id,omitempty"`
	AppName   string               `json:"app_name,omitempty"`
	IconUrl   string               `json:"icon_url"`
	Points    []*PointInfo         `json:"points"`
	PointsMap map[int64]*PointInfo `json:"-"` // key pointId
}

type Coin struct {
	CoinId          int64   `json:"coin_id,omitempty"`
	CoinSymbol      string  `json:"coin_symbol,omitempty"`
	ContractAddress string  `json:"contract_address,omitempty"`
	IconUrl         string  `json:"icon_url,omitempty"`
	ExchangeFees    float64 `json:"exchange_fees"`
}

type AppCoin struct {
	AppID int64 `json:"app_id"`
	Coin
}

type DB struct {
	Mysql        *basedb.Mysql
	MssqlAccount *basedb.Mssql
	Cache        *basedb.Cache

	MssqlPoints map[int64]*basedb.Mssql

	PointDoc map[string]*MemberPointInfo

	//PointList     map[int64]PointInfo // 전체 포인트 종류
	ScanPointsMap map[int64]PointInfo // 전체 포인트 종류 : key PointId

	AppPointsMap map[int64]*AppPointInfo // 전체 app과 포인트 : key AppId

	AppCoins map[int64][]*AppCoin // 전체 app에 속한 CoinID 정보
	Coins    map[int64]*Coin      // 전체 coin 정보
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
	gDB.ScanPointsMap = make(map[int64]PointInfo)
	gDB.AppPointsMap = make(map[int64]*AppPointInfo)
	gDB.AppCoins = make(map[int64][]*AppCoin)
	gDB.Coins = make(map[int64]*Coin)
	gDB.MssqlPoints = pointdbs

	// sequence is important
	gDB.GetPointList()
	gDB.GetAppCoins()
	gDB.GetCoins()
	gDB.GetApps()
	gDB.GetAppPoints()
}

func GetDB() *DB {
	return gDB
}

func MakeDbError(resp *base.BaseResponse, errCode int, err error) {
	resp.Return = errCode
	resp.Message = resultcode.ResultCodeText[errCode] + " : " + err.Error()
}
