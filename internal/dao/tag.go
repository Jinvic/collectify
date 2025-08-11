package dao

import (
	model "collectify/internal/model/db"

	"gorm.io/gorm"
)

func AddTagToItem(tx *gorm.DB, itemID uint, tagID uint) error {
	return tx.Model(&model.Item{}).Where("id = ?", itemID).Association("Tags").Append(&model.Tag{Model: gorm.Model{ID: tagID}})
}

func RemoveTagFromItem(tx *gorm.DB, itemID uint, tagID uint) error {
	return tx.Model(&model.Item{}).Where("id = ?", itemID).Association("Tags").Delete(&model.Tag{Model: gorm.Model{ID: tagID}})
}

func GetItemTags(tx *gorm.DB, itemID uint) ([]model.Tag, error) {
	var item model.Item
	if err := tx.Preload(model.Tag{}.TableName()).First(&item, itemID).Error; err != nil {
		return nil, err
	}
	return item.Tags, nil
}

func GetTagItems(tx *gorm.DB, tagID uint) ([]model.Item, error) {
	var tag model.Tag
	if err := tx.Preload(model.Item{}.TableName()).First(&tag, tagID).Error; err != nil {
		return nil, err
	}
	return tag.Items, nil
}
