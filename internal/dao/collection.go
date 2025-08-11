package dao

import (
	model "collectify/internal/model/db"

	"gorm.io/gorm"
)

func AddItemToCollection(tx *gorm.DB, collectionID uint, itemID uint) error {
	return tx.Model(&model.Collection{}).Where("id = ?", collectionID).Association("Items").Append(&model.Item{Model: gorm.Model{ID: itemID}})
}

func RemoveItemFromCollection(tx *gorm.DB, collectionID uint, itemID uint) error {
	return tx.Model(&model.Collection{}).Where("id = ?", collectionID).Association("Items").Delete(&model.Item{Model: gorm.Model{ID: itemID}})
}

func GetCollectionItems(tx *gorm.DB, collectionID uint) ([]model.Item, error) {
	var collection model.Collection
	if err := tx.Preload(model.Item{}.TableName()).First(&collection, collectionID).Error; err != nil {
		return nil, err
	}
	return collection.Items, nil
}

func GetItemCollections(tx *gorm.DB, itemID uint) ([]model.Collection, error) {
	var item model.Item
	if err := tx.Preload(model.Collection{}.TableName()).First(&item, itemID).Error; err != nil {
		return nil, err
	}
	return item.Collections, nil
}
