package dao

import (
	"collectify/internal/config"
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	"encoding/json"
	"strings"

	"gorm.io/gorm"
)

func UpdateItem(tx *gorm.DB, id uint, updates map[string]interface{}) error {
	var item model.Item
	if err := tx.First(&item, id).Error; err != nil {
		return err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(item.Metadata, &metadata); err != nil {
		return err
	}

	for key, value := range updates {
		metadata[key] = value
	}

	updatedMetadata, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	item.Metadata = updatedMetadata
	return tx.Save(&item).Error
}

func GetItemList(tx *gorm.DB, filters map[string]interface{}, p common.Pagination) ([]model.Item, int64, error) {
	var items []model.Item

	query := tx.Model(&model.Item{})

	for key, value := range filters {
		switch key {
		case "category":
			query = query.Where("category_id = ?", value)
		case "tag":
			query = query.Joins("JOIN item_tags ON items.id = item_tags.item_id").
				Where("item_tags.tag_id = ?", value)
		case "collection":
			query = query.Joins("JOIN collection_items ON items.id = collection_items.item_id").
				Where("collection_items.collection_id = ?", value)
		case "title":
			query = query.Where("title LIKE ?", "%"+value.(string)+"%")
		default:
			query = queryMetadata(query, key, value)
		}

	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func queryMetadata(tx *gorm.DB, key string, value interface{}) *gorm.DB {
	if strings.HasPrefix(key, "metadata.") {
		cfg := config.GetConfig()
		jsonPath := "$." + strings.TrimPrefix(key, "metadata.")

		switch cfg.Database.Type {
		case "sqlite":
			// SQLite: json_extract(metadata, '$.color') = 'red'
			return tx.Where("json_extract(metadata, ?) = ?", jsonPath, value)

		case "mysql":
			// MySQL: JSON_UNQUOTE(JSON_EXTRACT(metadata, '$.color')) = 'red'
			return tx.Where("JSON_UNQUOTE(JSON_EXTRACT(metadata, ?)) = ?", jsonPath, value)

		case "postgres", "postgresql":
			// PostgreSQL: metadata ->> 'color' = 'red'
			field := strings.TrimPrefix(jsonPath, "$.")
			return tx.Where("metadata ->> ? = ?", field, value)

		default:
			// 未知数据库类型，不作处理
			return tx
		}
	}
	return tx
}
