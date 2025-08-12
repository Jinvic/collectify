package model

import "gorm.io/gorm"

// Tag 标签
type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null" json:"name"` // 标签名唯一

	// 反向关联
	Items []Item `gorm:"many2many:item_tags;"` // 使用该标签的藏品
}

func (t Tag) TableName() string {
	return "tags"
}

func (t Tag) GetID() uint {
	return t.Model.ID
}

func (t Tag) IsDeleted() bool {
	return t.Model.DeletedAt.Valid
}
