package model

import (
	"encoding/json"
	"errors"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-point-manager/rest_server/config"
	"github.com/LumiWave/inno-point-manager/rest_server/controllers/context"
)

func MakeSwapWalletKey() string {
	return config.GetInstance().DBPrefix + ":SWAP-WALLET"
}

// swap 정보 저장
func (o *DB) CacheSetSwapWallets(swapInfos map[string]interface{}) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.HMSet(MakeSwapWalletKey(), swapInfos)
}

func (o *DB) CacheSetSwapWallet(swapInfo *context.ReqSwapInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	swapInfos := make(map[string]interface{})

	if swapInfo.TxType == context.EventID_P2C {
		swapInfos[swapInfo.SwapToCoin.WalletAddress] = swapInfo
	} else if swapInfo.TxType == context.EventID_C2P {
		swapInfos[swapInfo.SwapFromCoin.WalletAddress] = swapInfo
	} else if swapInfo.TxType == context.EventID_C2C { // C2C는 from(실제코인 전송 콜백 확인용), to(수수료 콜백 확인용) 둘다 남긴다.
		swapInfos[swapInfo.SwapToCoin.WalletAddress] = swapInfo
		swapInfos[swapInfo.SwapFromCoin.WalletAddress] = swapInfo
	}

	return o.Cache.HMSet(MakeSwapWalletKey(), swapInfos)
}

// 전체 swap 정보 조회
func (o *DB) CacheGetSwapWallets() (map[string]*context.ReqSwapInfo, []*context.ReqSwapInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	rMap, err := o.Cache.GetDB().HGetAll(MakeSwapWalletKey())
	if err != nil {
		return nil, nil, err
	}

	resMap := make(map[string]*context.ReqSwapInfo)
	resList := []*context.ReqSwapInfo{}
	for _, value := range rMap {
		loadData := &context.ReqSwapInfo{}
		if err := json.Unmarshal([]byte(value), loadData); err != nil {
			log.Errorf("CacheGetSwapWallets unmarshal err : %v", err)
			loadData = nil
		} else {
			if loadData.TxType == context.EventID_P2C {
				resMap[loadData.SwapToCoin.WalletAddress] = loadData
			} else if loadData.TxType == context.EventID_C2P || loadData.TxType == context.EventID_C2C {
				resMap[loadData.SwapFromCoin.WalletAddress] = loadData
			}

			//resMap[loadData.WalletAddress] = loadData
			resList = append(resList, loadData)
		}
	}

	return resMap, resList, nil
}

// 단일 swap 정보 조회
func (o *DB) CacheGetSwapWallet(walletAddress string) (*context.ReqSwapInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	retList, err := o.Cache.GetDB().HMGet(MakeSwapWalletKey(), walletAddress) // 주의 : 값이 없으면 retList[0]와 err에 모두 nil이 응답된다.
	if err != nil {
		return nil, err
	}
	if retList[0] == nil {
		return nil, errors.New("not exist")
	}
	loadData := &context.ReqSwapInfo{}
	if err := json.Unmarshal([]byte(retList[0].(string)), loadData); err != nil {
		log.Errorf("CacheGetSwapWallet unmarshal err : %v", err)
		loadData = nil
	}

	return loadData, err
}

// 단일 swap 정보 삭제
func (o *DB) CacheDelSwapWallet(walletAddress string) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.HDel(MakeSwapWalletKey(), walletAddress)
}

func (o *DB) CacheDelAllSwapWallet() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	return o.Cache.GetDB().HDelAll(MakeSwapWalletKey())
}
