package service

import (
	"collectify/internal/config"
	"collectify/internal/conn"
	"collectify/internal/dao"
	"collectify/internal/model/common"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateItem 创建收藏品
func CreateItem(item *model.Item, values []define.ItemFieldValue) error {
	db := conn.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建收藏品
		if err := dao.Create(tx, item); err != nil {
			return err
		}

		// 获取分类信息，并预加载字段
		uniqueFields := map[string]interface{}{"id": item.CategoryID}
		preloads := []string{"Fields"}
		category, err := dao.Get[model.Category](tx, uniqueFields, preloads...)
		if err != nil {
			return err
		}

		var fieldMap = make(map[uint]model.Field)
		for _, field := range category.Fields {
			fieldMap[field.ID] = field
		}

		// 遍历并保存每个字段值
		for _, value := range values {
			// 检查字段是否存在
			field, ok := fieldMap[value.FieldID]
			if !ok {
				return fmt.Errorf("field not found: %d", value.FieldID)
			}

			// 创建字段值
			creator := dao.NewFieldValueCreator(tx, item.ID, field, value.Value)
			if err := creator.Create(); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// UpdateItem 更新收藏品信息
func UpdateItem(item *model.Item, values []define.ItemFieldValue) error {
	db := conn.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		uniqueFields := map[string]interface{}{"id": item.ID}
		preloads := []string{
			"Category",        // 预加载所属分类
			"Category.Fields", // 预加载分类的字段
			"Values",          // 预加载已有字段值
		}
		oldItem, err := dao.Get[model.Item](tx, uniqueFields, preloads...)
		if err != nil {
			return err
		}

		// 更新收藏品信息
		updateFields := map[string]interface{}{
			"name":        item.Name,
			"status":      item.Status,
			"rating":      item.Rating,
			"description": item.Description,
			"notes":       item.Notes,
			"cover_url":   item.CoverURL,
			"source_url":  item.SourceURL,
			"priority":    item.Priority,
		}

		// 更新完成时间
		now := time.Now()
		if oldItem.Status != model.ItemStatusCompleted &&
			item.Status == model.ItemStatusCompleted {
			updateFields["completed_at"] = &now
		}
		if oldItem.Status == model.ItemStatusCompleted &&
			item.Status != model.ItemStatusCompleted {
			updateFields["completed_at"] = nil
		}

		if err := dao.Update[model.Item](tx, uniqueFields, updateFields); err != nil {
			return err
		}

		fieldMap := make(map[uint]model.Field)
		for _, field := range oldItem.Category.Fields {
			fieldMap[field.ID] = field
		}

		// 删除原有的字段值
		uniqueFields = map[string]interface{}{"item_id": item.ID}
		err = dao.Delete[model.ItemFieldValue](tx, uniqueFields, false) // 硬删除字段值
		if err != nil {
			return err
		}

		// 创建新的字段值
		for _, value := range values {
			// 检查字段是否存在
			field, ok := fieldMap[value.FieldID]
			if !ok {
				return fmt.Errorf("field not found: %d", value.FieldID)
			}

			// 创建字段值
			creator := dao.NewFieldValueCreator(tx, item.ID, field, value.Value)
			if err := creator.Create(); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// DeleteItem 删除收藏品
func DeleteItem(itemID uint) error {
	db := conn.GetDB()
	cfg := config.GetConfig()
	isSoftDelete := cfg.RecycleBin.Enable

	err := db.Transaction(func(tx *gorm.DB) error {
		var uniqueFields map[string]interface{}
		var err error

		// 删除收藏品下的字段值
		uniqueFields = map[string]interface{}{"item_id": itemID}
		err = dao.Delete[model.ItemFieldValue](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		// 删除收藏品
		uniqueFields = map[string]interface{}{"id": itemID}
		err = dao.Delete[model.Item](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// RestoreItem 恢复收藏品
func RestoreItem(itemID uint) error {
	db := conn.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		var uniqueFields map[string]interface{}
		var err error

		// 尝试恢复分类和分类下的字段
		var item model.Item
		err = tx.Unscoped().Model(&item).Where("id = ?", itemID).First(&item).Error
		if err != nil {
			return err
		}
		err = tryRestoreCategoryWithFields(tx, item.CategoryID)
		if err != nil {
			return err
		}

		// 恢复收藏品
		uniqueFields = map[string]interface{}{"id": itemID}
		err = dao.Restore[model.Item](tx, uniqueFields)
		if err != nil {
			return err
		}

		// 恢复收藏品下的字段值
		uniqueFields = map[string]interface{}{"item_id": itemID}
		err = dao.Restore[model.ItemFieldValue](tx, uniqueFields)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// ListItems 列出收藏品
func ListItems(p common.Pagination) ([]model.Item, int64, error) {
	db := conn.GetDB()

	filters := []dao.Filter{}
	orderBy := []dao.OrderBy{
		{
			Column: "updated_at",
			Desc:   true,
		},
	}
	preloads := []string{
		"Category",
		"Tags",
	}

	items, total, err := dao.GetList[model.Item](db, filters, orderBy, p, preloads...)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// SearchItems 搜索收藏品
func SearchItems(categoryID uint, name string, tagIDs []uint, collectionIDs []uint, fieldFilters map[uint]interface{}, p common.Pagination) ([]model.Item, int64, error) {
	db := conn.GetDB()

	// 预加载关联表
	joins := []dao.Join{
		{
			Table: "item_tags",
			On:    "items.id = item_tags.item_id",
		},
		{
			Table: "tags",
			On:    "item_tags.tag_id = tags.id",
		},
		{
			Table: "collection_items",
			On:    "items.id = collection_items.item_id",
		},
		{
			Table: "collections",
			On:    "collection_items.collection_id = collections.id",
		},
		{
			Table: "item_field_values",
			On:    "items.id = item_field_values.item_id",
		},
	}

	
	filters := []dao.Filter{}
	
	// 筛选条件
	if categoryID > 0 {
		filters = append(filters, dao.Filter{
			Where: "items.category_id = ?",
			Args:  []interface{}{categoryID},
		})
	}
	if name != "" {
		filters = append(filters, dao.Filter{
			Where: "items.name LIKE ?",
			Args:  []interface{}{"%" + name + "%"},
		})
	}
	if len(tagIDs) > 0 {
		filters = append(filters, dao.Filter{
			Where: "tags.id IN ?",
			Args:  []interface{}{tagIDs},
		})
	}
	if len(collectionIDs) > 0 {
		filters = append(filters, dao.Filter{
			Where: "collections.id IN ?",
			Args:  []interface{}{collectionIDs},
		})
	}

	// 更新时间逆序排序
	orderBy := []dao.OrderBy{
		{
			Column: "updated_at",
			Desc:   true,
		},
	}

	// 预加载关联表
	preloads := []string{
		"Category",
		"Tags",
		"Collections",
		"Values",
		"Values.Field",
	}

	var items []model.Item
	var total int64
	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取分类信息，并预加载字段
		fieldMap := make(map[uint]model.Field)
		category, err := dao.Get[model.Category](tx, map[string]interface{}{"id": categoryID}, "Fields")
		if err != nil {
			return err
		}
		for _, field := range category.Fields {
			fieldMap[field.ID] = field
		}

		// 遍历并添加字段值过滤条件
		for key, value := range fieldFilters {
			field, ok := fieldMap[key]
			if !ok {
				return fmt.Errorf("field not found: %d", key)
			}

			builder := dao.NewFieldValueQueryBuilder(tx, field, value)
			filter, err := builder.Build()
			if err != nil {
				return err
			}
			filters = append(filters, filter)
		}

		// 先查询出所有符合条件的收藏品ID，再预加载关联表，避免笛卡尔积查询
		itemIDs, err := dao.Pluck[model.Item, uint](tx, "items.id", joins, filters, true)
		if err != nil {
			return err
		}

		filters = []dao.Filter{
			{
				Where: "items.id IN ?",
				Args:  []interface{}{itemIDs},
			},
		}
		items, total, err = dao.GetList[model.Item](tx, filters, orderBy, p, preloads...)
		if err != nil {
			return err
		}

		return nil
	})

	return items, total, err
}
