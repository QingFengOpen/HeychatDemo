package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"HyechatDemo/heybotclient/model"
)

func SendChannelIM(httpClient *http.Client, sendData *model.ChannelImSendReq, token string) error {
	url := model.HttpReqHost + "/chatroom/v2/channel_msg/send?chat_os_type=bot&chat_version=1.27.0&nonce=%v"
	url = fmt.Sprintf(url, time.Now().Unix())
	data, err := json.Marshal(sendData)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(data))

	req, err := http.NewRequest(http.MethodPost, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("token", token)
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			println(err)
		}
	}(res.Body)
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
