package model

import (
	originCtx "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_DatabaseServers   = "[dbo].[USPAU_Scan_DatabaseServers]"
	USPAU_Scan_Points            = "[dbo].[USPAU_Scan_Points]"
	USPAU_Scan_ApplicationCoins  = "[dbo].[USPAU_Scan_ApplicationCoins]"
	USPAU_Scan_ApplicationPoints = "[dbo].[USPAU_Scan_ApplicationPoints]"
	USPAU_Scan_Applications      = "[dbo].[USPAU_Scan_Applications]"
	USPAU_Scan_Coins             = "[dbo].[USPAU_Scan_Coins]"
	USPAU_Scan_BaseCoins         = "[dbo].[USPAU_Scan_BaseCoins]"
)

// point database 리스트 요청
func (o *DB) GetPointDatabases() (map[int64]*context.PointDB, error) {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_DatabaseServers, &rs)

	if err != nil {
		log.Errorf("USPAU_Scan_DatabaseServers QueryContext error : %v", err)
		return nil, err
	}

	defer rows.Close()

	pointdbs := make(map[int64]*context.PointDB)

	pointdb := new(context.PointDB)
	for rows.Next() {
		rows.Scan(&pointdb.DatabaseID, &pointdb.DatabaseName, &pointdb.ServerName)
		pointdbs[pointdb.DatabaseID] = pointdb
	}

	return pointdbs, nil
}

// point 전체 list
func (o *DB) GetPointList() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_Points, &rs)
	if err != nil {
		log.Errorf("USPAU_Scan_Points QueryContext error : %v", err)
		return err
	}

	defer rows.Close()

	o.ScanPointsMap = make(map[int64]PointInfo)

	var pointId int64
	var pointName, iconPath string
	for rows.Next() {
		if err := rows.Scan(&pointId, &pointName, &iconPath); err == nil {
			info := PointInfo{
				PointId:   pointId,
				PointName: pointName,
				IconUrl:   iconPath,
			}
			o.ScanPointsMap[pointId] = info
		} else {
			log.Warnf("USPAU_Scan_Points Scan err : %v", err)
		}
	}

	return nil
}

// 전체 app coinid list
func (o *DB) GetAppCoins() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_ApplicationCoins, &rs)
	if err != nil {
		log.Errorf("USPAU_Scan_ApplicationCoins QueryContext error : %v", err)
		return err
	}

	defer rows.Close()

	o.AppCoins = make(map[int64][]*AppCoin)

	for rows.Next() {
		appCoin := &AppCoin{}
		if err := rows.Scan(&appCoin.AppID, &appCoin.CoinId, &appCoin.BaseCoinID); err == nil {
			o.AppCoins[appCoin.AppID] = append(o.AppCoins[appCoin.AppID], appCoin)
		} else {
			log.Errorf("USPAU_Scan_ApplicationCoins Scan error : %v", err)
		}
	}

	return nil
}

// 전체 coin info list
func (o *DB) GetCoins() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_Coins, &rs)
	if err != nil {
		log.Errorf("USPAU_Scan_Coins QueryContext error : %v", err)
		return err
	}

	defer rows.Close()

	o.Coins = make(map[int64]*Coin)
	o.CoinsBySymbol = make(map[string]*Coin)

	for rows.Next() {
		coin := &Coin{}
		if err := rows.Scan(&coin.CoinId,
			&coin.BaseCoinID,
			&coin.CoinName,
			&coin.CoinSymbol,
			&coin.ContractAddress,
			&coin.ExplorePath,
			&coin.IconUrl,
			&coin.DailyLimitedAcqExchangeQuantity,
			&coin.ExchangeFees,
			&coin.IsRechargeable); err == nil {
			o.Coins[coin.CoinId] = coin
			o.CoinsBySymbol[coin.CoinSymbol] = coin
		} else {
			log.Errorf("USPAU_Scan_Coins Scan error : %v", err)
		}
	}

	for _, appCoins := range o.AppCoins {
		for _, appCoin := range appCoins {
			for _, coin := range o.Coins {
				if appCoin.CoinId == coin.CoinId {
					appCoin.BaseCoinID = coin.BaseCoinID
					appCoin.CoinName = coin.CoinName
					appCoin.CoinSymbol = coin.CoinSymbol
					appCoin.ContractAddress = coin.ContractAddress
					appCoin.IconUrl = coin.IconUrl
					appCoin.DailyLimitedAcqExchangeQuantity = coin.DailyLimitedAcqExchangeQuantity
					appCoin.ExchangeFees = coin.ExchangeFees
					break
				}
			}
		}
	}

	return nil
}

