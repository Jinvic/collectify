package service

import (
	"collectify/internal/config"
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"
	define "collectify/internal/model/define"
	"fmt"

	"gorm.io/gorm"
)

// CreateItem 创建收藏品
func CreateItem(categoryID uint, title string, values []define.ItemFieldValue) error {
	db := db.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建收藏品
		item := &model.Item{
			CategoryID: categoryID,
			Title:      title,
		}
		if err := dao.Create(tx, item); err != nil {
			return err
		}

		// 获取分类信息，并预加载字段
		uniqueFields := map[string]interface{}{"id": categoryID}
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
func UpdateItem(itemID uint, title string, values []define.ItemFieldValue) error {
	db := db.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		uniqueFields := map[string]interface{}{"id": itemID}
		preloads := []string{
			"Category",        // 预加载所属分类
			"Category.Fields", // 预加载分类的字段
			"Values",          // 预加载已有字段值
		}
		item, err := dao.Get[model.Item](tx, uniqueFields, preloads...)
		if err != nil {
			return err
		}

		// 如果标题有变化，则更新标题
		if title != item.Title {
			updateFields := map[string]interface{}{"title": title}
			if err := dao.Update[model.Item](tx, uniqueFields, updateFields); err != nil {
				return err
			}
		}

		fieldMap := make(map[uint]model.Field)
		for _, field := range item.Category.Fields {
			fieldMap[field.ID] = field
		}

		// 删除原有的字段值
		uniqueFields = map[string]interface{}{"item_id": itemID}
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
	db := db.GetDB()
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
	db := db.GetDB()

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
