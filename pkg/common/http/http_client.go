/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/27 10:31).
 */
package http

import (
	"Share-Wallet/pkg/common/constant"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// no share
type PostAPI struct {
	BaseAddress string
	// Todo: Implement hmac authentication
	Token string
}

func NewPostApi(token string, BaseURL string) *PostAPI {
	return &PostAPI{BaseAddress: BaseURL, Token: token}
}
func Get(url string) (response []byte, err error) {
	client := http.Client{Timeout: constant.APITimeout * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// application/json; charset=utf-8
func Post(url string, data interface{}, timeOutSecond int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: time.Duration(timeOutSecond) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func PostReturn(url string, input, output interface{}, timeOut int) error {
	b, err := Post(url, input, timeOut)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, output); err != nil {
		return err
	}
	return nil
}

// application/json; charset=utf-8
func (p *PostAPI) PostWalletAPI(url string, data interface{}, timeOutSecond int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// baseAddress := "http://cms.devwallet.com"
	// if url != constant.CreateAccountInformationURL {
	// 	baseAddress = "http://devwalletapi.ddns.net:81"
	// }
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", p.BaseAddress, url), bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: time.Duration(timeOutSecond) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
