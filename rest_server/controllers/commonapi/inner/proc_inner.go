package inner

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
	uuid "github.com/satori/go.uuid"
)

func UpdateAppPoint(req *context.ReqPointAppUpdate, appId int64) (*context.Point, error) {
	// 1. redis lock
	Lockkey := model.MakeMemberPointListLockKey(req.MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		return nil, err
	}

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
						point.TodayQuantity = val.TodayLimitedQuantity
						if t, err := time.Parse("2006-01-02T15:04:05Z", val.ResetDate); err != nil {
							log.Errorf("time.Parse [err%v]", err)
						} else {
							point.ResetDate = t.Format("2006-01-02")
						}
					}
				}
			}

			find := false
			findIdx := 0
			for idx, point := range points {
				if point.PointID == req.PointID {
					if point.Quantity == req.PreQuantity { // last 수량 비교
						//일일 최대 적립 포인트량 비교
						enable := checkTodayPoint(point, appId, &req.AdjustQuantity)
						if !enable {
							err = errors.New(resultcode.ResultCodeText[resultcode.Result_Error_Exceeded_TodayPoints_earned])
							break
						}

						// if points[idx].PreQuantity == 0 {
						// 	points[idx].PreQuantity = req.PreQuantity
						// }

						points[idx].AdjustQuantity += req.AdjustQuantity
						points[idx].Quantity += req.AdjustQuantity

						find = true
						findIdx = idx
					} else {
						err = errors.New(resultcode.ResultCodeText[resultcode.Result_Error_NotEqual_PreviousQuantity])
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
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo, appId, false)
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
					//일일 최대 적립 포인트량 비교
					enable := checkTodayPoint(point, appId, &req.AdjustQuantity)
					if !enable {
						err = errors.New(resultcode.ResultCodeText[resultcode.Result_Error_Exceeded_TodayPoints_earned])
						break
					}

					// if points[idx].PreQuantity == 0 {
					// 	points[idx].PreQuantity = req.PreQuantity
					// }
					points[idx].AdjustQuantity += req.AdjustQuantity
					points[idx].Quantity += req.AdjustQuantity
					find = true
					findIdx = idx
				} else {
					err = errors.New(resultcode.ResultCodeText[resultcode.Result_Error_NotEqual_PreviousQuantity])
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

func LoadPointList(MUID, DatabaseID, appId int64) (*context.PointInfo, error) {
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
						point.TodayQuantity = val.TodayLimitedQuantity
						if t, err := time.Parse("2006-01-02T15:04:05Z", val.ResetDate); err != nil {
							log.Errorf("time.Parse [err%v]", err)
						} else {
							point.ResetDate = t.Format("2006-01-02")
						}
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
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo, appId, true)
		}
	} else {
		// redis에 존재 한다면 내가 관리하는 thread check, 내 관리가 아니면 그냥 값만 리턴
	}

	return pointInfo, nil
}

