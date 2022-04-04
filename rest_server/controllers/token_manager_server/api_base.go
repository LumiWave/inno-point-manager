package token_manager_server

import (
	"bytes"
	"crypto/tls"
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
	Api_post_sendfrom_userWallet   = 1 // 특정지갑에서 출금 PostSendFromUserWallet

	Api_get_coin_fee = 3 // 코인 가스비 조회
)

type ApiInfo struct {
	ApiType          api_kind
	Method           string
	Uri              string
	ResponseType     interface{}
	ResponseFuncType func() interface{}
	client           *http.Client
}

var ApiList = map[api_kind]ApiInfo{
	Api_post_sendfrom_parentwallet: ApiInfo{ApiType: Api_post_sendfrom_parentwallet, Method: "POST", Uri: "/token/transfer", ResponseType: new(ResSendFromParentWallet),
		ResponseFuncType: func() interface{} { return new(ResSendFromParentWallet) }, client: NewClient()},
	Api_post_sendfrom_userWallet: ApiInfo{ApiType: Api_post_sendfrom_userWallet, Method: "POST", Uri: "/token/transfer/user", ResponseType: new(ResSendFromUserWallet),
		ResponseFuncType: func() interface{} { return new(ResSendFromUserWallet) }, client: NewClient()},

	Api_get_coin_fee: ApiInfo{ApiType: Api_post_sendfrom_userWallet, Method: "GET", Uri: "/token/coin/fee", ResponseType: new(ResCoinFeeInfo),
		ResponseFuncType: func() interface{} { return new(ResCoinFeeInfo) }, client: NewClient()},
}

func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100
	t.IdleConnTimeout = 30 * time.Second
	t.DisableKeepAlives = false
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	return client
}

func MakeHttp(callUrl string, auth string, method string, body *bytes.Buffer, queryStr string) *http.Request {
	req, err := http.NewRequest(method, callUrl, body)
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(auth) > 0 {
		req.Header.Add("Authorization", "Bearer "+auth)
	}
	if len(queryStr) > 0 {
		req.URL.RawQuery = queryStr
	}

	return req
}

func HttpCall(client *http.Client, callUrl string, auth string, method string, kind api_kind, body *bytes.Buffer, queryStruct interface{}, response interface{}) (interface{}, error) {
	var v url.Values
	var queryStr string
	if queryStruct != nil {
		v, _ = query.Values(queryStruct)
		queryStr = v.Encode()
	}

	req := MakeHttp(callUrl, auth, method, body, queryStr)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp, kind, response)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseResponse(resp *http.Response, kind api_kind, response interface{}) (interface{}, error) {
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(response)
	if err != nil {
		return nil, err
	}
	return response, err
}
