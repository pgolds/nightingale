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
 *	@Desc 获取外部系统调用webhook记录总数
 *	@Date 2021-10-25
 */
func ExternalAlertEventTotal(channel string, hasSend bool) (total int64, err error) {
	cond := builder.NewCond()
	cond = cond.And(builder.Eq{"has_send": hasSend})
	if channel != "" {
		cond = cond.And(builder.Eq{"channel": channel})
	}
	num, err := DB.Where(cond).Count(new(ExternalAlertEvent))
	if err != nil {
		logger.Errorf("mysql.error: count external_alert_event fail: %v", err)
		return 0, internalServerError
	}
	return num, nil
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

/**
 *	@Desc 更新通知发送返回结果
 *	@Date 2021-10-20
 */
func ExternalAlertEventUpdateResult(id int64, result string) error {
	_, err := DB.Exec("UPDATE external_alert_event set result=?, has_send=1 where id =?", result, id)
	if err != nil {
		logger.Errorf("mysql.error: update external_alert_event result fail: %s", err)
		return internalServerError
	}
	return nil
}