package handler

import (
	define "collectify/internal/model/define"
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func Restore(c *gin.Context) {
	var req define.DeletedReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	for _, item := range req.List {
		uniqueFields := map[string]interface{}{
			"id": item.ID,
		}
		err := service.Restore(item.Type, uniqueFields)
		if err != nil {
			Fail(c, err)
			return
		}
	}

	Success(c)
}

func HardDelete(c *gin.Context) {
	var req define.DeletedReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	for _, item := range req.List {
		uniqueFields := map[string]interface{}{
			"id": item.ID,
		}
		err := service.HardDelete(item.Type, uniqueFields)
		if err != nil {
			Fail(c, err)
			return
		}
	}

	Success(c)
}