// 전체 base coin list 조회
func (o *DB) GetBaseCoins() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_BaseCoins, &rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return err
	}

	defer rows.Close()

	o.BaseCoinMapByCoinID = make(map[int64]*context.BaseCoinInfo)
	o.BaseCoinMapBySymbol = make(map[string]*context.BaseCoinInfo)
	o.BaseCoins.Coins = nil
	for rows.Next() {
		baseCoin := &context.BaseCoinInfo{}
		if err := rows.Scan(&baseCoin.BaseCoinID, &baseCoin.BaseCoinName, &baseCoin.BaseCoinSymbol, &baseCoin.IsUsedParentWallet); err == nil {
			o.BaseCoinMapByCoinID[baseCoin.BaseCoinID] = baseCoin
			o.BaseCoinMapBySymbol[baseCoin.BaseCoinSymbol] = baseCoin
			o.BaseCoins.Coins = append(o.BaseCoins.Coins, baseCoin)
		}
	}

	return nil
}

// 전체 app list 조회
func (o *DB) GetApps() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(originCtx.Background(), USPAU_Scan_Applications, &rs)
	if err != nil {
		log.Errorf("USPAU_Scan_Applications QueryContext error : %v", err)
		return err
	}

	defer rows.Close()

	o.AppPointsMap = make(map[int64]*AppPointInfo)
	for rows.Next() {
		appInfo := &AppPointInfo{}
		if err := rows.Scan(&appInfo.AppId, &appInfo.AppName, &appInfo.IconUrl,
			&appInfo.GooglePlayPath, &appInfo.AppleStorePath, &appInfo.BrandingPagePath); err == nil {
			o.AppPointsMap[appInfo.AppId] = appInfo
			o.AppPointsMap[appInfo.AppId].PointsMap = make(map[int64]*PointInfo)
		} else {
			log.Errorf("USPAU_Scan_Applications Scan error : %v", err)
		}
	}

	return nil
}

// 전체 app 과 포인트 list 조회
func (o *DB) GetAppPoints() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.GetDB().QueryContext(originCtx.Background(), USPAU_Scan_ApplicationPoints, &rs)
	if err != nil {
		log.Errorf("USPAU_Scan_ApplicationPoints QueryContext error : %v", err)
		return err
	}

	defer rows.Close()

	o.ScanPointsOfApp = make(map[int64]*AppPointInfo)

	var appId, pointId, minExchangeQuantity, daliyLimiteAcqQuantity, dailyLimitedAcqExchangeQuantity sql.NullInt64
	var exchangeRatio sql.NullFloat64
	for rows.Next() {
		if err := rows.Scan(&appId, &pointId, &minExchangeQuantity, &exchangeRatio, &daliyLimiteAcqQuantity, &dailyLimitedAcqExchangeQuantity); err == nil {
			temp := o.ScanPointsMap[pointId.Int64]
			temp.ExchangeRatio = exchangeRatio.Float64
			temp.MinExchangeQuantity = minExchangeQuantity.Int64
			temp.DaliyLimitedAcqQuantity = daliyLimiteAcqQuantity.Int64
			temp.DailyLimitedAcqExchangeQuantity = dailyLimitedAcqExchangeQuantity.Int64

			o.AppPointsMap[appId.Int64].Points = append(o.AppPointsMap[appId.Int64].Points, &temp)
			o.AppPointsMap[appId.Int64].PointsMap[pointId.Int64] = &temp
			o.ScanPointsOfApp[pointId.Int64] = o.AppPointsMap[appId.Int64]
		} else {
			log.Errorf("USPAU_Scan_ApplicationPoints Scan error : %v", err)
		}
	}

	return nil
}
