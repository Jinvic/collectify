package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	FieldTypeString = iota + 1
	FieldTypeInt
	FieldTypeBool
	FieldTypeDatetime
)

var FieldTypeNames = map[int]string{
	FieldTypeString:   "string",
	FieldTypeInt:      "int",
	FieldTypeBool:     "bool",
	FieldTypeDatetime: "datetime",
}

// Field 字段
type Field struct {
	gorm.Model
	CategoryID uint   `gorm:"index;uniqueIndex:idx_field_category_name"`
	Name       string `gorm:"not null;uniqueIndex:idx_field_category_name"`
	Type       int    `gorm:"not null"`
	IsArray    bool   `gorm:"default:false"`
	Required   bool   `gorm:"default:false"`
}

func (f Field) TableName() string {
	return "fields"
}

func (f Field) GetID() uint {
	return f.ID
}

func (f Field) IsDeleted() bool {
	return f.DeletedAt.Valid
}

// ItemFieldValue 自定义字段值
type ItemFieldValue struct {
	gorm.Model
	ItemID  uint `gorm:"index"`
	FieldID uint `gorm:"index"`

	ValueString *string    `gorm:"index"`
	ValueInt    *int       `gorm:"index"`
	ValueBool   *bool      `gorm:"index"`
	ValueTime   *time.Time `gorm:"index"`

	Item  Item  `gorm:"foreignKey:ItemID"`
	Field Field `gorm:"foreignKey:FieldID"`
}

func (i ItemFieldValue) TableName() string {
	return "item_field_values"
}

func (i ItemFieldValue) GetID() uint {
	return i.ID
}

func (i ItemFieldValue) IsDeleted() bool {
	return i.DeletedAt.Valid
}
