package model

import (
	contextR "context"
	"errors"
	"fmt"
	"strconv"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_ExchangePointToCoinTiers          = "[dbo].[USPAU_Scan_ExchangePointToCoinTiers]"
	USPAU_Scan_ExchangePointToCoinTierConditions = "[dbo].[USPAU_Scan_ExchangePointToCoinTierConditions]"
)

func (o *DB) USPAU_Scan_ExchangePointToCoinTiers() error {
	var returnValue orginMssql.ReturnStatus
	proc := USPAU_Scan_ExchangePointToCoinTiers
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), proc,
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%v QueryContext error : %v", proc, err)
		return nil
	}

	o.SwapP2CTiers = []*context.SwapP2CTier{}
	o.SwapP2CTiersMap = make(map[string]map[int64]*context.SwapP2CTier)

	for rows.Next() {
		tier := &context.SwapP2CTier{}

		if err := rows.Scan(
			&tier.FromID,
			&tier.ToID,
			&tier.TierID,
			&tier.MinimumExchangeQuantity,
			&tier.ExchangeRatio); err != nil {
			log.Errorf("USPAU_Scan_ExchangePointToCoinTiers Scan error : %v", err)
			return err
		} else {
			o.SwapP2CTiers = append(o.SwapP2CTiers, tier)
			key := fmt.Sprintf("%v_%v", tier.FromID, tier.ToID)
			if o.SwapP2CTiersMap[key] == nil {
				o.SwapP2CTiersMap[key] = make(map[int64]*context.SwapP2CTier)
			}
			o.SwapP2CTiersMap[key][tier.TierID] = tier
		}
	}

	if returnValue != 1 {
		log.Errorf("%v returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return nil
}

func (o *DB) USPAU_Scan_ExchangePointToCoinTierConditions() error {
	var returnValue orginMssql.ReturnStatus
	proc := USPAU_Scan_ExchangePointToCoinTierConditions
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), proc,
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%v QueryContext error : %v", proc, err)
		return nil
	}

	o.SwapP2CTierConditions = []*context.SwapP2CTierCondition{}
	o.SwapP2CTierConditionsMap = make(map[string]map[int64][]*context.SwapP2CTierCondition)

	for rows.Next() {
		tier := &context.SwapP2CTierCondition{}

		if err := rows.Scan(
			&tier.FromID,
			&tier.ToID,
			&tier.TierID,
			&tier.ConditionType,
			&tier.ConditionID,
			&tier.Quantity); err != nil {
			log.Errorf("USPAU_Scan_ExchangePointToCoinTierConditions Scan error : %v", err)
			return err
		} else {
			o.SwapP2CTierConditions = append(o.SwapP2CTierConditions, tier)
			key := fmt.Sprintf("%v_%v", tier.FromID, tier.ToID)
			if o.SwapP2CTierConditionsMap[key] == nil {
				o.SwapP2CTierConditionsMap[key] = make(map[int64][]*context.SwapP2CTierCondition)
			}
			o.SwapP2CTierConditionsMap[key][tier.TierID] = append(o.SwapP2CTierConditionsMap[key][tier.TierID], tier)
		}
	}

	if returnValue != 1 {
		log.Errorf("%v returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return nil
}
