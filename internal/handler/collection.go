package handler

import (
	"collectify/internal/config"
	"collectify/internal/conn"
	"collectify/internal/dao"
	"collectify/internal/model/common"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"collectify/internal/pkg/e"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCollection(c *gin.Context) {
	var req define.CreateCollectionReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	uniqueFields := map[string]interface{}{"name": req.Name}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.Collection](conn.GetDB(), uniqueFields, filters)
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
	collection := &model.Collection{
		Name: req.Name,
	}
	err = dao.Create(conn.GetDB(), collection)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func DeleteCollection(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	isSoftDelete := config.GetConfig().RecycleBin.Enable
	uniqueFields := map[string]interface{}{"id": id}
	err = dao.Delete[model.Collection](conn.GetDB(), uniqueFields, isSoftDelete)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RestoreCollection(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	err = dao.Restore[model.Collection](conn.GetDB(), uniqueFields)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c)
}

func UpdateCollection(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	var req define.UpdateCollectionReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	collection, err := dao.Get[model.Collection](conn.GetDB(), uniqueFields)
	if err != nil {
		Fail(c, err)
		return
	}

	// 如果名称不同，则检查是否重复
	if req.Name != "" && collection.Name != req.Name {
		uniqueFields := map[string]interface{}{
			"name": req.Name,
		}
		filters := []dao.Filter{
			{
				Where: "id != ?",
				Args:  []interface{}{id},
			},
		}
		id, isDeleted, err := dao.DuplicateCheck[model.Tag](conn.GetDB(), uniqueFields, filters)
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
	}

	updateFields := map[string]interface{}{}
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}

	if len(updateFields) == 0 {
		Success(c)
		return
	}

	err = dao.Update[model.Collection](conn.GetDB(), uniqueFields, updateFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func GetCollection(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	collection, err := dao.Get[model.Collection](conn.GetDB(), uniqueFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, e.ErrNotFound)
			return
		}
		Fail(c, err)
		return
	}

	collectionInfo := define.Collection{}
	collectionInfo.FromDB(&collection)

	SuccessWithData(c, collectionInfo)
}

func ListCollection(c *gin.Context) {
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
	collections, total, err := dao.GetList[model.Collection](conn.GetDB(), filters, orderBy, pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	collectionInfos := make([]define.Collection, len(collections))
	for idx, collection := range collections {
		collectionInfos[idx].FromDB(&collection)
	}

	SuccessWithData(c, define.SearchResp{
		List:  collectionInfos,
		Total: total,
	})
}
