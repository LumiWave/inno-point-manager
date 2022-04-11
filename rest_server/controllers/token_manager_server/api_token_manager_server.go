package token_manager_server

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (o *TokenManagerServerInfo) PostSendFromParentWallet(req *ReqSendFromParentWallet) (*ResSendFromParentWallet, error) {
	urlInfo := ApiList[Api_post_sendfrom_parentwallet]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, urlInfo.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(urlInfo.client, callUrl, o.ApiKey, urlInfo.Method, urlInfo.ApiType, buff, nil, urlInfo.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*ResSendFromParentWallet), nil
}

func (o *TokenManagerServerInfo) PostSendFromUserWallet(req *ReqSendFromUserWallet) (*ResSendFromUserWallet, error) {
	urlInfo := ApiList[Api_post_sendfrom_userWallet]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, urlInfo.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(urlInfo.client, callUrl, o.ApiKey, urlInfo.Method, urlInfo.ApiType, buff, nil, urlInfo.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*ResSendFromUserWallet), nil
}

func (o *TokenManagerServerInfo) GetBalance(req *ReqBalance) (*ResBalanc, error) {
	urlInfo := ApiList[Api_get_balance]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, urlInfo.Uri)

	data, err := HttpCall(urlInfo.client, callUrl, o.ApiKey, urlInfo.Method, urlInfo.ApiType, bytes.NewBuffer(nil), req, urlInfo.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*ResBalanc), nil
}

func (o *TokenManagerServerInfo) GetCoinFee(req *ReqCoinFee) (*ResCoinFeeInfo, error) {
	urlInfo := ApiList[Api_get_coin_fee]
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, urlInfo.Uri)

	data, err := HttpCall(urlInfo.client, callUrl, o.ApiKey, urlInfo.Method, urlInfo.ApiType, bytes.NewBuffer(nil), req, urlInfo.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*ResCoinFeeInfo), nil
}
