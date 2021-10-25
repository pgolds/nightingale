package notify

import (
	"bytes"
	"encoding/json"
	model "github.com/didi/nightingale/v5/models"
	"github.com/toolkits/pkg/logger"
	"io/ioutil"
	"net/http"
)

const DingTalkTokenKey = "dingtalk_robot_token"
const WeComTokenKey = "wecom_robot_token"
const DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token="
const WeComUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

type MsgType string
const (
	Text = MsgType("text")
	Markdown = MsgType("markdown")
)

type DingTalkMessage struct {
	Msgtype MsgType	`json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text string	`json:"text"`
	}	`json:"markdown"`
	Text struct {
		Content string	`json:"content"`
	}	`json:"text"`
	At struct {
		AtMobiles []string	`json:"atMobiles"`
		IsAtAll bool	`json:"IsAtAll"`
	}	`json:"at"`
}

type WeComMessage struct {
	Msgtype MsgType	`json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	}	`json:"markdown"`
	Text struct {
		Content string	`json:"content"`
		MentionedMobileList []string	`json:"mentioned_mobile_list"`
	}	`json:"text"`
}

func PostToDingTalk(text string, msgtype MsgType, users []model.User, id int64) {
	tokenMap := make(map[string]DingTalkMessage)
	for _, user := range users {
		if user.Contacts == nil {
			continue
		}
		logger.Infof("user.Contacts: %s", user.Contacts)
		var contactKeys map[string]string
		if err := json.Unmarshal(user.Contacts, &contactKeys); err != nil {
			continue
		}
		atMobile := "@" + user.Phone
		// 判断用户是否设置dingtalk token
		if _, ok := contactKeys[DingTalkTokenKey]; ok {
			token := contactKeys[DingTalkTokenKey]
			//	提取一致的token的消息群体
			if _, ok := tokenMap[token]; ok {
				atMobiles := tokenMap[token].At.AtMobiles
				atMobiles = append(atMobiles, atMobile)
			} else {
				var atMobiles []string
				atMobiles = append(atMobiles, atMobile)
				dingTalkMessage := DingTalkMessage{
					Msgtype: msgtype,
					Text: struct {
						Content string	`json:"content"`
					}{Content: text},
					At: struct {
						AtMobiles []string	`json:"atMobiles"`
						IsAtAll bool	`json:"IsAtAll"`
					}{AtMobiles: atMobiles, IsAtAll: false},
				}
				tokenMap[token] = dingTalkMessage
			}
		}
	}

	for token, msg := range tokenMap {
		postMessage, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err.Error())
		}
		Post(DingTalkUrl + token, postMessage, "DingTalk", id)
	}
}

func PostToWeCom(text string, msgtype MsgType, users []model.User, id int64) {
	tokenMap := make(map[string]WeComMessage)
	for _, user := range users {
		if user.Contacts == nil {
			continue
		}
		logger.Infof("user.Contacts: %s", user.Contacts)
		var contactKeys map[string]string
		if err := json.Unmarshal(user.Contacts, &contactKeys); err != nil {
			continue
		}
		atMobile := user.Phone
		// 判断用户是否设置dingtalk token
		if _, ok := contactKeys[WeComTokenKey]; ok {
			token := contactKeys[WeComTokenKey]
			//	提取一致的token的消息群体
			if _, ok := tokenMap[token]; ok {
				atMobiles := tokenMap[token].Text.MentionedMobileList
				atMobiles = append(atMobiles, atMobile)
			} else {
				var atMobiles []string
				atMobiles = append(atMobiles, atMobile)
				weComMessage := WeComMessage{
					Msgtype: msgtype,
					Text: struct {
						Content string	`json:"content"`
						MentionedMobileList []string	`json:"mentioned_mobile_list"`
					}{Content: text, MentionedMobileList: atMobiles},
				}
				tokenMap[token] = weComMessage
			}
		}
	}

	for token, msg := range tokenMap {
		postMessage, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err.Error())
		}
		Post(WeComUrl + token, postMessage, "WeCom", id)
	}
}

func Post(url string, message []byte, logsign string, id int64) {
	reader := bytes.NewReader(message)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		logger.Errorf("【%s】消息发送失败：%s", logsign, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("【%s】消息发送失败：%s", logsign, err)
	}
	_ = model.ExternalAlertEventUpdateResult(id, string(body))
	logger.Infof("【%s】消息发送完成,服务器返回内容：%s", logsign, string(body))
}
