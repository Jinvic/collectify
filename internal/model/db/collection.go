package model

import "gorm.io/gorm"

// Collection 收藏夹
type Collection struct {
	gorm.Model
	Name        string `gorm:"not null;index" json:"name"` // 收藏夹名称
	Description string `json:"description"`                // 描述

	// 关联的藏品
	Items []Item `gorm:"many2many:collection_items;" json:"items"` // 包含的藏品（可跨类型）
}

func (c Collection) TableName() string {
	return "collections"
}
