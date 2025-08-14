package define

import model "collectify/internal/model/db"

type ItemFieldValue struct {
	FieldID uint        `json:"field_id" form:"field_id" binding:"required,gt=0"`
	Value   interface{} `json:"value" form:"value" binding:"required"`
}

type Item struct {
	Title       string           `json:"title" form:"title" binding:"required"`
	Status      int              `json:"status" form:"status" binding:"required,oneof=1 2 3 4 5"`
	Rating      *float64         `json:"rating" form:"rating" binding:"omitempty,min=0,max=10"`
	Description string           `json:"description" form:"description"`
	Notes       string           `json:"notes" form:"notes"`
	CoverURL    string           `json:"cover_url" form:"cover_url" binding:"omitempty,url"`
	SourceURL   string           `json:"source_url" form:"source_url" binding:"omitempty,url"`
	Priority    int              `json:"priority" form:"priority" binding:"omitempty,min=0"`
	Values      []ItemFieldValue `json:"values" form:"values" binding:"omitempty,dive"`
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
