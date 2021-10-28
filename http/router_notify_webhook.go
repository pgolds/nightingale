package http

import (
	model "github.com/didi/nightingale/v5/models"
	"github.com/didi/nightingale/v5/notify"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type channel string

const (
	DingTalk = channel("DingTalk")
	WeCom = channel("WeCom")
	SMS = channel("SMS")
	Voice = channel("Voice")
)

type level string

const (
	INFO = level("Info")
	WARNING = level("Warning")
	CRITICAL = level("Critical")
)

type NotifyWebHookForm struct {
	Message string	`json:"message" validate:"required"`
	MsgType notify.MsgType	`json:"msgType" validate:"required"`
	Channel channel	`json:"channel" validate:"required"`
	Level level	`json:"level" validate:"required"`
	Contacts []string	`json:"contacts" validate:"required"`
}

func NotifyWebHook(c *gin.Context) {
	var nwf NotifyWebHookForm
	bind(c, &nwf)
	if len(nwf.Contacts) == 0 {
		renderMessage(c, map[string]string{
			"err": "无联系人信息!",
		})
		return
	}
	eae := &model.ExternalAlertEvent{
		Level: string(nwf.Level),
		Msgtype: string(nwf.MsgType),
		Message: nwf.Message,
		Contacts: strings.Join(nwf.Contacts, ","),
		CreateAt: time.Now(),
	}
	_ = eae.Add()

	users, err := model.UserPhoneGetByUsername(nwf.Contacts)
	dangerous(err)
	switch nwf.Level {
	case CRITICAL:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
	case WARNING:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
	case INFO:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
	}
	renderMessage(c, nil)
}