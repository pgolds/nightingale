package models

import (
	"github.com/toolkits/pkg/logger"
	"time"
	"xorm.io/builder"
)

type ExternalAlertEvent struct {
	Id       int64  `json:"id"`
	Channel	 string	`json:"channel"`
	Msgtype  string	`json:"msgtype"`
	Message  string `json:"message"`
	Contacts string `json:"contacts"`
	Result   string `json:"result"`
	CreateAt time.Time `json:"createAt"`
	HasSend  bool	`json:"hasSend"`
}

func (hae *ExternalAlertEvent) Add() error {
	return DBInsertOne(hae)
}

/**
 *  @Desc 获取外部系统调用webhook记录列表
 *	@Date 2021-10-20
 */
func ExternalAlertEventGets(channel string, hasSend bool, limit, offset int) ([]ExternalAlertEvent, error) {
	cond := builder.NewCond()
	cond = cond.And(builder.Eq{"has_send": hasSend})
	if channel != "" {
		cond = cond.And(builder.Eq{"channel": channel})
	}
	var objs []ExternalAlertEvent
	err := DB.Where(cond).Desc("create_at").Limit(limit, offset).Find(&objs)
	if err != nil {
		logger.Errorf("mysql.error: gets external_alert_event fail: %v", err)
		return nil, internalServerError
	}
	return objs, nil
}