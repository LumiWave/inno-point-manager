package commonapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/commonapi/inner"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
)

func PostWalletWebHookETHDeposit(params *context.ReqPostWalletETHResult, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := inner.TransferResultDepositWallet(params.FromAddr, params.ToAddr, params.Value, params.Symbol, params.TxHash, int64(params.TransactionFee), 18); err != nil {
		resp = err
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostWalletWebHookETHWithdrawal(params *context.ReqPostWalletETHResult, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	fee := strconv.FormatUint(params.TransactionFee, 10)

	if err := inner.TransferResultWithdrawalWallet(params.FromAddr, params.ToAddr, params.Value, fee,
		params.Symbol, params.TxHash, params.Status, 18); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostWalletWebHookSUIDeposit(params *context.SUI_CB_Balance_Changes, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	fromAddr := params.Sender
	toAddr := ""
	value := ""
	symbol := ""
	for _, balance := range params.Balances {
		if strings.EqualFold(balance.Owner, params.Sender) {
			continue // sender 본인의 변동(차감 토큰 + SUI 가스/리베이트)은 스킵
		}
		if !strings.Contains(balance.Amount, "-") {
			toAddr = balance.Owner
			value = balance.Amount
			symbol = balance.Symbol
		}
	}

	comput, _ := strconv.ParseInt(params.GasUsed.ComputationCost, 10, 64)
	cost, _ := strconv.ParseInt(params.GasUsed.StorageCost, 10, 64)
	rebate, _ := strconv.ParseInt(params.GasUsed.StorageRebate, 10, 64)

	gasFee := comput + cost - rebate

	if err := inner.TransferResultDepositWallet(fromAddr, toAddr, value, symbol, params.Digest, gasFee, 9); err != nil {
		resp = err
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostWalletWebHookSUIWithdrawal(params *context.SUI_CB_Balance_Changes, ctx *context.PointManagerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	fromAddr := params.Sender
	toAddr := ""
	value := ""
	symbol := ""
	for _, balance := range params.Balances {
		if strings.EqualFold(balance.Owner, params.Sender) {
			continue // sender 본인의 변동(원본 토큰 차감 + SUI 가스/리베이트)은 스킵
		}
		if !strings.Contains(balance.Amount, "-") {
			toAddr = balance.Owner
			value = balance.Amount
			symbol = balance.Symbol
		}
	}
	ComputationCost, _ := strconv.ParseInt(params.GasUsed.ComputationCost, 10, 64)
	StorageCost, _ := strconv.ParseInt(params.GasUsed.StorageCost, 10, 64)
	StorageRebate, _ := strconv.ParseInt(params.GasUsed.StorageRebate, 10, 64)
	fee := ComputationCost + StorageCost - StorageRebate

	fe := strconv.FormatInt(fee, 10)

	if err := inner.TransferResultWithdrawalWallet(fromAddr, toAddr, value, fe,
		symbol, params.Digest, params.Status, 9); err != nil {
		resp = err
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
