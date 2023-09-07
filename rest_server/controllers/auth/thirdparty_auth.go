package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
)

var gTransport http.Transport
var gClient http.Client

func InitHttpClient() {
	gTransport.MaxIdleConns = 100
	gTransport.MaxIdleConnsPerHost = 100
	gTransport.IdleConnTimeout = 30 * time.Second
	gTransport.DisableKeepAlives = false

	gClient.Timeout = 60 * time.Second
	gClient.Transport = &gTransport
}

// verify auth token respoinse
type VerifyAuthToken struct {
	CompanyID int64  `json:"company_id"`
	AppID     int64  `json:"app_id"`
	LoginType int64  `json:"login_type"`
	Uuid      string `json:"uuid"`
	InnoUID   string `json:"inno_uid"`
	AUID      int64  `json:"au_id"`
}

type AuthResponse struct {
	Return  int             `json:"return"`
	Message string          `json:"message"`
	Value   VerifyAuthToken `json:"value"`
}

// ///////////////////////
func CheckAuthToken(authToken string) (*VerifyAuthToken, *AuthResponse, error) {
	conf := config.GetInstance()

	callURL := fmt.Sprintf("%s%s", conf.Auth.ApiAuthDomain, conf.Auth.ApiAuthVerify)

	req, err := http.NewRequest("GET", callURL, bytes.NewBuffer(nil))
	if err != nil {
		log.Errorf("%v", err)
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp, err := gClient.Do(req)

	if err != nil {
		log.Errorf("membership resp: %v, err: %v", resp, err)
		return nil, nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Errorf("resp: %v, body close err: %v", resp, err)
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	baseResp := new(AuthResponse)
	err = decoder.Decode(baseResp)
	if err != nil {
		log.Errorf("resp: %v, docode err: %v", resp, err)
		return nil, nil, err
	}

	if baseResp.Message != "success" {
		//err := errors.New(baseResp.Message)
		//log.Errorf("resp: %v, body close err: %v", resp, err)
		return nil, baseResp, nil
	}

	return &baseResp.Value, baseResp, nil
}
