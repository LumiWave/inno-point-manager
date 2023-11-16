package inner

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/token_manager_server"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/model"
)

// auid에 해당하는 모든 지갑의 balance 정보 수집
func GetBalanceAll(auid int64) *base.BaseResponse {
	resp := new(base.BaseResponse)
	resp.Success()

	if wallets, _, err := model.GetDB().USPAU_GetList_AccountWallets(auid); err == nil {
		resBalanceAll := &context.ResBalanceAll{}
		resBalanceAll.Balances = []*context.Balance{}
		for _, wallet := range wallets {
			if wallet.ConnectionStatus == 1 {
				for _, coin := range model.GetDB().Coins {
					if wallet.BaseCoinID != coin.BaseCoinID || coin.CoinId >= 20000 { // nft는 제외
						continue
					}
					req := &token_manager_server.ReqBalance{
						BaseSymbol: model.GetDB().BaseCoinMapByCoinID[wallet.BaseCoinID].BaseCoinSymbol,
						Contract: func() string {
							// 코인 타입이면 contract 정보를 를 보내지 않는다.
							if strings.EqualFold(model.GetDB().BaseCoinMapByCoinID[wallet.BaseCoinID].BaseCoinSymbol, model.GetDB().Coins[coin.CoinId].CoinSymbol) {
								return ""
							}
							return model.GetDB().Coins[coin.CoinId].ContractAddress
						}(),
						Address: wallet.WalletAddress,
					}

					if res, err := token_manager_server.GetInstance().GetBalance(req); err != nil {
						resp.SetReturn(resultcode.ResultInternalServerError)
						return resp
					} else {
						if res.Return != 0 { // token manager 전송 에러
							log.Errorf("GetBalance return : %v, msg:%v", res.Return, res.Message)
						} else {
							log.Debugf("Balance symbol:%v, bal:%v, decimal:%v", coin.CoinSymbol, res.Balance, res.Decimal)
							resBalanceAll.Balances = append(resBalanceAll.Balances, &context.Balance{
								CoinID:     coin.CoinId,
								BaseCoinID: coin.BaseCoinID,
								Symbol:     coin.CoinSymbol,
								Balance:    res.Balance,
								Address:    wallet.WalletAddress,
								Decimal:    res.Decimal,
							})
						}
					}
				}
			}
		}
		resp.Value = resBalanceAll
	} else {
		resp.SetReturn(resultcode.Result_Error_Db_GetAccountWallets)
	}

	return resp
}
