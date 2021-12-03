package notify

import (
	"bytes"
	"encoding/json"
	"github.com/didi/nightingale/v5/config"
	model "github.com/didi/nightingale/v5/models"
	"github.com/toolkits/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const DingTalkTokenKey = "dingtalk_robot_token"
const WeComTokenKey = "wecom_robot_token"
const DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token="
const WeComUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

type Channel string

const (
	DingTalk = Channel("dingtalk")
	WeCom = Channel("weCom")
	SMS = Channel("sms")
	Voice = Channel("voice")
)

type MsgType string
const (
	Text = MsgType("text")
	Markdown = MsgType("markdown")
)

type DingTalkMessage struct {
	Msgtype MsgType	`json:"msgtype"`
	Markdown struct {
		Title string `json:"title,omitempty"`
		Text string	`json:"text,omitempty"`
	}	`json:"markdown,omitempty"`
	Text struct {
		Content string	`json:"content,omitempty"`
	}	`json:"text,omitempty"`
	At struct {
		AtMobiles []string	`json:"atMobiles"`
		IsAtAll bool	`json:"isAtAll"`
	}	`json:"at"`
}

type WeComMessage struct {
	Msgtype MsgType	`json:"msgtype"`
	Markdown struct {
		Content string `json:"content,omitempty"`
	}	`json:"markdown,omitempty"`
	Text struct {
		Content string	`json:"content,omitempty"`
		MentionedMobileList []string	`json:"mentioned_mobile_list,omitempty"`
	}	`json:"text"`
}

type VoiceMessage struct {
	Phones []string	`json:"phones"`
	Text	string	`json:"text"`
}

// 钉钉通知
func PostToDingTalk(text string, msgtype MsgType, users []*model.User, id int64) {
	tokenMap := make(map[string]DingTalkMessage)
	for _, user := range users {
		if user.Contacts == nil {
			continue
		}
		var contactKeys map[string]string
		if err := json.Unmarshal(user.Contacts, &contactKeys); err != nil {
			continue
		}
		atMobile := user.Phone
		// 判断用户是否设置dingtalk token
		if _, ok := contactKeys[DingTalkTokenKey]; ok {
			token := contactKeys[DingTalkTokenKey]
			if token == "" {
				continue
			}
			//	提取一致的token的消息群体
			if _, ok := tokenMap[token]; ok {
				atMobiles := tokenMap[token].At.AtMobiles
				atMobiles = append(atMobiles, atMobile)
			} else {
				// markdown格式指定@某人text里必须包含@信息
				var atMobiles []string
				atMobiles = append(atMobiles, atMobile)
				dingTalkMessage := DingTalkMessage{
					Msgtype: msgtype,
					At: struct {
						AtMobiles []string	`json:"atMobiles"`
						IsAtAll bool	`json:"isAtAll"`
					}{AtMobiles: atMobiles, IsAtAll: false},
				}
				switch msgtype {
				case Markdown: {
					dingTalkMessage.Markdown.Title = "钉钉通知"
					dingTalkMessage.Markdown.Text = text
				}
				case Text: {
					dingTalkMessage.Text.Content = text
				}
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
		Post(DingTalkUrl + token, postMessage, string(DingTalk), id, msg.At.AtMobiles)
	}
}

// 企业微信通知
func PostToWeCom(text string, msgtype MsgType, users []*model.User, id int64) {
	tokenMap := make(map[string]WeComMessage)
	for _, user := range users {
		if user.Contacts == nil {
			continue
		}
		var contactKeys map[string]string
		if err := json.Unmarshal(user.Contacts, &contactKeys); err != nil {
			continue
		}
		atMobile := user.Phone
		// 判断用户是否设置dingtalk token
		if _, ok := contactKeys[WeComTokenKey]; ok {
			token := contactKeys[WeComTokenKey]
			if token == "" {
				continue
			}
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
						Content string	`json:"content,omitempty"`
						MentionedMobileList []string	`json:"mentioned_mobile_list,omitempty"`
					}{Content: text, MentionedMobileList: atMobiles},
				}
				switch msgtype {
				case Markdown: {
					weComMessage.Markdown.Content = text
				}
				case Text: {
					weComMessage.Text.Content = text
				}
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
		Post(WeComUrl + token, postMessage, string(WeCom), id, msg.Text.MentionedMobileList)
	}
}

// 电话通知
func PostToVoice(text string, users []*model.User, id int64) {
	var phones []string
	for _, user := range users {
		phone := user.Phone
		if phone != "" {
			phones = append(phones, phone)
		}
	}
	voiceMessage := &VoiceMessage{
		phones,
		text,
	}
	postMessage, err := json.Marshal(voiceMessage)
	if err != nil {
		logger.Error(err.Error())
	}
	dialUrl := config.Config.AlertConfig.DialUrl
	Post(dialUrl, postMessage, string(Voice), id, phones)
}

func Post(url string, message []byte, logsign string, id int64, contacts []string) {
	reader := bytes.NewReader(message)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		logger.Errorf("【%s】消息发送失败：%s", logsign, err)
	}
	if resp == nil {
		logger.Error("resp is nil")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("【%s】连接失败：%s", logsign, err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("【%s】消息发送失败：%s", logsign, err)
	}
	// 发送记录
	sendResult := model.ExternalAlertResult{
		Event_id: id,
		Channel: string(DingTalk),
		Contacts: strings.Join(contacts, ","),
		Result: string(body),
	}
	_ = sendResult.Add()
	// 更新已发送
	_ = model.ExternalAlertEventUpdateStatus(id)
	logger.Infof("【%s】消息发送完成,服务器返回内容：%s", logsign, string(body))
}
