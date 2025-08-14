package dao

import (
	model "collectify/internal/model/db"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type FieldValueCreator struct {
	tx     *gorm.DB
	itemID uint
	field  model.Field
	value  interface{}
}

func NewFieldValueCreator(tx *gorm.DB, itemID uint, field model.Field, value interface{}) *FieldValueCreator {
	return &FieldValueCreator{
		tx:     tx,
		itemID: itemID,
		field:  field,
		value:  value,
	}
}

func (c *FieldValueCreator) Create() error {
	// 必填校验统一处理
	if c.field.Required {
		if err := c.validateRequired(); err != nil {
			return err
		}
	}

	// 根据是否为数组分发处理
	if c.field.IsArray {
		return c.createArrayValues()
	}
	return c.createSingleValue()
}

func (c *FieldValueCreator) validateRequired() error {
	switch c.field.Type {
	case model.FieldTypeString:
		strs := cast.ToStringSlice(c.value)
		if len(strs) == 0 || (len(strs) == 1 && strs[0] == "") {
			return fmt.Errorf("field %s is required", c.field.Name)
		}
	case model.FieldTypeInt:
		ints := cast.ToIntSlice(c.value)
		if len(ints) == 0 {
			return fmt.Errorf("field %s is required", c.field.Name)
		}
	case model.FieldTypeBool:
		// bool 可能为 false，所以只检查是否为 nil
		if c.value == nil {
			return fmt.Errorf("field %s is required", c.field.Name)
		}
	case model.FieldTypeDatetime:
		// 假设 value 是 time.Time 或 []time.Time
		if c.value == nil {
			return fmt.Errorf("field %s is required", c.field.Name)
		}
	}
	return nil
}

func (c *FieldValueCreator) createSingleValue() error {
	ifv := &model.ItemFieldValue{
		ItemID:  c.itemID,
		FieldID: c.field.ID,
	}

	switch c.field.Type {
	case model.FieldTypeString:
		s := cast.ToString(c.value)
		if s == "" {
			return fmt.Errorf("field %s is required", c.field.Name)
		}
		ifv.ValueString = &s
	case model.FieldTypeInt:
		i := cast.ToInt(c.value)
		ifv.ValueInt = &i
	case model.FieldTypeBool:
		b := cast.ToBool(c.value)
		ifv.ValueBool = &b
	case model.FieldTypeDatetime:
		t, ok := c.value.(time.Time)
		if !ok || t.IsZero() {
			return fmt.Errorf("invalid datetime value for field %s", c.field.Name)
		}
		ifv.ValueTime = &t
	default:
		return fmt.Errorf("unsupported field type: %d", c.field.Type)
	}

	return Create(c.tx, ifv)
}

func (c *FieldValueCreator) createArrayValues() error {
	var err error
	switch c.field.Type {
	case model.FieldTypeString:
		for _, s := range cast.ToStringSlice(c.value) {
			if s == "" {
				continue
			}
			err = c.createSingleValueWith(c.tx, &model.ItemFieldValue{
				ItemID:      c.itemID,
				FieldID:     c.field.ID,
				ValueString: &s,
			})
		}
	case model.FieldTypeInt:
		for _, i := range cast.ToIntSlice(c.value) {
			err = c.createSingleValueWith(c.tx, &model.ItemFieldValue{
				ItemID:   c.itemID,
				FieldID:  c.field.ID,
				ValueInt: &i,
			})
		}
	case model.FieldTypeBool:
		for _, b := range cast.ToBoolSlice(c.value) {
			err = c.createSingleValueWith(c.tx, &model.ItemFieldValue{
				ItemID:    c.itemID,
				FieldID:   c.field.ID,
				ValueBool: &b,
			})
		}
	case model.FieldTypeDatetime:
		for _, t := range c.value.([]time.Time) { // 注意类型断言
			if t.IsZero() {
				continue
			}
			err = c.createSingleValueWith(c.tx, &model.ItemFieldValue{
				ItemID:    c.itemID,
				FieldID:   c.field.ID,
				ValueTime: &t,
			})
		}
	default:
		return fmt.Errorf("unsupported array field type: %d", c.field.Type)
	}
	return err
}

// createSingleValueWith 用于数组场景，避免重复写 dao.Create
func (c *FieldValueCreator) createSingleValueWith(tx *gorm.DB, ifv *model.ItemFieldValue) error {
	return Create(tx, ifv)
}
