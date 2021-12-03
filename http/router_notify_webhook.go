package http

import (
	model "github.com/didi/nightingale/v5/models"
	"github.com/didi/nightingale/v5/notify"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type level string

const (
	P1 = level("P1")
	P2 = level("P2")
	P3 = level("P3")
)

type NotifyWebHookForm struct {
	Message string	`json:"message" binding:"required"`
	MsgType notify.MsgType	`json:"msgType" binding:"required"`
	Channel notify.Channel	`json:"channel"`
	Level level	`json:"level" binding:"required"`
	Contacts []string	`json:"contacts" binding:"required"`
}

func NotifyWebHook(c *gin.Context) {
	var nwf NotifyWebHookForm
	bind(c, &nwf)
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
	case P3:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
	case P2:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
	case P1:
		go notify.PostToDingTalk(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToWeCom(nwf.Message, nwf.MsgType, users, eae.Id)
		go notify.PostToVoice(nwf.Message, users, eae.Id)
	}
	renderMessage(c, nil)
}