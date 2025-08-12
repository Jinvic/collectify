package model

type GormModel interface {
	GetID() uint
	IsDeleted() bool
}

const (
	ModelTypeCategory   = "category"
	ModelTypeCollection = "collection"
	ModelTypeField      = "field"
	ModelTypeItem       = "item"
	ModelTypeTag        = "tag"
)
