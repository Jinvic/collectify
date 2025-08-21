package service

import (
	"collectify/internal/config"
	"collectify/internal/dao"
	"collectify/internal/conn"
	model "collectify/internal/model/db"

	"gorm.io/gorm"
)

// DeleteCategory 删除分类
func DeleteCategory(categoryID uint) error {
	db := conn.GetDB()
	cfg := config.GetConfig()
	isSoftDelete := cfg.RecycleBin.Enable

	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取分类下的字段ID
		var fieldIDs []uint
		err := tx.Model(&model.Field{}).Where("category_id = ?", categoryID).Pluck("id", &fieldIDs).Error
		if err != nil {
			return err
		}

		// 删除分类下的字段值
		filters := []dao.Filter{
			{
				Where: "field_id IN (?)",
				Args:  []interface{}{fieldIDs},
			},
		}
		err = dao.DeleteByFilter[model.ItemFieldValue](tx, filters, isSoftDelete)
		if err != nil {
			return err
		}

		var uniqueFields map[string]interface{}

		// 删除分类下的字段
		uniqueFields = map[string]interface{}{"category_id": categoryID}
		err = dao.Delete[model.Field](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		// 删除分类下的收藏品
		uniqueFields = map[string]interface{}{"category_id": categoryID}
		err = dao.Delete[model.Item](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		// 删除分类
		uniqueFields = map[string]interface{}{"id": categoryID}
		err = dao.Delete[model.Category](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// RestoreCategory 恢复分类
func RestoreCategory(categoryID uint) error {
	db := conn.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取分类下的字段ID
		var fieldIDs []uint
		err := tx.Unscoped().Model(&model.Field{}).Where("category_id = ?", categoryID).Pluck("id", &fieldIDs).Error
		if err != nil {
			return err
		}

		var uniqueFields map[string]interface{}

		// 恢复分类
		uniqueFields = map[string]interface{}{"id": categoryID}
		err = dao.Restore[model.Category](tx, uniqueFields)
		if err != nil {
			return err
		}

		// 恢复分类下的字段
		uniqueFields = map[string]interface{}{"category_id": categoryID}
		err = dao.Restore[model.Field](tx, uniqueFields)
		if err != nil {
			return err
		}

		// 恢复分类下的收藏品
		uniqueFields = map[string]interface{}{"category_id": categoryID}
		err = dao.Restore[model.Item](tx, uniqueFields)
		if err != nil {
			return err
		}

		// 恢复分类下的字段值
		filters := []dao.Filter{
			{
				Where: "field_id IN (?)",
				Args:  []interface{}{fieldIDs},
			},
		}
		err = dao.RestoreByFilter[model.ItemFieldValue](tx, filters)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// tryRestoreCategory 尝试仅恢复分类
func tryRestoreCategory(tx *gorm.DB, categoryID uint) error {
	var err error

	// 检查分类是否被删除
	isDeleted, err := dao.IsDeleted[model.Category](tx, map[string]interface{}{"id": categoryID})
	if err != nil {
		return err
	}
	// 如果分类未被删除，则不恢复
	if !isDeleted {
		return nil
	}

	// 恢复分类
	uniqueFields := map[string]interface{}{"id": categoryID}
	err = dao.Restore[model.Category](tx, uniqueFields)
	if err != nil {
		return err
	}

	return nil
}

// tryRestoreCategory 尝试恢复分类和分类下的字段，不恢复分类下的收藏品和字段值
//
// 注意：此函数会恢复分类下的所有字段，因为：
// 1. 字段是分类的元数据，与分类强绑定
// 2. ItemFieldValue 依赖字段，恢复藏品时需确保字段存在
// 如果某些字段不应被恢复，应在删除分类前硬删除它们。
func tryRestoreCategoryWithFields(tx *gorm.DB, categoryID uint) error {
	var err error

	// 检查分类是否被删除
	isDeleted, err := dao.IsDeleted[model.Category](tx, map[string]interface{}{"id": categoryID})
	if err != nil {
		return err
	}
	// 如果分类未被删除，则不恢复
	if !isDeleted {
		return nil
	}

	// 获取分类下的字段ID
	var fieldIDs []uint
	err = tx.Unscoped().Model(&model.Field{}).Where("category_id = ?", categoryID).Pluck("id", &fieldIDs).Error
	if err != nil {
		return err
	}

	var uniqueFields map[string]interface{}

	// 恢复分类
	uniqueFields = map[string]interface{}{"id": categoryID}
	err = dao.Restore[model.Category](tx, uniqueFields)
	if err != nil {
		return err
	}

	// 恢复分类下的字段
	uniqueFields = map[string]interface{}{"category_id": categoryID}
	err = dao.Restore[model.Field](tx, uniqueFields)
	if err != nil {
		return err
	}

	return nil
}
