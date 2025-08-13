package handler

type DeletedReq struct {
	List []DeletedReqItem `json:"list" form:"list" binding:"required,dive"`
}

type DeletedReqItem struct {
	ID   uint   `json:"id" form:"id" binding:"required,gt=0"`
	Type string `json:"type" form:"type" binding:"required,oneof=category collection field item tag"`
}

type CreateCategoryReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type RenameCategoryReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}
