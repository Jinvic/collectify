package service

import (
	"collectify/internal/config"
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"

	"gorm.io/gorm"
)

// DeleteField 删除字段
func DeleteField(fieldID uint) error {
	db := db.GetDB()
	cfg := config.GetConfig()
	isSoftDelete := cfg.RecycleBin.Enable

	err := db.Transaction(func(tx *gorm.DB) error {
		var uniqueFields map[string]interface{}
		var err error

		// 删除字段值
		uniqueFields = map[string]interface{}{"field_id": fieldID}
		err = dao.Delete[model.ItemFieldValue](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		// 删除字段
		uniqueFields = map[string]interface{}{"id": fieldID}
		err = dao.Delete[model.Field](tx, uniqueFields, isSoftDelete)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// RestoreField 恢复字段
func RestoreField(fieldID uint) error {
	db := db.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		var uniqueFields map[string]interface{}
		var err error

		// 尝试恢复分类
		var field model.Field
		err = tx.Unscoped().Model(&field).Where("id = ?", fieldID).First(&field).Error
		if err != nil {
			return err
		}
		err = tryRestoreCategory(tx, field.CategoryID)
		if err != nil {
			return err
		}

		// 恢复字段
		uniqueFields = map[string]interface{}{"id": fieldID}
		err = dao.Restore[model.Field](tx, uniqueFields)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
