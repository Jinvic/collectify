package handler

type DeletedReq struct {
	ID   uint   `json:"id" form:"id" binding:"required,gt=0"`
	Type string `json:"type" form:"type" binding:"required,oneof=category collection field item tag"`
}

type IDReq struct {
	ID uint `json:"id" form:"id" binding:"required,gt=0"`
}

type CreateCategoryReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type RenameCategoryReq struct {
	ID   uint   `json:"id" form:"id" binding:"required,gt=0"`
	Name string `json:"name" form:"name" binding:"required"`
}
