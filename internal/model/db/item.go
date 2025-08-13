package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Item 收藏品
type Item struct {
	gorm.Model
	Title      string          `gorm:"not null;index" json:"title"`       // 名称
	CategoryID uint            `gorm:"not null;index" json:"category_id"` // 关联的类别ID
	Metadata   json.RawMessage `gorm:"type:json" json:"metadata"`         // 自定义字段 (JSON)

	// 关联关系
	Category    Category     `gorm:"foreignKey:CategoryID"`                          // 所属类别
	Tags        []Tag        `gorm:"many2many:item_tags;" json:"tags"`               // 多个标签
	Collections []Collection `gorm:"many2many:collection_items;" json:"collections"` // 所属的收藏夹
}

func (i Item) TableName() string {
	return "items"
}

func (i Item) GetID() uint {
	return i.ID
}

func (i Item) IsDeleted() bool {
	return i.DeletedAt.Valid
}
