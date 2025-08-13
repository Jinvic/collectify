package handler

import (
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"
	"collectify/internal/pkg/e"
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var req CreateCategoryReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	id, isDeleted, err := dao.DuplicateCheck[model.Category](db.GetDB(), map[string]interface{}{"name": req.Name})
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

	// 创建
	category := &model.Category{
		Name: req.Name,
	}
	_, err = dao.Create(db.GetDB(), category)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func DeleteCategory(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	err = service.SoftDelete(model.ModelTypeCategory, map[string]interface{}{"id": id})
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RenameCategory(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	var req RenameCategoryReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	id, isDeleted, err := dao.DuplicateCheck[model.Category](db.GetDB(), map[string]interface{}{
		"name":    req.Name,
		"id != ?": id,
	})
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

	uniqueFields := map[string]interface{}{"id": id}
	updateFields := map[string]interface{}{"name": req.Name}
	err = dao.Update[model.Category](db.GetDB(), uniqueFields, updateFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}
