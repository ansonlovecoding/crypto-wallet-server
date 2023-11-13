package push

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/common/push/jpush/common"
	"Share-Wallet/pkg/common/push/jpush/requestBody"
	"Share-Wallet/pkg/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	JPushClient *JPush
)

func init() {
	JPushClient = newGetuiClient()
}

type JPush struct{}

func newGetuiClient() *JPush {
	return &JPush{}
}

func (j *JPush) Auth(apiKey, secretKey string, timeStamp int64) (token string, err error) {
	return token, nil
}

func (j *JPush) SetAlias(cid, alias string) (resp string, err error) {
	return resp, nil
}

func (j *JPush) Push(accounts []string, alert, detailContent, operationID string) (string, error) {

	var pf requestBody.Platform
	pf.SetAll()
	var au requestBody.Audience
	au.SetAlias(accounts)
	var no requestBody.Notification

	no.IOSEnableMutableContent()
	no.SetAlert(alert)
	var extras requestBody.Extras
	extras.MsgContent = detailContent
	no.SetExtras(extras)
	var me requestBody.Message
	me.SetMsgContent(detailContent)
	var o requestBody.Options
	o.SetApnsProduction(config.Config.Push.Jpns.IsProduct)
	var po requestBody.PushObj
	po.SetPlatform(&pf)
	po.SetAudience(&au)
	po.SetNotification(&no)
	po.SetMessage(&me)
	po.SetOptions(&o)

	con, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.Push.Jpns.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", common.GetAuthorization(config.Config.Push.Jpns.AppKey, config.Config.Push.Jpns.MasterSecret))

	log.NewInfo(operationID, utils.GetSelfFuncName(), "push req:", req.Header, utils.Bytes2String(con))
	resp, err := client.Do(req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "client.Do failed:", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "ioutil.ReadAll failed:", err.Error())
		return "", err
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "push result:", string(result), "accounts", accounts)
	return string(result), nil
}
