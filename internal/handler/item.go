package handler

import (
	"collectify/internal/conn"
	"collectify/internal/dao"
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"collectify/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {
	var req define.CreateItemReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	item := req.Item.ToDB()
	item.CategoryID = req.CategoryID

	err := service.CreateItem(item, req.Item.Values)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func UpdateItem(c *gin.Context) {
	var req define.UpdateItemReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	item := req.Item.ToDB()
	item.ID = req.ID

	err := service.UpdateItem(item, req.Item.Values)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func DeleteItem(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	err = service.DeleteItem(id)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func RestoreItem(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	err = service.RestoreItem(id)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}

func ListItems(c *gin.Context) {
	pagination, err := GetPagination(c)
	if err != nil {
		Fail(c, err)
		return
	}

	items, total, err := service.ListItems(pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	itemInfos := make([]define.Item, len(items))
	for i, item := range items {
		itemInfos[i].FromDB(&item)
	}

	SuccessWithData(c, define.SearchResp{
		List:  itemInfos,
		Total: total,
	})
}

func SearchItems(c *gin.Context) {
	var req define.SearchItemsReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, err)
		return
	}

	pagination := common.Pagination{
		Disable: req.NoPaging,
		Page:    req.Page,
		Size:    req.PageSize,
	}

	// 字段筛选仅在指定类别时可用
	if req.CategoryID == 0 {
		req.Filters = nil
	}

	items, total, err := service.SearchItems(req.CategoryID, req.Name, req.TagIDs, req.CollectionIDs, req.Filters, pagination)
	if err != nil {
		Fail(c, err)
		return
	}

	itemDetails := make([]define.ItemDetail, len(items))
	for i, item := range items {
		itemDetails[i].FromDB(&item)
	}

	SuccessWithData(c, define.SearchResp{
		List:  itemDetails,
		Total: total,
	})
}

func GetItem(c *gin.Context) {
	id, err := GetID(c, "id")
	if err != nil {
		Fail(c, err)
		return
	}

	uniqueFields := map[string]interface{}{"id": id}
	preloads := []string{"Category", "Tags", "Collections", "Values", "Values.Field"}
	item, err := dao.Get[model.Item](conn.GetDB(), uniqueFields, preloads...)
	if err != nil {
		Fail(c, err)
		return
	}

	itemDetail := define.ItemDetail{}
	itemDetail.FromDB(&item)

	SuccessWithData(c, itemDetail)
}
