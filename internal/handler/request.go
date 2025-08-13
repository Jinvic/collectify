package handler

type DeletedReq struct {
	List []DeletedReqItem `json:"list" form:"list" binding:"required,dive"`
}

type SearchReq struct {
	NoPaging bool                   `json:"no_paging" form:"no_paging"` // 是否禁用分页
	Page     int                    `json:"page" form:"page"`
	PageSize int                    `json:"page_size" form:"page_size"`
	Filters  map[string]interface{} `json:"filters" form:"filters"`
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
