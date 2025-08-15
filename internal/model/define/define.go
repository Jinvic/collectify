package define

import (
	model "collectify/internal/model/db"
	"time"
)

type ItemFieldValue struct {
	FieldID uint        `json:"field_id" form:"field_id" binding:"required,gt=0"`
	Value   interface{} `json:"value" form:"value" binding:"required"`

	FieldName string `json:"field_name"`
	FieldType int    `json:"field_type"`
}

type Item struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	Title       string           `json:"title" form:"title" binding:"required"`
	Status      int              `json:"status" form:"status" binding:"required,oneof=1 2 3 4 5"`
	Rating      *float64         `json:"rating" form:"rating" binding:"omitempty,min=0,max=10"`
	Description string           `json:"description" form:"description"`
	Notes       string           `json:"notes" form:"notes"`
	CoverURL    string           `json:"cover_url" form:"cover_url" binding:"omitempty,url"`
	SourceURL   string           `json:"source_url" form:"source_url" binding:"omitempty,url"`
	Priority    int              `json:"priority" form:"priority" binding:"omitempty,min=0"`
	Values      []ItemFieldValue `json:"values" form:"values" binding:"omitempty,dive"`

	Category Category `json:"category"`
}

func (i Item) ToDB() *model.Item {
	return &model.Item{
		Title:       i.Title,
		Status:      i.Status,
		Rating:      i.Rating,
		Description: i.Description,
		Notes:       i.Notes,
		CoverURL:    i.CoverURL,
		SourceURL:   i.SourceURL,
		Priority:    i.Priority,
	}
}

func (i *Item) FromDB(item *model.Item) {
	i.ID = item.ID
	i.CreatedAt = item.CreatedAt
	i.UpdatedAt = item.UpdatedAt
	i.DeletedAt = item.DeletedAt.Time

	i.Title = item.Title
	i.Status = item.Status
	i.Rating = item.Rating
	i.Description = item.Description
	i.Notes = item.Notes
	i.CoverURL = item.CoverURL
	i.SourceURL = item.SourceURL
	i.Priority = item.Priority

	i.Category.FromDB(&item.Category)
}

type ItemDetail struct {
	Item
	Tags        []Tag        `json:"tags"`
	Collections []Collection `json:"collections"`
}

func (i *ItemDetail) FromDB(item *model.Item) {
	i.Item.FromDB(item)

	i.Tags = make([]Tag, len(item.Tags))
	i.Collections = make([]Collection, len(item.Collections))
	i.Values = make([]ItemFieldValue, len(item.Values))

	for idx, tag := range item.Tags {
		i.Tags[idx].FromDB(&tag)
	}

	for idx, collection := range item.Collections {
		i.Collections[idx].FromDB(&collection)
	}

	var fieldValueMap = make(map[uint][]model.ItemFieldValue)
	var fieldIsArray = make(map[uint]bool)
	for _, value := range item.Values {
		if _, ok := fieldValueMap[value.FieldID]; !ok {
			fieldValueMap[value.FieldID] = []model.ItemFieldValue{}
		}
		fieldValueMap[value.FieldID] = append(fieldValueMap[value.FieldID], value)
		fieldIsArray[value.FieldID] = value.Field.IsArray
	}

	getFieldValue := func(value model.ItemFieldValue) interface{} {
		switch value.Field.Type {
		case model.FieldTypeString:
			return value.ValueString
		case model.FieldTypeInt:
			return value.ValueInt
		case model.FieldTypeBool:
			return value.ValueBool
		case model.FieldTypeDatetime:
			return value.ValueTime
		default:
			return nil
		}
	}

	for fieldID, values := range fieldValueMap {
		ifv := ItemFieldValue{
			FieldID:   fieldID,
			FieldName: values[0].Field.Name,
			FieldType: values[0].Field.Type,
		}
		if fieldIsArray[fieldID] {
			var valueArray []interface{}
			for _, value := range values {
				valueArray = append(valueArray, getFieldValue(value))
			}
			ifv.Value = valueArray
		} else {
			ifv.Value = getFieldValue(values[0])
		}
		i.Values = append(i.Values, ifv)
	}
}

type Field struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
	IsArray bool `json:"is_array"`
	Required bool `json:"required"`
}

func (f *Field) FromDB(field *model.Field) {
	f.ID = field.ID
	f.Name = field.Name
	f.Type = field.Type
	f.IsArray = field.IsArray
	f.Required = field.Required
}

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	Fields []Field `json:"fields"`
}

func (c *Category) FromDB(category *model.Category) {
	c.ID = category.ID
	c.Name = category.Name

	c.Fields = make([]Field, len(category.Fields))
	for idx, field := range category.Fields {
		c.Fields[idx].FromDB(&field)
	}
}

type Tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (t *Tag) FromDB(tag *model.Tag) {
	t.ID = tag.ID
	t.Name = tag.Name
}

type Collection struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (c *Collection) FromDB(collection *model.Collection) {
	c.ID = collection.ID
	c.Name = collection.Name
}
