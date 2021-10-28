package models

import (
	"github.com/toolkits/pkg/logger"
)

type ExternalAlertResult struct {
	Id int64	`json:"id"`
	Event_id int64	`json:"eventId"`
	Channel string	`json:"channel"`
	Contacts string	`json:"contacts"`
	Result string	`json:"result"`
}

func (hae *ExternalAlertResult) Add() error {
	return DBInsertOne(hae)
}

/**
 *	@Desc 获取发送结果
 *	@Date 2021-10-28
 */
func GetExternalAlertResultsByIds(ids []int64) ([]ExternalAlertResult, error) {
	var externalAlertResult []ExternalAlertResult
	err := DB.In("event_id", ids).Find(&externalAlertResult)
	if err != nil {
		logger.Errorf("mysql.error: count external_alert_result fail: %v", err)
		return nil, internalServerError
	}
	return externalAlertResult, nil
}