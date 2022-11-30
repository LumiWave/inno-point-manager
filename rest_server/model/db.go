package model

import (
	"strconv"
	"sync"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type PointInfo struct {
	PointId                         int64   `json:"point_id,omitempty"`
	PointName                       string  `json:"point_name,omitempty"`
	IconUrl                         string  `json:"icon_url,omitempty"`
	MinExchangeQuantity             int64   `json:"minimum_exchange_quantity"`
	ExchangeRatio                   float64 `json:"exchange_ratio"`
	DaliyLimitedAcqQuantity         int64   `json:"daliy_limited_acq_quantity,omitempty"`
	DailyLimitedAcqExchangeQuantity int64   `json:"daily_limited_acq_exchange_quantity,omitempty"`
}

type AppPointInfo struct {
	AppId            int64                `json:"app_id,omitempty"`
	AppName          string               `json:"app_name,omitempty"`
	IconUrl          string               `json:"icon_url"`
	GooglePlayPath   string               `json:"google_play_path"`
	AppleStorePath   string               `json:"apple_store_path"`
	BrandingPagePath string               `json:"branding_page_path"`
	Points           []*PointInfo         `json:"points"`
	PointsMap        map[int64]*PointInfo `json:"-"` // key pointId
}

type Coin struct {
	BaseCoinID                      int64   `json:"base_coin_id"`
	CoinId                          int64   `json:"coin_id,omitempty"`
	CoinName                        string  `json:"coin_name"`
	CoinSymbol                      string  `json:"coin_symbol,omitempty"`
	ContractAddress                 string  `json:"contract_address,omitempty"`
	ExplorePath                     string  `json:"explore_path"`
	IconUrl                         string  `json:"icon_url,omitempty"`
	DailyLimitedAcqExchangeQuantity float64 `json:"daily_limited_acq_exchange_quantity"`
	ExchangeFees                    float64 `json:"exchange_fees"`
	IsRechargeable                  bool    `json:"is_rechargeable"`
}

type AppCoin struct {
	AppID int64 `json:"app_id"`
	Coin
}

type DB struct {
	MssqlAccountAll  *basedb.Mssql
	MssqlAccountRead *basedb.Mssql
	Cache            *basedb.CacheV8

	MssqlPointsAll  map[int64]*basedb.Mssql
	MssqlPointsRead map[int64]*basedb.Mssql

	PointDoc    map[string]*MemberPointInfo
	PointDocMtx sync.Mutex

	//PointList     map[int64]PointInfo // 전체 포인트 종류
	ScanPointsMap   map[int64]PointInfo     // 전체 포인트 종류 : key PointId
	ScanPointsOfApp map[int64]*AppPointInfo // 포인트별 app 정보 key pointID

	AppPointsMap map[int64]*AppPointInfo // 전체 app과 포인트 : key AppId

	AppCoins map[int64][]*AppCoin // 전체 app에 속한 CoinID 정보
	Coins    map[int64]*Coin      // 전체 coin 정보 : key coinID

	BaseCoinMapByCoinID map[int64]*context.BaseCoinInfo  // 전체 base coin 정보 : key coin symbol
	BaseCoinMapBySymbol map[string]*context.BaseCoinInfo // 전체 base coin 정보 : key coin symbol
	BaseCoins           context.BaseCoinList

	RedSync *redsync.Redsync
}

type DBType int

const (
	ACCOUNT DBType = iota
	POINT
)

var gDB *DB

func GetDB() *DB {
	return gDB
}

func InitDB(conf *config.ServerConfig) (err error) {
	cache := basedb.GetCacheV8(&conf.Cache)
	gDB = &DB{
		Cache: cache,
	}
	pool := goredis.NewPool(cache.GetDB().RedisClient())
	gDB.RedSync = redsync.New(pool)

	if err := ConnectAllDB(conf); err != nil {
		log.Errorf("InitDB Error : %v", err)
		return err
	}

	// point db create
	gDB.MssqlPointsAll = make(map[int64]*basedb.Mssql)
	gDB.MssqlPointsRead = make(map[int64]*basedb.Mssql)
	var getPointDBs map[int64]*context.PointDB

	getPointDBs, err = gDB.GetPointDatabases()

	if err != nil {
		return err
	} else {
		for _, pointDB := range getPointDBs {
			mssqlDBAll, err := gDB.ConnectDBOfPoint(&conf.MssqlDBPointAll, pointDB)
			if err != nil {
				log.Errorf("err: %v, val: %v, %v, %v, %v",
					err, pointDB.ServerName, conf.MssqlDBPointAll.ID, conf.MssqlDBPointAll.Password, pointDB.DatabaseName)
				return err
			}

			mssqlDBRead, err := gDB.ConnectDBOfPoint(&conf.MssqlDBPointRead, pointDB)
			if err != nil {
				log.Errorf("err: %v, val: %v, %v, %v, %v",
					err, pointDB.ServerName, conf.MssqlDBPointRead.ID, conf.MssqlDBPointRead.Password, pointDB.DatabaseName)
				return err
			}

			gDB.MssqlPointsAll[pointDB.DatabaseID] = mssqlDBAll
			gDB.MssqlPointsRead[pointDB.DatabaseID] = mssqlDBRead
		}
	}

	go func() {
		for {
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			timer.Stop()

			// DB ping 체크 후 오류 시 재 연결
			if db := CheckPingDB(gDB.MssqlAccountAll, conf.MssqlDBAccountAll, ACCOUNT, nil); db != nil {
				gDB.MssqlAccountAll = db
			}

			if db := CheckPingDB(gDB.MssqlAccountRead, conf.MssqlDBAccountRead, ACCOUNT, nil); db != nil {
				gDB.MssqlAccountRead = db
			}

			for _, pointDB := range getPointDBs {
				if db := CheckPingDB(gDB.MssqlPointsAll[pointDB.DatabaseID], conf.MssqlDBPointAll, POINT, pointDB); db != nil {
					gDB.MssqlPointsAll[pointDB.DatabaseID] = db
				}
				if db := CheckPingDB(gDB.MssqlPointsRead[pointDB.DatabaseID], conf.MssqlDBPointRead, POINT, pointDB); db != nil {
					gDB.MssqlPointsRead[pointDB.DatabaseID] = db
				}
			}
		}
	}()

	LoadDBPoint()

	go gDB.ListenSubscribeEvent()

	return nil
}

func LoadDBPoint() {
	gDB.PointDoc = make(map[string]*MemberPointInfo)
	gDB.ScanPointsMap = make(map[int64]PointInfo)
	gDB.ScanPointsOfApp = make(map[int64]*AppPointInfo)
	gDB.AppPointsMap = make(map[int64]*AppPointInfo)
	gDB.AppCoins = make(map[int64][]*AppCoin)
	gDB.Coins = make(map[int64]*Coin)
	gDB.BaseCoinMapByCoinID = make(map[int64]*context.BaseCoinInfo)
	gDB.BaseCoinMapBySymbol = make(map[string]*context.BaseCoinInfo)

	// sequence is important
	gDB.GetPointList()
	gDB.GetAppCoins()
	gDB.GetCoins()
	gDB.GetApps()
	gDB.GetAppPoints()
	gDB.GetBaseCoins()
}

func MakeDbError(resp *base.BaseResponse, errCode int, err error) {
	resp.Return = errCode
	resp.Message = resultcode.ResultCodeText[errCode] + " : " + err.Error()
}

func (o *DB) ConnectDB(conf *baseconf.DBAuth) (*basedb.Mssql, error) {
	port, _ := strconv.ParseInt(conf.Port, 10, 32)
	mssqlDB, err := basedb.NewMssql(conf.Database, "", conf.ID, conf.Password, conf.Host, int(port),
		conf.ApplicationIntent, conf.Timeout, conf.ConnectRetryCount, conf.ConnectRetryInterval)
	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v, %v, %v",
			err, conf.Host, conf.ID, conf.Password, conf.Database, conf.PoolSize, conf.IdleSize)
		return nil, err
	}
	idleSize, _ := strconv.ParseInt(conf.IdleSize, 10, 32)
	mssqlDB.GetDB().SetMaxOpenConns(int(idleSize))
	mssqlDB.GetDB().SetMaxIdleConns(int(idleSize))
	return mssqlDB, nil
}

