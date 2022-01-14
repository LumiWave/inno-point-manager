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

func UpdateAppPoint(req *context.ReqPointAppUpdate, appId int64) (*context.Point, error) {
	// 1. redis lock
	Lockkey := model.MakeMemberPointListLockKey(req.MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

	// if value, ok := model.GetDB().AppPointsMap[appId]; ok {

	// }

	// 1-1. redis unlock
	defer unLock()

	respPoint := new(context.Point)
	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakeMemberPointListKey(req.MUID)
	pointInfo, err := model.GetDB().GetCacheMemberPointList(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointAppList(req.MUID, req.DatabaseID); err != nil {
			return nil, err
		} else {
			// 2-1-1. Account points 로드
			if accountPoint, err := model.GetDB().GetListAccountPoints(0, req.MUID); err != nil {
				return nil, err
			} else {
				// merge
				for _, point := range points {
					if val, ok := accountPoint[point.PointID]; ok {
						point.DailyQuantity = val.DailyQuantity
						point.ResetDate = val.ResetDate
					}
				}
			}

			find := false
			findIdx := 0
			for idx, point := range points {
				if point.PointID == req.PointID {
					if point.Quantity == req.PreQuantity { // last 수량 비교
						//points[idx].PreQuantity += req.PreQuantity
						if points[idx].PreQuantity == 0 {
							points[idx].PreQuantity = req.PreQuantity
						}
						points[idx].AdjustQuantity += req.AdjustQuantity
						points[idx].Quantity += req.AdjustQuantity

						find = true
						findIdx = idx
					} else {
						err = errors.New("not equal previous quantity")
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
			if err := model.GetDB().SetCacheMemberPointList(key, pointInfo); err != nil {
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
				if point.Quantity == req.PreQuantity { // last 수량 비교
					if points[idx].PreQuantity == 0 {
						points[idx].PreQuantity = req.PreQuantity
					}
					points[idx].AdjustQuantity += req.AdjustQuantity
					points[idx].Quantity += req.AdjustQuantity
					find = true
					findIdx = idx
				} else {
					err = errors.New("not equal previous quantity")
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
		if err := model.GetDB().SetCacheMemberPointList(key, pointInfo); err != nil {
			return nil, err
		}
		respPoint = points[findIdx]
	}

	return respPoint, nil
}

func LoadPointList(MUID, DatabaseID int64) (*context.PointInfo, error) {
	// 1. redis lock
	Lockkey := model.MakeMemberPointListLockKey(MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

	// 1-1. redis unlock
	defer unLock()

	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakeMemberPointListKey(MUID)
	pointInfo, err := model.GetDB().GetCacheMemberPointList(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointAppList(MUID, DatabaseID); err != nil {
			return nil, err
		} else {
			// 2-1-1. Account points 로드
			if accountPoint, err := model.GetDB().GetListAccountPoints(0, MUID); err != nil {
				return nil, err
			} else {
				// merge
				for _, point := range points {
					if val, ok := accountPoint[point.PointID]; ok {
						point.DailyQuantity = val.DailyQuantity
						point.ResetDate = val.ResetDate
					}
				}
			}

			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: DatabaseID,

				MUID:   MUID,
				Points: points,
			}

			// 2-2. redis 에 write
			if err := model.GetDB().SetCacheMemberPointList(key, pointInfo); err != nil {
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

func LoadPoint(MUID, PointID, DatabaseID int64) (*context.PointInfo, error) {
	// 1. redis lock
	Lockkey := model.MakeMemberPointListLockKey(MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

	// 1-1. redis unlock
	defer unLock()

	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakeMemberPointListKey(MUID)
	pointInfo, err := model.GetDB().GetCacheMemberPointList(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointAppList(MUID, DatabaseID); err != nil {
			return nil, err
		} else {
			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: DatabaseID,

				MUID: MUID,
			}

			// 2-1-1. Account points 로드
			if accountPoint, err := model.GetDB().GetListAccountPoints(0, MUID); err != nil {
				return nil, err
			} else {
				// merge
				for _, point := range points {
					if val, ok := accountPoint[point.PointID]; ok {
						point.DailyQuantity = val.DailyQuantity
						point.ResetDate = val.ResetDate
					}
				}
			}

			pointInfo = &context.PointInfo{
				MyUuid:     uuid.NewV4().String(),
				DatabaseID: DatabaseID,

				MUID:   MUID,
				Points: points,
			}

			// 2-2. redis 에 write
			if err := model.GetDB().SetCacheMemberPointList(key, pointInfo); err != nil {
				return nil, err
			}

			// 2-3. redis update thread 생성
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo)

			// 2-4. 요청한 point id만 응답해준다.
			temp := context.Point{}
			bFind := false
			for _, point := range pointInfo.Points {
				if point.PointID == PointID {
					temp = *point
					bFind = true
					break
				}
			}

			pointInfo.Points = []*context.Point{}
			if bFind {
				pointInfo.Points = append(pointInfo.Points, &temp)
			}
		}
	} else {
		// redis에 존재 한다면 내가 관리하는 thread check
		// 요청한 point id만 응답해준다.
		temp := context.Point{}
		bFind := false
		for _, point := range pointInfo.Points {
			if point.PointID == PointID {
				temp = *point
				bFind = true
				break
			}
		}

		pointInfo.Points = []*context.Point{}
		if bFind {
			pointInfo.Points = append(pointInfo.Points, &temp)
		}
	}

	return pointInfo, nil
}
