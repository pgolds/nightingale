package http

import (
	"github.com/didi/nightingale/v5/models"
	"github.com/gin-gonic/gin"
)

func externalAlertEventGets(c *gin.Context) {
	channel := queryStr(c, "channel", "")
	hasSend := queryBool(c, "hasSend", true)
	limit := queryInt(c, "limit", defaultLimit)

	total, err := models.ExternalAlertEventTotal(channel, hasSend)
	dangerous(err)
	list, err := models.ExternalAlertEventGets(channel, hasSend, limit, offset(c, limit))
	dangerous(err)

	if len(list) == 0 {
		renderZeroPage(c)
		return
	}

	renderData(c, map[string]interface{}{
		"total": total,
		"list":  list,
	}, nil)
}