func LoadPoint(MUID, PointID, DatabaseID, appId int64) (*context.PointInfo, error) {
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
						point.TodayQuantity = val.TodayLimitedQuantity
						if t, err := time.Parse("2006-01-02T15:04:05Z", val.ResetDate); err != nil {
							log.Error(err)
						} else {
							point.ResetDate = t.Format("2006-01-02")
						}
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
			model.GetDB().PointDoc[key] = model.NewMemberPointInfo(pointInfo, appId, true)

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

func Swap(params *context.ReqSwapInfo) *base.BaseResponse {

	// 보유 포인트 수량 가져오기
	// if memberInfo, err := point_manager_server.GetInstance().GetPointAppList(member.MUID, member.DatabaseID); err == nil {
	// 	for _, point := range memberInfo.Points {
	// 		if point.PointID == swapInfo.PointID {
	// 			swapInfo.PreviousPointQuantity = point.Quantity
	// 			swapInfo.PointQuantity = swapInfo.PreviousPointQuantity + swapInfo.AdjustPointQuantity
	// 			break
	// 		}
	// 	}

	// 	// 코인으로 전환시 0이 아닌 수량을 가지고 있는지 확인
	// 	if swapInfo.PreviousCoinQuantity == 0 {
	// 		// 내 포인트를 찾지 못함
	// 		log.Errorf("not find me point id [point_id:%v]", swapInfo.PointID)
	// 		resp.SetReturn(resultcode.Result_Invalid_PointID_Error)
	// 		return ctx.EchoContext.JSON(http.StatusOK, resp)
	// 	}
	// } else {
	// 	// point manager server 호출 에러
	// 	log.Errorf("point_manager_server GetPointAppList error : %v", err)
	// 	return ctx.EchoContext.JSON(http.StatusOK, base.BaseResponseInternalServerError())
	// }

	// // 코인으로 전환은 point의 최소 수량만 확인
	// if appPointInfo, ok := model.GetDB().AppPointsMap[swapInfo.AppID]; ok {
	// 	for _, point := range appPointInfo.Points {
	// 		if point.PointId == swapInfo.PointID {
	// 			if swapInfo.PreviousPointQuantity < point.MinExchangeQuantity || swapInfo.PreviousPointQuantity < swapInfo.AdjustPointQuantity{
	// 				resp.SetReturn(resultcode.Result_Invalid_PointQuantity_Error)
	// 				return ctx.EchoContext.JSON(http.StatusOK, resp)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	// 존재하지 않는 appid를 전환 시도 에러
	// 	resp.SetReturn(resultcode.Result_Invalid_AppID_Error)
	// 	return ctx.EchoContext.JSON(http.StatusOK, resp)
	// }

	resp := new(base.BaseResponse)
	resp.Success()

	// LogID, PreviousPointQuantity, PointQuantity 정보를 찾아서 params에 추가 해줘야 함
	// point정보는 redis lock을 걸고 조회 해야 무결성이 유지됨

	// 1. redis lock
	Lockkey := model.MakeMemberPointListLockKey(params.MUID)
	unLock, err := model.AutoLock(Lockkey)
	if err != nil {
		resp.SetReturn(resultcode.Result_RedisError_Lock_fail)
		return resp
	}

	// 1-1. redis unlock
	defer unLock()

	// 2. redis에 해당 포인트 정보 존재하는지 check
	key := model.MakeMemberPointListKey(params.MUID)
	mePointInfo, err := model.GetDB().GetCacheMemberPointList(key)
	if err != nil {
		// 2-1. redis에 존재하지 않는다면 db에서 로드
		if points, err := model.GetDB().GetPointAppList(params.MUID, params.DatabaseID); err != nil {
			log.Errorf("GetPointAppList error : %v", err)
			resp.SetReturn(resultcode.Result_Error_DB_GetPointAppList)
			return resp
		} else {
			for _, point := range points {
				if point.PointID == params.PointID {
					params.PreviousPointQuantity = point.Quantity
					params.PointQuantity = params.PreviousPointQuantity + params.AdjustPointQuantity
					break
				}
			}

		}
	} else {
		// redis에 존재 한다면 강제로 db에 먼저 write
		for _, point := range mePointInfo.Points {
			var eventID context.EventID_type
			if point.AdjustQuantity >= 0 {
				eventID = context.EventID_add
			} else {
				eventID = context.EventID_sub
			}
			if todayLimitedQuantity, resetDate, err := model.GetDB().UpdateAppPoint(mePointInfo.DatabaseID, mePointInfo.MUID, point.PointID,
				point.PreQuantity, point.AdjustQuantity, point.Quantity, context.LogID_cp, eventID); err != nil {
				log.Errorf("UpdateAppPoint error : %v", err)
				resp.SetReturn(resultcode.Result_Error_DB_UpdateAppPoint)
				return resp
			} else {
				//현재 일일 누적량, 날짜 업데이트
				point.TodayQuantity = todayLimitedQuantity
				point.ResetDate = resetDate

				point.AdjustQuantity = 0
				point.PreQuantity = point.Quantity

				// swap point quantity에 업데이트
				if params.PointID == point.PointID && params.MUID == mePointInfo.MUID {
					params.PreviousPointQuantity = point.Quantity
					params.PointQuantity = params.PreviousPointQuantity + params.AdjustPointQuantity
				}
			}
		}
	}

	pointInfo := model.GetDB().AppPointsMap[params.AppID].PointsMap[params.PointID]
	if params.EventID == context.EventID_toCoin {
		// 코인으로 전환시 체크
		// 포인트 보유수량이 전환량 보다 큰지 확인
		if params.PreviousPointQuantity <= 0 || // 보유 포인트량이 0일경우
			params.PreviousPointQuantity < params.AdjustPointQuantity || // 전환 할 수량보다 보유 수량이 적을 경우
			pointInfo.MinExchangeQuantity > params.AdjustPointQuantity { // 전환 최소 수량 에러
			// 전환할 포인트 수량이 없음 에러
			log.Errorf("not find me point id [point_id:%v][PointQuantity:%v]", params.PointID, params.PreviousPointQuantity)
			resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
			return resp
		}
		// 전환 비율 계산 후 타당성 확인
		exchangeCoin := float64(params.AdjustPointQuantity) + float64(params.AdjustPointQuantity)*pointInfo.ExchangeRatio
		exchangeCoin = toFixed(exchangeCoin, 4)
		if params.AdjustCoinQuantity != exchangeCoin {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return resp
		}

	} else if params.EventID == context.EventID_toPoint {
		// 코인 보유 수량이 전환량 보다 큰지 확인
		if params.PreviousCoinQuantity <= 0 || // 보유 코인량이 0인경우
			params.PreviousCoinQuantity < params.AdjustCoinQuantity {
			log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_MinCoinQuantity]+" [coin_id:%v][coin_quantity:%v]", params.CoinID, params.PreviousCoinQuantity)
			resp.SetReturn(resultcode.Result_Error_MinCoinQuantity)
			return resp
		}
		// 전환 비율 계산 후 타당성 확인
		exchangePoint := params.AdjustCoinQuantity * pointInfo.ExchangeRatio
		exchangePoint = toFixed(exchangePoint, 0)
		if params.AdjustPointQuantity != int64(exchangePoint) {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToCoin)
			return resp
		}
	}

	params.LogID = context.LogID_exchange

	// swap 후에 redis 삭제
	if err := model.GetDB().PostPointCoinSwap(params); err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
	}

	model.GetDB().DelCacheMemberPointList(key)

	return nil
}

func checkTodayPoint(point *context.Point, appId int64, reqAdjustQuantity *int64) bool {
	if strings.EqualFold(point.ResetDate, time.Now().Format("2006-01-02")) { // 날짜가 바뀌었는지 체크
		if point.TodayQuantity+point.AdjustQuantity >= model.GetDB().AppPointsMap[appId].PointsMap[point.PointID].DaliyLimitedQuantity {
			// 이미 다 채운 상태라면 에러 리턴
			return false
		} else {
			if point.TodayQuantity + +point.AdjustQuantity + *reqAdjustQuantity > model.GetDB().AppPointsMap[appId].PointsMap[point.PointID].DaliyLimitedQuantity {
				// 초과시 가능 포인트만 적립
				*reqAdjustQuantity = model.GetDB().AppPointsMap[appId].PointsMap[point.PointID].DaliyLimitedQuantity - point.TodayQuantity
			}
		}
	}

	return true
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
