package handler

import (
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func Restore(c *gin.Context) {
	var req DeletedReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{
		"id": req.ID,
	}

	err := service.Restore(req.Type, uniqueFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func HardDelete(c *gin.Context) {
	var req DeletedReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{
		"id": req.ID,
	}

	err := service.HardDelete(req.Type, uniqueFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}
