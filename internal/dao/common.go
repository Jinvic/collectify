package dao

import (
	common "collectify/internal/model/common"

	"gorm.io/gorm"
)

func IsExists[T any](tx *gorm.DB, id uint) (bool, error) {
	var count int64
	var t T
	err := tx.Model(&t).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func Create[T any](tx *gorm.DB, data *T) error {
	return tx.Create(data).Error
}

func SoftDelete[T any](tx *gorm.DB, id uint) error {
	var t T
	return tx.Delete(&t, id).Error
}

func Restore[T any](tx *gorm.DB, id uint) error {
	var t T
	return tx.Unscoped().Model(&t).Where("id = ?", id).Update("deleted_at", nil).Error
}

func HardDelete[T any](tx *gorm.DB, id uint) error {
	var t T
	return tx.Unscoped().Delete(&t, id).Error
}

func Get[T any](tx *gorm.DB, id uint) (T, error) {
	var t T
	return t, tx.First(&t, id).Error
}

func GetList[T any](tx *gorm.DB, filters map[string]interface{}, p common.Pagination) ([]T, int64, error) {
	var t []T
	query := tx.Model(&t)

	if filters["name"] != nil {
		query = query.Where("name = ?", filters["name"])
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&t).Error; err != nil {
		return nil, 0, err
	}
	return t, total, nil
}
