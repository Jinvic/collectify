package dao

import (
	common "collectify/internal/model/common"
	model "collectify/internal/model/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Filter struct {
	Where string
	Args  []interface{}
}

type OrderBy struct {
	Column string
	Desc   bool
}

type Join struct {
	Table string
	On    string
}

// uniqueFields: 业务上需要检查唯一性的字段，如 name, email
// filters: 查询时的附加条件，如排除当前记录、状态过滤等
func DuplicateCheck[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, filters []Filter) (id uint, isDeleted bool, err error) {
	var t T
	query := tx.Unscoped().Model(&t).Where(uniqueFields)
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

func IsDeleted[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) (bool, error) {
	var t T
	err := tx.Unscoped().Model(&t).Where(uniqueFields).First(&t).Error
	if err != nil {
		return false, err
	}
	return t.IsDeleted(), nil
}

func Create[T model.GormModel](tx *gorm.DB, data *T) error {
	err := tx.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func Restore[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}) error {
	var t T
	return tx.Unscoped().Model(&t).Where(uniqueFields).Update("deleted_at", nil).Error
}

func RestoreByFilter[T model.GormModel](tx *gorm.DB, filters []Filter) error {
	var t T

	query := tx.Unscoped().Model(&t)
	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}
	return query.Update("deleted_at", nil).Error
}

func Delete[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, isSoftDelete bool) error {
	var t T

	query := tx.Model(&t)
	if !isSoftDelete {
		query = query.Unscoped()
	}
	return query.Where(uniqueFields).Delete(&t).Error
}

func DeleteByFilter[T model.GormModel](tx *gorm.DB, filters []Filter, isSoftDelete bool) error {
	var t T

	query := tx.Model(&t)
	if !isSoftDelete {
		query = query.Unscoped()
	}

	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}

	return query.Delete(&t).Error
}

func Get[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, preloads ...string) (T, error) {
	var t T
	query := tx.Model(&t)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return t, query.Where(uniqueFields).First(&t).Error
}

func GetList[T model.GormModel](tx *gorm.DB, filters []Filter, orderBy []OrderBy, p common.Pagination, preloads ...string) ([]T, int64, error) {
	var t []T
	query := tx.Model(&t)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	for _, orderBy := range orderBy {
		var sort string
		if orderBy.Desc {
			sort = "DESC"
		} else {
			sort = "ASC"
		}
		query = query.Order(orderBy.Column + " " + sort)
	}

	if !p.Disable {
		query = query.Offset(p.GetOffset()).Limit(p.GetLimit())
	}

	if err := query.Find(&t).Error; err != nil {
		return nil, 0, err
	}
	return t, total, nil
}

func Pluck[T model.GormModel, R any](tx *gorm.DB, column string, joins []Join, filters []Filter, distinct bool) ([]R, error) {
	var t []T
	var r []R
	query := tx.Model(&t)
	for _, join := range joins {
		joinStr := fmt.Sprintf("LEFT JOIN %s ON %s", join.Table, join.On)
		query = query.Joins(joinStr)
	}
	for _, filter := range filters {
		query = query.Where(filter.Where, filter.Args...)
	}
	if distinct {
		query = query.Distinct(column)
	}
	err := query.Pluck(column, &r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Update[T model.GormModel](tx *gorm.DB, uniqueFields map[string]interface{}, updateFields map[string]interface{}) error {
	var t T
	return tx.Model(&t).Where(uniqueFields).Updates(updateFields).Error
}

func Associate[T1 model.GormModel, T2 model.GormModel](tx *gorm.DB, id1, id2 uint, association string) error {
	uniqueFields := map[string]interface{}{"id": id1}
	t1, err := Get[T1](tx, uniqueFields)
	if err != nil {
		return err
	}

	uniqueFields = map[string]interface{}{"id": id2}
	t2, err := Get[T2](tx, uniqueFields)
	if err != nil {
		return err
	}

	err = tx.Model(&t1).Association(association).Append(&t2)
	if err != nil {
		return err
	}

	return nil
}

func Disassociate[T1 model.GormModel, T2 model.GormModel](tx *gorm.DB, id1, id2 uint, association string) error {
	uniqueFields := map[string]interface{}{"id": id1}
	t1, err := Get[T1](tx, uniqueFields)
	if err != nil {
		return err
	}
	
	uniqueFields = map[string]interface{}{"id": id2}
	t2, err := Get[T2](tx, uniqueFields)
	if err != nil {
		return err
	}

	err = tx.Model(&t1).Association(association).Delete(&t2)
	if err != nil {
		return err
	}

	return nil
}
