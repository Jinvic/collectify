package model

import "gorm.io/gorm"

// Category 类别
type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null" json:"name"` // 类别名称唯一

	// 反向关联
	Items  []Item  `gorm:"foreignKey:CategoryID"` // 使用该类别的藏品
	Fields []Field `gorm:"foreignKey:CategoryID"` // 类别包含的字段
}
