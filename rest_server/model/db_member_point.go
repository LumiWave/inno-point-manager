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
	//1. update timeout check
	ticker := time.NewTicker(10000 * time.Millisecond)

	defer func() {
		ticker.Stop()
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				{
					//2. redis lock
					key := MakePointLockKey(o.CUID, o.AppID)
					unLock, err := AutoLock(key)
					if err != nil {
						return
					}

					defer unLock() // redis unlock

					//3. redis read
					pointInfo, err := GetDB().GetPoint(key)
					if err != nil {
						return
					}
					//4. myuuid check else go func end
					if !strings.EqualFold(o.MyUuid, pointInfo.MyUuid) {
						return
					}
					//5. db update
					for _, point := range *pointInfo.Points {
						if err := GetDB().UpdatePoint(pointInfo.CUID, pointInfo.AppID, point.PointID, point.Quantity, pointInfo.DatabaseID); err != nil {
							return
						}
					}

					//6. local save
					o.PointInfo = pointInfo
				}
			}
		}
	}()
}
