package model

import (
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
)

type MemberPointInfo struct {
	*context.PointInfo
	AppId             int64
	BackUpCurQuantity map[int64]int64 `json:"backup_current_quantity"`
}

func NewMemberPointInfo(pointInfo *context.PointInfo, appId int64, load bool) *MemberPointInfo {
	memberPointInfo := &MemberPointInfo{
		PointInfo: pointInfo,
		AppId:     appId,
	}

	memberPointInfo.BackUpCurQuantity = make(map[int64]int64)
	for _, point := range memberPointInfo.Points {
		if load {
			memberPointInfo.BackUpCurQuantity[point.PointID] = point.Quantity
		} else {
			memberPointInfo.BackUpCurQuantity[point.PointID] = 0
		}
	}

	memberPointInfo.UpdateRun()

	return memberPointInfo
}

func (o *MemberPointInfo) UpdateRun() {
	go func() {

		defer func() {
			key := MakeMemberPointListKey(o.MUID)
			delete(GetDB().PointDoc, key)
		}()

		for {
			timer := time.NewTimer(10 * time.Second)
			<-timer.C

			//2. redis lock
			Lockkey := MakeMemberPointListLockKey(o.MUID)
			unLock, err := AutoLock(Lockkey)
			if err != nil {
				log.Errorf("redis lock fail [lockkey:%v][err:%v]", Lockkey, err)
				return
			}

			key := MakeMemberPointListKey(o.MUID)
			//3. redis read
			pointInfo, err := GetDB().GetCacheMemberPointList(key)
			if err != nil {
				unLock() // redis unlock
				log.Errorf("GetCacheMemberPointList [key:%v][err:%v]", key, err)
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
				if o.BackUpCurQuantity[point.PointID] != point.Quantity { // 포인트 정보가 변경된 경우에만 db 업데이트 처리
					var eventID context.EventID_type
					if point.AdjustQuantity >= 0 {
						eventID = context.EventID_add
					} else {
						eventID = context.EventID_sub
					}
					if todayLimitedQuantity, resetDate, err := GetDB().UpdateAppPoint(pointInfo.DatabaseID, pointInfo.MUID, point.PointID,
						point.PreQuantity, point.AdjustQuantity, point.Quantity, context.LogID_cp, eventID); err != nil {
						unLock() // redis unlock
						log.Errorf("UpdateAppPoint [err:%v]", err)
					} else {
						// 업데이트 성공시 BackUpCurQuantity 최신으로 업데이트
						o.BackUpCurQuantity[point.PointID] = point.Quantity

						//현재 일일 누적량, 날짜 업데이트
						point.TodayQuantity = todayLimitedQuantity
						point.ResetDate = resetDate

						point.AdjustQuantity = 0
						point.PreQuantity = point.Quantity
					}
				}
			}

			//6. local save
			o.PointInfo = pointInfo

			// 7. redis 에 write
			if err := GetDB().SetCacheMemberPointList(key, pointInfo); err != nil {
				log.Errorf("SetCacheMemberPointList [err:%v]", err)
			}

			timer.Stop()
			unLock() // redis unlock
		}
	}()
}
