package dao

import (
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	"errors"

	"gorm.io/gorm"
)

type Filter struct {
	Where string
	Args  []interface{}
}

// uniqueFields: 业务上需要检查唯一性的字段，如 name, email
// filters: 查询时的附加条件，如排除当前记录、状态过滤等
func DuplicateCheck[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, filters []Filter) (id uint, isDeleted bool, err error) {
	var t T
	query := tx.Model(&t).Where(uniqueFields)
	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}
	err = query.First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return t.GetID(), t.IsDeleted(), nil
}

func Create[T model.GormModel](tx *gorm.DB, data *T) error {
	err := tx.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func SoftDelete[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) error {
	var t T
	return tx.Model(&t).Where(uniqueFields).Delete(&t).Error
}

func Restore[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) error {
	var t T
	return tx.Unscoped().Model(&t).Where(uniqueFields).Update("deleted_at", nil).Error
}

func HardDelete[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) error {
	var t T
	return tx.Unscoped().Model(&t).Where(uniqueFields).Delete(&t).Error
}

func Get[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, preloads ...string) (T, error) {
	var t T
	query := tx.Model(&t).Where(uniqueFields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return t, query.First(&t).Error
}

func GetList[T model.GormModel](tx *gorm.DB, filters []Filter, p common.Pagination) ([]T, int64, error) {
	var t []T
	query := tx.Model(&t)

	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if !p.Disable {
		query = query.Offset(p.GetOffset()).Limit(p.GetLimit())
	}

	if err := query.Find(&t).Error; err != nil {
		return nil, 0, err
	}
	return t, total, nil
}

func Update[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, updateFields map[string]interface{}) error {
	var t T
	return tx.Model(&t).Where(uniqueFields).Updates(updateFields).Error
}
