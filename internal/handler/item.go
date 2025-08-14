package handler

import (
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"collectify/internal/pkg/e"
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {
	var req define.CreateItemReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{
		"category_id": req.CategoryID,
		"title":       req.Title,
	}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.Item](db.GetDB(), uniqueFields, filters)
	if err != nil {
		Fail(c, err)
		return
	}
	if id != 0 {
		FailWithData(c, e.ErrDuplicated, map[string]interface{}{
			"id":        id,
			"isDeleted": isDeleted,
		})
		return
	}

	err = service.CreateItem(req.CategoryID, req.Title, req.Values)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}
