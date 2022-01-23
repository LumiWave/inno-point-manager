package inner

import (
	"math"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

func Swap(params *context.ReqSwapInfo) *base.BaseResponse {
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

			if point.AdjustQuantity != 0 {
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
				}
			} else {
				point.AdjustQuantity = 0
				point.PreQuantity = point.Quantity
			}

			// swap point quantity에 업데이트
			if params.PointID == point.PointID && params.MUID == mePointInfo.MUID {
				params.PreviousPointQuantity = point.Quantity
				params.PointQuantity = params.PreviousPointQuantity + params.AdjustPointQuantity
			}
		}
	}

	pointInfo := model.GetDB().AppPointsMap[params.AppID].PointsMap[params.PointID]
	if params.EventID == context.EventID_toCoin {
		// 코인으로 전환시 체크
		// 포인트 보유수량이 전환량 보다 큰지 확인
		absAdjustPointQuantity := int64(math.Abs(float64(params.AdjustPointQuantity)))
		if params.PreviousPointQuantity <= 0 || // 보유 포인트량이 0일경우
			params.PreviousPointQuantity < params.AdjustPointQuantity || // 전환 할 수량보다 보유 수량이 적을 경우
			pointInfo.MinExchangeQuantity > absAdjustPointQuantity { // 전환 최소 수량 에러
			// 전환할 포인트 수량이 없음 에러
			log.Errorf("not find me point id [point_id:%v][PointQuantity:%v]", params.PointID, params.PreviousPointQuantity)
			resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
			return resp
		}
		// 전환 비율 계산 후 타당성 확인
		exchangeCoin := float64(absAdjustPointQuantity) * pointInfo.ExchangeRatio
		exchangeCoin = toFixed(exchangeCoin, 4)
		if params.AdjustCoinQuantity != exchangeCoin {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return resp
		}

	} else if params.EventID == context.EventID_toPoint {
		// 코인 보유 수량이 전환량 보다 큰지 확인
		absAdjustCoinQuantity := math.Abs(params.AdjustCoinQuantity)
		if params.PreviousCoinQuantity <= 0 || // 보유 코인량이 0인경우
			params.PreviousCoinQuantity < absAdjustCoinQuantity {
			log.Errorf(resultcode.ResultCodeText[resultcode.Result_Error_MinCoinQuantity]+" [coin_id:%v][coin_quantity:%v]", params.CoinID, params.PreviousCoinQuantity)
			resp.SetReturn(resultcode.Result_Error_MinCoinQuantity)
			return resp
		}
		// 전환 비율 계산 후 타당성 확인
		exchangePoint := absAdjustCoinQuantity / pointInfo.ExchangeRatio
		exchangePoint = toFixed(exchangePoint, 0)
		if params.AdjustPointQuantity != int64(exchangePoint) {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToCoin)
			return resp
		}
	}

	// swap 후에 redis 삭제
	if err := model.GetDB().PostPointCoinSwap(params); err != nil {
		resp.SetReturn(resultcode.Result_Error_DB_PostPointCoinSwap)
	}

	model.GetDB().DelCacheMemberPointList(key)
	resp.Value = params
	return resp
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
