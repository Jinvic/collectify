package service

import (
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"
	"collectify/internal/pkg/e"

	"gorm.io/gorm"
)

var restoreFuncs = map[string]func(tx *gorm.DB, uniqueFields map[string]interface{}) error{
	model.ModelTypeCategory:   dao.Restore[model.Category],
	model.ModelTypeCollection: dao.Restore[model.Collection],
	model.ModelTypeField:      dao.Restore[model.Field],
	model.ModelTypeItem:       dao.Restore[model.Item],
	model.ModelTypeTag:        dao.Restore[model.Tag],
}

var hardDeleteFuncs = map[string]func(tx *gorm.DB, uniqueFields map[string]interface{}) error{
	model.ModelTypeCategory:   dao.HardDelete[model.Category],
	model.ModelTypeCollection: dao.HardDelete[model.Collection],
	model.ModelTypeField:      dao.HardDelete[model.Field],
	model.ModelTypeItem:       dao.HardDelete[model.Item],
	model.ModelTypeTag:        dao.HardDelete[model.Tag],
}

var softDeleteFuncs = map[string]func(tx *gorm.DB, uniqueFields map[string]interface{}) error{
	model.ModelTypeCategory:   dao.SoftDelete[model.Category],
	model.ModelTypeCollection: dao.SoftDelete[model.Collection],
	model.ModelTypeField:      dao.SoftDelete[model.Field],
	model.ModelTypeItem:       dao.SoftDelete[model.Item],
	model.ModelTypeTag:        dao.SoftDelete[model.Tag],
}

func Restore(typ string, uniqueFields map[string]interface{}) error {
	db := db.GetDB()
	fn, ok := restoreFuncs[typ]
	if !ok {
		return e.ErrInvalidParams
	}
	return fn(db, uniqueFields)
}

func HardDelete(typ string, uniqueFields map[string]interface{}) error {
	db := db.GetDB()
	fn, ok := hardDeleteFuncs[typ]
	if !ok {
		return e.ErrInvalidParams
	}
	return fn(db, uniqueFields)
}

func SoftDelete(typ string, uniqueFields map[string]interface{}) error {
	db := db.GetDB()
	fn, ok := softDeleteFuncs[typ]
	if !ok {
		return e.ErrInvalidParams
	}
	return fn(db, uniqueFields)
}
