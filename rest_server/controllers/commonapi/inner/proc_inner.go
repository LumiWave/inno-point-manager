package inner

import (
	"errors"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	uuid "github.com/satori/go.uuid"
)

// var gPointDoc map[string]MemberPoint

// type MemberPoint struct {
// 	MUID  int64
// 	AppId int64
// }

func UpdateAppPoint(req *context.ReqPointAppUpdate) (*context.Point, error) {
	// 1. redis lock
	Lockkey := model.MakePointLockKey(req.MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

	// 1-1. redis unlock
	defer unLock()

	respPoint := new(context.Point)
	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakePointKey(req.MUID)
	pointInfo, err := model.GetDB().GetCachePoint(key)
	if err != nil {
		// redis에 존재 하지 않으면 로그인 유저가 로그인 하지 않았다고 판단 하고 에러 리턴
		//return nil, err
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointApp(req.MUID, req.DatabaseID); err != nil {
			return nil, err
		} else {
			find := false
			findIdx := 0
			for idx, point := range points {
				if point.PointID == req.PointID {
					if point.Quantity == req.LastQuantity { // last 수량 비교
						points[idx].Quantity += req.ChangeQuantity
						find = true
						findIdx = idx
					} else {
						err = errors.New("not equal lastest quantity")
					}
					break
				}
			}

			if err != nil {
				log.Errorf("%v ", err)
				return nil, err
			}
			if !find { // point id를 못찼았을경우
				err = errors.New("invalid point ID")
				log.Errorf("%v ", err)
				return nil, err
			}

			respPoint = points[findIdx]

			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: req.DatabaseID,

				MUID:   req.MUID,
				Points: points,
			}

			// 2-2. redis 에 write
			if err := model.GetDB().SetCachePoint(key, pointInfo); err != nil {
				return nil, err
			}

			// 2-3. redis update thread 생성
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo)
		}
	} else {
		// redis 에 존재하면 업데이트
		points := pointInfo.Points

		err = nil
		find := false
		findIdx := 0
		for idx, point := range points {
			if point.PointID == req.PointID {
				if point.Quantity == req.LastQuantity { // last 수량 비교
					points[idx].Quantity = req.LastQuantity + req.ChangeQuantity
					find = true
					findIdx = idx
				} else {
					err = errors.New("not equal lastest quantity")
				}
				break
			}
		}
		if err != nil {
			log.Errorf("%v ", err)
			return nil, err
		}
		if !find { // point id를 못찼았을경우
			err = errors.New("invalid point ID")
			log.Errorf("%v ", err)
			return nil, err
		}

		pointInfo.Points = points
		if err := model.GetDB().SetCachePoint(key, pointInfo); err != nil {
			return nil, err
		}
		respPoint = points[findIdx]
	}

	return respPoint, nil
}

func LoadPoint(MUID, DatabaseID int64) (*context.PointInfo, error) {
	// 1. redis lock
	Lockkey := model.MakePointLockKey(MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

	// 1-1. redis unlock
	defer unLock()

	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakePointKey(MUID)
	pointInfo, err := model.GetDB().GetCachePoint(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointApp(MUID, DatabaseID); err != nil {
			return nil, err
		} else {
			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: DatabaseID,

				MUID:   MUID,
				Points: points,
			}

			// 2-2. redis 에 write
			if err := model.GetDB().SetCachePoint(key, pointInfo); err != nil {
				return nil, err
			}

			// 2-3. redis update thread 생성
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo)
		}
	} else {
		// redis에 존재 한다면 내가 관리하는 thread check, 내 관리가 아니면 그냥 값만 리턴
	}

	return pointInfo, nil
}
