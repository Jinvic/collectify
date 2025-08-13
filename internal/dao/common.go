package dao

import (
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	"errors"

	"gorm.io/gorm"
)

func DuplicateCheck[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) (id uint, isDeleted bool, err error) {
	var t T
	err = tx.Model(&t).Where(uniqueFields).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return t.GetID(), t.IsDeleted(), nil
}

func Create[T model.GormModel](tx *gorm.DB, data *T) (*T, error) {
	err := tx.Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
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

func Get[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) (T, error) {
	var t T
	return t, tx.Model(&t).Where(uniqueFields).First(&t).Error
}

func GetList[T model.GormModel](tx *gorm.DB, filters map[string]interface{}, p common.Pagination) ([]T, int64, error) {
	var t []T
	query := tx.Model(&t)

	query = query.Where(filters)

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
