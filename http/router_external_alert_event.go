package http

import (
	"github.com/didi/nightingale/v5/models"
	"github.com/gin-gonic/gin"
)

func externalAlertEventGets(c *gin.Context) {
	level := queryStr(c, "level", "")
	hasSend := queryBool(c, "hasSend", true)
	limit := queryInt(c, "limit", defaultLimit)

	total, err := models.ExternalAlertEventTotal(level, hasSend)
	dangerous(err)
	list, err := models.ExternalAlertEventGets(level, hasSend, limit, offset(c, limit))
	dangerous(err)

	if len(list) == 0 {
		renderZeroPage(c)
		return
	}

	var ids []int64
	for _, event := range list {
		ids = append(ids, event.Id)
	}
	results, err := models.GetExternalAlertResultsByIds(ids)
	dangerous(err)
	for i, event := range list {
		var res []models.ExternalAlertResult
		for _, result := range results {
			if result.Event_id == event.Id {
				res = append(res, result)
			}
		}
		event.Result = res
		list[i] = event
	}
	renderData(c, map[string]interface{}{
		"total": total,
		"list":  list,
	}, nil)
}
