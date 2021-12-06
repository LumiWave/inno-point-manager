package inner

import (
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	uuid "github.com/satori/go.uuid"
)

var gPointDoc map[string]MemberPoint

type MemberPoint struct {
	CUID  string
	AppId int64
}

func LoadPoint(CUID string, appid, databaseid int64) (*context.PointInfo, error) {
	// 1. redis lock
	key := model.MakePointLockKey(CUID, appid)
	unLock, err := model.AutoLock(key)
	if err != nil {
		return nil, err
	}

	// 1-1. redis unlock
	defer unLock()

	// 2. redis에 해당 포인트 정보 존재하는지 check
	pointInfo, err := model.GetDB().GetPoint(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointMember(CUID, appid, databaseid); err != nil {
			return nil, err
		} else {
			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: databaseid,

				CUID:   CUID,
				AppID:  appid,
				Points: points,
			}

			// 2-2. redis 에 write
			if err := model.GetDB().SetPoint(key, pointInfo); err != nil {
				return nil, err
			}

			// 2-3. redis update thread 생성
			model.GetDB().PointDoc[key] = &model.MemberPointInfo{
				PointInfo: pointInfo,
			}
			model.GetDB().PointDoc[key].UpdateRun()
		}
	} else {
		// redis에 존재 한다면 내가 관리하는 thread check, 내 관리가 아니면 그냥 값만 리턴
	}

	return pointInfo, nil
}
