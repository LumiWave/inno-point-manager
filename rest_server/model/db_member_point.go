package model

import (
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
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
				return
			}

			key := MakePointKey(o.MUID)
			//3. redis read
			pointInfo, err := GetDB().GetCachePoint(key)
			if err != nil {
				unLock() // redis unlock
				return
			}
			//4. myuuid check else go func end
			if !strings.EqualFold(o.MyUuid, pointInfo.MyUuid) {
				unLock() // redis unlock
				return
			}
			//5. db update
			for _, point := range pointInfo.Points {
				if err := GetDB().UpdateAppPoint(pointInfo.MUID, point.PointID, point.Quantity, pointInfo.DatabaseID); err != nil {
					unLock() // redis unlock
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