func (o *DB) ConnectDBOfPoint(conf *baseconf.DBAuth, pointDB *context.PointDB) (*basedb.Mssql, error) {
	port, _ := strconv.ParseInt(conf.Port, 10, 32)
	mssqlDB, err := basedb.NewMssql(pointDB.DatabaseName, "pointDB", conf.ID, conf.Password, pointDB.ServerName, int(port),
		conf.ApplicationIntent, conf.Timeout, conf.ConnectRetryCount, conf.ConnectRetryInterval)

	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v",
			err, pointDB.ServerName, conf.ID, conf.Password, pointDB.DatabaseName)
		return nil, err
	}

	idleSize, _ := strconv.ParseInt(conf.IdleSize, 10, 32)
	mssqlDB.GetDB().SetMaxOpenConns(int(idleSize))
	mssqlDB.GetDB().SetMaxIdleConns(int(idleSize))
	return mssqlDB, nil
}

func ConnectAllDB(conf *config.ServerConfig) error {
	var err error
	gDB.MssqlAccountAll, err = gDB.ConnectDB(&conf.MssqlDBAccountAll)
	if err != nil {
		return err
	}

	gDB.MssqlAccountRead, err = gDB.ConnectDB(&conf.MssqlDBAccountRead)
	if err != nil {
		return err
	}
	return nil
}

func CheckPingDB(db *basedb.Mssql, conf baseconf.DBAuth, dbType DBType, pointDB *context.PointDB) *basedb.Mssql {
	// 연결이 안되어있거나, DB Connection이 끊어진 경우에는 재연결 시도
	if db == nil || !db.Connection.IsConnect {
		if dbType == ACCOUNT {
			newDB, err := gDB.ConnectDB(&conf)
			if err == nil {
				log.Errorf("ACCOUNT connect DB OK")
			}
			return newDB
		} else if dbType == POINT {
			newDB, err := gDB.ConnectDBOfPoint(&conf, pointDB)
			if err == nil {
				log.Errorf("POINT connect DB OK")
			}
			return newDB
		}
	}

	// 연결이 되어있는 상태면 ping
	if db.Connection.IsConnect {
		if err := db.GetDB().Ping(); err != nil {
			// 재시도 횟수
			db.Connection.RetryCount += 1
			if dbType == ACCOUNT {
				log.Errorf("%v DB Ping err RetryCount(%v)", conf.Database, db.Connection.RetryCount)
			} else {
				log.Errorf("%v DB Ping err RetryCount(%v)", pointDB.DatabaseName, db.Connection.RetryCount)
			}

			// ping 2회 시도해도 안되면 close
			if db.Connection.RetryCount >= 2 {
				db.Connection.IsConnect = false
				// DB Close
				if err = db.GetDB().Close(); err == nil {
					if dbType == ACCOUNT {
						log.Errorf("%v DB Closed", conf.Database)
					} else {
						log.Errorf("%v DB Closed", pointDB.DatabaseName)
					}
				}
			}
		}
	}
	return nil
}
