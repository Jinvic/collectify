package define

type ItemFieldValue struct {
	FieldID uint        `json:"field_id" form:"field_id" binding:"required,gt=0"`
	Value   interface{} `json:"value" form:"value" binding:"required"`
}
