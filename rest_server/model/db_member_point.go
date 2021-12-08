package model

import (
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
)

type MemberPointInfo struct {
	*context.PointInfo
}

func (o *MemberPointInfo) UpdateRun() {

	go func() {
		defer func() {
			key := MakePointKey(o.MUID)
			GetDB().PointDoc[key] = nil
		}()

		for {
			timer := time.NewTimer(10 * time.Second)
			<-timer.C

			//2. redis lock
			Lockkey := MakePointLockKey(o.MUID)
			unLock, err := AutoLock(Lockkey)
			if err != nil {
				log.Errorf("redis lock fail [lockkey:%v][err:%v]", Lockkey, err)
				return
			}

			key := MakePointKey(o.MUID)
			//3. redis read
			pointInfo, err := GetDB().GetCachePoint(key)
			if err != nil {
				unLock() // redis unlock
				log.Errorf("GetCachePoint [key:%v][err:%v]", key, err)
				return
			}
			//4. myuuid check else go func end
			if !strings.EqualFold(o.MyUuid, pointInfo.MyUuid) {
				log.Errorf("Myuuid diffrent [my_uuid:%v][cache_uuid:%v]", o.MyUuid, pointInfo.MyUuid)
				unLock() // redis unlock
				return
			}
			//5. db update
			for _, point := range pointInfo.Points {
				if err := GetDB().UpdateAppPoint(pointInfo.MUID, point.PointID, point.Quantity, pointInfo.DatabaseID); err != nil {
					unLock() // redis unlock
					log.Errorf("UpdateAppPoint [err:%v]", err)
					return
				}
			}

			//6. local save
			o.PointInfo = pointInfo

			timer.Stop()
			unLock() // redis unlock
		}
	}()
}
