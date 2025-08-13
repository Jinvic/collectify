package handler

import (
	"collectify/internal/config"
	"collectify/internal/dao"
	"collectify/internal/db"
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	"collectify/internal/pkg/e"
	"collectify/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func CreateCategory(c *gin.Context) {
	var req CreateCategoryReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	uniqueFields := map[string]interface{}{"name": req.Name}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.Category](db.GetDB(), uniqueFields, filters)
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

	if config.GetConfig().RecycleBin.Enable {
		err = service.SoftDelete(model.ModelTypeCategory, map[string]interface{}{"id": id})
		if err != nil {
			Fail(c, err)
			return
		}
	} else {
		err = service.HardDelete(model.ModelTypeCategory, map[string]interface{}{"id": id})
		if err != nil {
			Fail(c, err)
			return
		}
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
	uniqueFields := map[string]interface{}{
		"name": req.Name,
	}
	filters := []dao.Filter{
		{
			Where: "id != ?",
			Args:  []interface{}{id},
		},
	}
	id, isDeleted, err := dao.DuplicateCheck[model.Category](db.GetDB(), uniqueFields, filters)
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

	uniqueFields = map[string]interface{}{"id": id}
	updateFields := map[string]interface{}{"name": req.Name}
	err = dao.Update[model.Category](db.GetDB(), uniqueFields, updateFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func GetCategory(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	category, err := dao.Get[model.Category](db.GetDB(), map[string]interface{}{"id": id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, e.ErrNotFound)
			return
		}
		Fail(c, err)
		return
	}

	SuccessWithData(c, category)
}

func SearchCategory(c *gin.Context) {
	var req SearchReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	pagination := common.Pagination{
		Disable: req.NoPaging,
		Page:    req.Page,
		Size:    req.PageSize,
	}

	var filters []dao.Filter

	if name, ok := req.Filters["name"]; ok {
		filters = append(filters, dao.Filter{
			Where: "name LIKE ?",
			Args:  []interface{}{"%" + cast.ToString(name) + "%"},
		})
	}

	categories, total, err := dao.GetList[model.Category](db.GetDB(), filters, pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	SuccessWithData(c, SearchResp{
		List:  categories,
		Total: total,
	})
}
