package token_manager_server

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (o *TokenManagerServerInfo) PostSendFromParentWallet(req *ReqSendFromParentWallet) (*ResSendFromParentWallet, error) {
	callUrl := fmt.Sprintf("%s%s%s", o.IntHostUri, o.IntVer, ApiList[Api_post_sendfrom_parentwallet].Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	data, err := HttpCall(callUrl, o.ApiKey, "POST", Api_post_sendfrom_parentwallet, buff, nil, &ResSendFromParentWallet{})
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

	data, err := HttpCall(callUrl, o.ApiKey, "POST", urlInfo.ApiType, buff, nil, urlInfo.ResponseFuncType())
	if err != nil {
		return nil, err
	}

	return data.(*ResSendFromUserWallet), nil
}
