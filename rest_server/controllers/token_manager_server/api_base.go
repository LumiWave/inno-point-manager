package token_manager_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

type api_kind int

const (
	Api_post_sendfrom_parentwallet = 0 // 대표지갑에서 출금 PostSendFromParentWallet
)

type ApiInfo struct {
	ApiType      api_kind
	Uri          string
	ResponseType interface{}
}

var ApiList = map[api_kind]ApiInfo{
	Api_post_sendfrom_parentwallet: ApiInfo{ApiType: Api_post_sendfrom_parentwallet, Uri: "/token/transfer", ResponseType: new(ResSendFromParentWallet)},
}

func MakeHttpClient(callUrl string, auth string, method string, body *bytes.Buffer, queryStr string) (*http.Client, *http.Request) {
	req, err := http.NewRequest(method, callUrl, body)
	if err != nil {
		return nil, nil
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(auth) > 0 {
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	if len(queryStr) > 0 {
		req.URL.RawQuery = queryStr
	}

	client := &http.Client{Timeout: 5 * time.Second}
	return client, req
}

func HttpCall(callUrl string, auth string, method string, kind api_kind, body *bytes.Buffer, queryStruct interface{}) (interface{}, error) {

	var v url.Values
	var queryStr string
	if queryStruct != nil {
		v, _ = query.Values(queryStruct)
		queryStr = v.Encode()
	}

	client, req := MakeHttpClient(callUrl, auth, method, body, queryStr)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp, kind)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseResponse(resp *http.Response, kind api_kind) (interface{}, error) {
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)

	strc := ApiList[kind]
	err := decoder.Decode(strc.ResponseType)
	if err != nil {
		return nil, err
	}
	return strc.ResponseType, err
}