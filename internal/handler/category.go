package handler

import (
	"collectify/internal/conn"
	"collectify/internal/dao"
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"collectify/internal/pkg/e"
	"collectify/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(c *gin.Context) {
	var req define.CreateCategoryReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	uniqueFields := map[string]interface{}{"name": req.Name}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.Category](conn.GetDB(), uniqueFields, filters)
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
	err = dao.Create(conn.GetDB(), category)
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

	err = service.DeleteCategory(id)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RestoreCategory(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	err = service.RestoreCategory(id)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RenameCategory(c *gin.Context) {
	categoryID, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	var req define.RenameCategoryReq
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
			Args:  []interface{}{categoryID},
		},
	}
	id, isDeleted, err := dao.DuplicateCheck[model.Category](conn.GetDB(), uniqueFields, filters)
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

	uniqueFields = map[string]interface{}{"id": categoryID}
	updateFields := map[string]interface{}{"name": req.Name}
	err = dao.Update[model.Category](conn.GetDB(), uniqueFields, updateFields)
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

	uniqueFields := map[string]interface{}{"id": id}
	preloads := []string{"Fields"}
	category, err := dao.Get[model.Category](conn.GetDB(), uniqueFields, preloads...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, e.ErrNotFound)
			return
		}
		Fail(c, err)
		return
	}

	categoryInfo := define.Category{}
	categoryInfo.FromDB(&category)

	SuccessWithData(c, categoryInfo)
}

func ListCategory(c *gin.Context) {
	name := c.Query("name")

	var filters []dao.Filter
	var orderBy []dao.OrderBy

	if name != "" {
		filters = append(filters, dao.Filter{
			Where: "name LIKE ?",
			Args:  []interface{}{"%" + name + "%"},
		})
	}

	// 不分页
	pagination := common.Pagination{
		Disable: true,
	}

	// 创建时间顺序排序
	orderBy = []dao.OrderBy{
		{
			Column: "created_at",
			Desc:   false,
		},
	}
	categories, total, err := dao.GetList[model.Category](conn.GetDB(), filters, orderBy, pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	categoryInfos := make([]define.Category, len(categories))
	for idx, category := range categories {
		categoryInfos[idx].FromDB(&category)
	}

	SuccessWithData(c, define.SearchResp{
		List:  categoryInfos,
		Total: total,
	})
}
