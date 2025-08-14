package service

import (
	"collectify/internal/dao"
	"collectify/internal/db"
	model "collectify/internal/model/db"

	"gorm.io/gorm"
)

var DeleteByFilterFuncs = map[string]func(tx *gorm.DB, filters []dao.Filter, isSoftDelete bool) error{
	model.ModelTypeCategory:   dao.DeleteByFilter[model.Category],
	model.ModelTypeCollection: dao.DeleteByFilter[model.Collection],
	model.ModelTypeField:      dao.DeleteByFilter[model.Field],
	model.ModelTypeItem:       dao.DeleteByFilter[model.Item],
	model.ModelTypeTag:        dao.DeleteByFilter[model.Tag],
	model.ModelTypeIFV:        dao.DeleteByFilter[model.ItemFieldValue],
}

func ClearRecycleBin() error {
	db := db.GetDB()

	filters := []dao.Filter{
		{
			Where: "deleted_at is not null",
		},
	}
	for _, fn := range DeleteByFilterFuncs {
		err := fn(db, filters, false)
		if err != nil {
			return err
		}
	}
	return nil
}
