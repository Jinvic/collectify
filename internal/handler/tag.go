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

func CreateTag(c *gin.Context) {
	var req define.CreateTagReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	// 检查是否重复
	uniqueFields := map[string]interface{}{"name": req.Name}
	filters := []dao.Filter{}
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

	// 创建
	tag := &model.Tag{
		Name: req.Name,
	}
	err = dao.Create(conn.GetDB(), tag)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func DeleteTag(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	isSoftDelete := config.GetConfig().RecycleBin.Enable
	uniqueFields := map[string]interface{}{"id": id}
	err = dao.Delete[model.Tag](conn.GetDB(), uniqueFields, isSoftDelete)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RestoreTag(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	err = dao.Restore[model.Tag](conn.GetDB(), uniqueFields)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c)
}

func RenameTag(c *gin.Context) {
	tagID, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	var req define.RenameTagReq
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
			Args:  []interface{}{tagID},
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

	uniqueFields = map[string]interface{}{"id": tagID}
	updateFields := map[string]interface{}{"name": req.Name}
	err = dao.Update[model.Tag](conn.GetDB(), uniqueFields, updateFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func GetTag(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	tag, err := dao.Get[model.Tag](conn.GetDB(), uniqueFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, e.ErrNotFound)
			return
		}
		Fail(c, err)
		return
	}

	tagInfo := define.Tag{}
	tagInfo.FromDB(&tag)

	SuccessWithData(c, tagInfo)
}

func ListTag(c *gin.Context) {
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
	tags, total, err := dao.GetList[model.Tag](conn.GetDB(), filters, orderBy, pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	tagInfos := make([]define.Tag, len(tags))
	for idx, tag := range tags {
		tagInfos[idx].FromDB(&tag)
	}

	SuccessWithData(c, define.SearchResp{
		List:  tagInfos,
		Total: total,
	})
}
