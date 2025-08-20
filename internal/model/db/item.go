package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ItemStatusTodo       = iota + 1 // 待完成
	ItemStatusInProgress            // 进行中
	ItemStatusPaused                // 暂停
	ItemStatusAbandoned             // 放弃
	ItemStatusCompleted             // 完成
)

// Item 收藏品
type Item struct {
	gorm.Model
	Name        string     `gorm:"not null;index" json:"name"`                                     // 名称
	CategoryID  uint       `gorm:"not null;index" json:"category_id"`                              // 关联的类别ID
	Status      int        `gorm:"not null;default:1;index" json:"status"`                         // 状态
	Rating      *float64   `gorm:"type:decimal(3,1);check:rating>=0 and rating<=10" json:"rating"` // 评分
	Description string     `gorm:"type:text" json:"description"`                                   // 简介
	Notes       string     `gorm:"type:text" json:"notes"`                                         // 感想
	CoverURL    string     `json:"cover_url"`                                                      // 封面图
	SourceURL   string     `json:"source_url"`                                                     // 外部链接
	CompletedAt *time.Time `json:"completed_at"`                                                   // 完成时间
	Priority    int        `gorm:"default:0" json:"priority"`                                      // 优先级

	// 关联关系
	Category    Category         `gorm:"foreignKey:CategoryID"`                          // 所属类别
	Tags        []Tag            `gorm:"many2many:item_tags;" json:"tags"`               // 多个标签
	Collections []Collection     `gorm:"many2many:collection_items;" json:"collections"` // 所属的收藏夹
	Values      []ItemFieldValue `gorm:"foreignKey:ItemID"`                              // 自定义字段值
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
