package model

import "gorm.io/gorm"

const (
	FieldTypeString   = "string"
	FieldTypeInt      = "int"
	FieldTypeBool     = "bool"
	FieldTypeDate     = "date"
	FieldTypeDatetime = "datetime"
)

// Field 字段
type Field struct {
	gorm.Model
	CategoryID uint   `gorm:"index"`
	Name       string `gorm:"not null"`
	Type       string `gorm:"not null"`
	IsArray    bool   `gorm:"default:false"`
	Required   bool
}

func (f Field) TableName() string {
	return "fields"
}

func (f Field) GetID() uint {
	return f.Model.ID
}

func (f Field) IsDeleted() bool {
	return f.Model.DeletedAt.Valid
}
