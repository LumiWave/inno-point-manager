package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-point-manager/rest_server/model"
)

func PostPointCoinSwap(params *context.ReqSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if !model.GetSwapEnable() {
		resp.SetReturn(resultcode.Result_Error_IsSwapMaintenance)
	} else if err := inner.SwapWallet(params, ctx.GetValue().InnoUID); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PutSwapStatus(params *context.ReqSwapStatus, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.PutSwapStatus(params); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetSwapInprogressNotExist(params *context.ReqSwapInprogress, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	log.Debugf("GetSwapInprogressNotExist auid : %v", params.AUID)
	// 내 지갑 정보를 가져와서 모든 지갑을 뒤져버 진행 중에 있는지 캐시 정보로 체크
	swapInfos := []*context.ReqSwapInfo{}
	mapWallet := make(map[string]string)
	if wallets, _, err := model.GetDB().USPAU_GetList_AccountWallets(params.AUID); err == nil {
		for _, wallet := range wallets {
			if _, ok := mapWallet[wallet.WalletAddress]; !ok {
				if swapInfo, err := model.GetDB().CacheGetSwapWallet(wallet.WalletAddress); err == nil {
					swapInfos = append(swapInfos, swapInfo)
					if swapInfo.TxType == context.EventID_P2C {
						mapWallet[swapInfo.SwapToCoin.WalletAddress] = swapInfo.SwapToCoin.WalletAddress
					} else if swapInfo.TxType == context.EventID_C2P || swapInfo.TxType == context.EventID_C2C {
						mapWallet[swapInfo.SwapFromCoin.WalletAddress] = swapInfo.SwapFromCoin.WalletAddress
					}
				}
			}
		}
		if len(swapInfos) > 0 {
			resp.Value = swapInfos
			resp.SetReturn(resultcode.Result_Error_Transfer_Inprogress)
		}
	} else {
		log.Errorf("USPAU_GetList_AccountWallets err : %v, auid:%v", err, params.AUID)
		resp.SetReturn(resultcode.Result_Error_Db_GetAccountWallets)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteDeleteSwapInfo(params *context.DeleteDeleteSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if len(params.WalletAddress) == 0 {
		model.GetDB().CacheDelAllSwapWallet()
	} else if err := model.GetDB().CacheDelSwapWallet(params.WalletAddress); err != nil {
		log.Errorf("CacheDelSwapWallet err:%v, wallet:%v", err, params.WalletAddress)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetGameChipSwapInfo(params *context.ReqGameChipSwapInfo, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	res := &context.ResGameChipSwapInfo{
		ExchangeRatio: config.GetInstance().GameSwap.ExchangeRatio,

		SSRToSSRMID: context.EventID_SSR2SSRM,
		SSRMToSSRID: context.EventID_SSRM2SSR,
	}

	_, members, err := model.GetDB().USPAU_GetList_Members(ctx.GetValue().AUID)
	if err != nil {
		log.Errorf("USPAU_GetList_Members err: %v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// SSR 포인트 조회
	appID := config.GetInstance().GameSwap.AppID
	pointID := config.GetInstance().GameSwap.PointID
	if member, ok := members[appID]; !ok {
		res.SSRQuantity = 0
	} else {
		if pointInfo, err := inner.LoadPoint(member.MUID, pointID, member.DatabaseID, member.AppID); err != nil {
			resp.SetReturn(resultcode.Result_DBError)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		} else {
			if len(pointInfo.Points) == 0 { // 포인트가 존재하지 않는다면 에러
				resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
				return ctx.EchoContext.JSON(http.StatusOK, resp)
			}
			res.SSRQuantity = pointInfo.Points[0].Quantity
		}
	}
	// SSRM chip 조회
	if isBlocked, chipQuantity, err := model.GetDB().USPG_Get_Users_By_InnoUID(ctx.GetValue().InnoUID); err != nil {
		resp.SetReturn(resultcode.Result_Error_NotExistMember)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	} else {
		if isBlocked {
			resp.SetReturn(resultcode.Result_Error_NotExistMember)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		} else {
			res.SSRMQuantity = chipQuantity
		}
	}

	resp.Value = res

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostGameChipSwap(params *context.ReqGameChipSwap, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	exchangeRate := config.GetInstance().GameSwap.ExchangeRatio

	// 수량 검증
	switch params.SwapType {
	case context.EventID_SSR2SSRM:
		if params.SSRAdjustPoint >= 0 || params.SSRMAdjustChip <= 0 {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}

		if int64(absInt64(params.SSRAdjustPoint)*exchangeRate/100) != params.SSRMAdjustChip {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
		// 보유량 확인 : SSR 보유량만 확인 하면 됨
		_, members, err := model.GetDB().USPAU_GetList_Members(ctx.GetValue().AUID)
		if err != nil {
			log.Errorf("USPAU_GetList_Members err: %v", err)
			resp.SetReturn(resultcode.Result_DBError)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
		appID := config.GetInstance().GameSwap.AppID
		pointID := config.GetInstance().GameSwap.PointID
		if member, ok := members[appID]; !ok {
			resp.SetReturn(resultcode.Result_Error_NotExistMember)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		} else {
			if pointInfo, err := inner.LoadPoint(member.MUID, pointID, member.DatabaseID, member.AppID); err != nil {
				resp.SetReturn(resultcode.Result_DBError)
				return ctx.EchoContext.JSON(http.StatusOK, resp)
			} else {
				if len(pointInfo.Points) == 0 { // 포인트가 존재하지 않는다면 에러
					resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				}
				if pointInfo.Points[0].PointID != pointID { // 포인트가 존재하지 않는다면 에러
					resp.SetReturn(resultcode.Result_Error_NotExistMember)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				}
				if pointInfo.Points[0].Quantity < absInt64(params.SSRAdjustPoint) { // 보유 ssr 포인트 수량이 부족하면 에러
					resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				}

				// 수량 업데이트
				if err := ProcGameChipSwap(params, pointID, pointInfo.MUID, pointInfo.DatabaseID, ctx.GetValue().InnoUID, pointInfo.Points[0].Quantity); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				}
			}
		}
	case context.EventID_SSRM2SSR:
		if params.SSRMAdjustChip >= 0 || params.SSRAdjustPoint <= 0 {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
		if int64(absInt64(params.SSRMAdjustChip)*100/exchangeRate) != params.SSRAdjustPoint {
			resp.SetReturn(resultcode.Result_Error_Exchangeratio_ToPoint)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
		// 보유량 확인 : SSRM 만 확인 하면 됨
		if isBlocked, chipQuantity, err := model.GetDB().USPG_Get_Users_By_InnoUID(ctx.GetValue().InnoUID); err != nil {
			resp.SetReturn(resultcode.Result_Error_NotExistMember)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		} else {
			if isBlocked {
				resp.SetReturn(resultcode.Result_Error_NotExistMember)
				return ctx.EchoContext.JSON(http.StatusOK, resp)
			} else {
				if chipQuantity < absInt64(params.SSRMAdjustChip) { // 교환 수량 부족
					resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				}
			}

			_, members, err := model.GetDB().USPAU_GetList_Members(ctx.GetValue().AUID)
			if err != nil {
				log.Errorf("USPAU_GetList_Members err: %v", err)
				resp.SetReturn(resultcode.Result_DBError)
				return ctx.EchoContext.JSON(http.StatusOK, resp)
			}
			appID := config.GetInstance().GameSwap.AppID
			pointID := config.GetInstance().GameSwap.PointID
			if member, ok := members[appID]; !ok {
				resp.SetReturn(resultcode.Result_Error_NotExistMember)
				return ctx.EchoContext.JSON(http.StatusOK, resp)
			} else {
				if pointInfo, err := inner.LoadPoint(member.MUID, pointID, member.DatabaseID, member.AppID); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				} else {
					if len(pointInfo.Points) == 0 { // 포인트가 존재하지 않는다면 에러
						resp.SetReturn(resultcode.Result_Error_MinPointQuantity)
						return ctx.EchoContext.JSON(http.StatusOK, resp)
					}
					if pointInfo.Points[0].PointID != pointID { // 포인트가 존재하지 않는다면 에러
						resp.SetReturn(resultcode.Result_Error_NotExistMember)
						return ctx.EchoContext.JSON(http.StatusOK, resp)
					}

					// 수량 업데이트
					if err := ProcGameChipSwap(params, pointID, pointInfo.MUID, pointInfo.DatabaseID, ctx.GetValue().InnoUID, pointInfo.Points[0].Quantity); err != nil {
						resp.SetReturn(resultcode.Result_DBError)
						return ctx.EchoContext.JSON(http.StatusOK, resp)
					}
				}
			}
		}
	default:
		resp.SetReturn(resultcode.Result_Error_Invalid_data)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func ProcGameChipSwap(params *context.ReqGameChipSwap, pointID int64, muid int64, databaseID int64, innoUID string, preSSRQ int64) error {
	// 레디스에 ssr 포인트가 존재 한다면 강제로 update 한다.
	preSSRQuantity := preSSRQ
	pointKey := model.MakeMemberPointListKey(muid)
	mePointInfo, err := model.GetDB().GetCacheMemberPointList(pointKey)
	if err == nil {
		// redis에 존재 한다면 강제로 db에 먼저 write
		for _, point := range mePointInfo.Points {
			var eventID context.EventID_type
			if point.AdjustQuantity >= 0 {
				eventID = context.EventID_add
			} else {
				eventID = context.EventID_sub
			}

			if point.AdjustQuantity != 0 {
				if _, _, err := model.GetDB().UpdateAppPoint(mePointInfo.DatabaseID, mePointInfo.MUID, point.PointID,
					point.PreQuantity, point.AdjustQuantity, point.Quantity, context.LogID_cp, eventID); err != nil {
					log.Errorf("UpdateAppPoint error : %v", err)
					return err
				} else {
					preSSRQuantity = point.Quantity
				}
			} else {
				preSSRQuantity = point.Quantity
			}
		}

		model.GetDB().DelCacheMemberPointList(pointKey)
	}

	// redis에도 없었다면 바로 업데이트 진행한다.
	if _, _, err := model.GetDB().UpdateAppPoint(databaseID, muid, pointID,
		preSSRQuantity, params.SSRAdjustPoint, preSSRQuantity+params.SSRAdjustPoint, context.LogID_external_inno_game, context.EventID_type(params.SwapType)); err != nil {
		log.Errorf("UpdateAppPoint error : %v", err)
		return err
	}

	// ssrm chip 업데이트 한다.
	if err := model.GetDB().USPG_Mod_Users_ChipQuantity(innoUID, params.SSRMAdjustChip, params.SwapType); err != nil {
		log.Errorf("USPG_Mod_Users_ChipQuantity : %v", err)
		return err
	}

	return nil
}

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
