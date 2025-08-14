package handler

import (
	"collectify/internal/config"
	"collectify/internal/dao"
	"collectify/internal/db"
	define "collectify/internal/model/define"
	model "collectify/internal/model/db"
	"collectify/internal/pkg/e"
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateField(c *gin.Context) {
	var req define.CreateFieldReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{
		"category_id": req.CategoryID,
		"name":        req.Name,
	}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.Field](db.GetDB(), uniqueFields, filters)
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

	field := &model.Field{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Type:       req.Type,
		IsArray:    req.IsArray,
		Required:   req.Required,
	}

	err = dao.Create(db.GetDB(), field)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func DeleteField(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	if config.GetConfig().RecycleBin.Enable {
		err = service.SoftDelete(model.ModelTypeField, map[string]interface{}{"id": id})
		if err != nil {
			Fail(c, err)
			return
		}
	} else {
		err = service.HardDelete(model.ModelTypeField, map[string]interface{}{"id": id})
		if err != nil {
			Fail(c, err)
			return
		}
	}

	Success(c)
}

