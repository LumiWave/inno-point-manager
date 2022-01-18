package model

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
)

type DaliyAppPointUnit struct {
	AppId             int64
	PointId           int64
	BackUpCurQuantity map[int64]int64 `json:"backup_current_quantity"`
}

func NewDaliyAppPointUnit(appId, pointId int64) *DaliyAppPointUnit {
	daliyAppPointUnit := &DaliyAppPointUnit{
		AppId:   appId,
		PointId: pointId,
	}

	daliyAppPointUnit.UpdateRun()

	return daliyAppPointUnit
}

func (o *DaliyAppPointUnit) UpdateRun() {
	go func() {

		defer func() {
			//key := MakeMemberPointListKey(o.MUID)
			//delete(GetDB().PointDoc, key)
		}()

		for {
			timer := time.NewTimer(10 * time.Second)
			<-timer.C
			timer.Stop()

			//2. redis lock
			Lockkey := MakeDailyAppPointLockKey(o.AppId, o.PointId)
			unLock, err := AutoLock(Lockkey)
			if err != nil {
				log.Errorf("redis lock fail [lockkey:%v][err:%v]", Lockkey, err)
				return
			}

			key := MakeDailyAppPointKey(o.AppId, o.PointId)
			//3. redis read
			dailyAppPoint, err := GetDB().GetCacheDailyAppPoint(key)
			if err != nil {
				unLock() // redis unlock
				log.Errorf("GetCacheDailyAppPoint [key:%v][err:%v]", key, err)
				return
			}

			//4. db update
			if dailyAppPoint.AdjustQuantity != 0 || dailyAppPoint.AdjustExchangeQuantity != 0 {
				dailyQuantity, dailyExchangeQuantity, resetData, err := GetDB().UpdateApplicationPoints(dailyAppPoint.AppId, dailyAppPoint.PointId, dailyAppPoint.AdjustQuantity, dailyAppPoint.AdjustExchangeQuantity)
				if err != nil {
					log.Errorf("UpdateApplicationPoints [err:%v]", err)
				} else {
					log.Infof("UpdateApplicationPoints [appId:%v][pointId:%v][adjQuantity:%v][adjExchangeQuantity:%v] res[dailyQuantity:%v][dailyExchangeQuantity:%v][resetDate:%v]",
						dailyAppPoint.AppId, dailyAppPoint.PointId, dailyAppPoint.AdjustQuantity, dailyAppPoint.AdjustExchangeQuantity, dailyQuantity, dailyExchangeQuantity, resetData)
				}
				// 7. redis clear
				dailyAppPoint.AdjustQuantity = 0
				dailyAppPoint.AdjustExchangeQuantity = 0
				if err := GetDB().SetCacheDailyAppPoint(key, dailyAppPoint); err != nil {
					log.Errorf("SetCacheDailyAppPoint [err:%v]", err)
				}
			}

			unLock() // redis unlock
		}
	}()
}
