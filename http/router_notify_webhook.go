package http

import (
	model "github.com/didi/nightingale/v5/models"
	"github.com/didi/nightingale/v5/notify"
	"github.com/gin-gonic/gin"
	"strings"
)

type channel string

const (
	DingTalk = channel("DingTalk")
	WeCom = channel("WeCom")
	SMS = channel("SMS")
	Voice = channel("Voice")
)

type NotifyWebHookForm struct {
	Message string	`json:"message" validate:"required"`
	MsgType notify.MsgType	`json:"msgType" validate:"required"`
	Channel channel	`json:"channel" validate:"required"`
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
	users, err := model.UserPhoneGetByUsername(strings.Join(nwf.Contacts, ","))
	dangerous(err)
	switch nwf.Channel {
	case DingTalk:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users)
	case WeCom:
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users)
	case SMS:
	case Voice:
	}
	renderMessage(c, nil)
}