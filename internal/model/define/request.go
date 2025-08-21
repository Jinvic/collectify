package define

import "time"

type DeletedReq struct {
	List []DeletedReqItem `json:"list" form:"list" binding:"required,dive"`
}

type ListReq struct {
	NoPaging bool `json:"no_paging" form:"no_paging"` // 是否禁用分页
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"page_size" form:"page_size"`
}
type SearchReq struct {
	ListReq
	Filters map[uint]interface{} `json:"filters" form:"filters"`
}

type SearchFilterDatetime struct {
	Start time.Time `json:"start" form:"start"`
	End   time.Time `json:"end" form:"end"`
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

type CreateFieldReq struct {
	CategoryID uint   `json:"category_id" form:"category_id" binding:"required,gt=0"`
	Name       string `json:"name" form:"name" binding:"required"`
	Type       int    `json:"type" form:"type" binding:"required,oneof=1 2 3 4"`
	IsArray    bool   `json:"is_array" form:"is_array"`
	Required   bool   `json:"required" form:"required"`
}

type CreateItemReq struct {
	CategoryID uint `json:"category_id" form:"category_id" binding:"required,gt=0"`
	Item       Item `json:"item" form:"item" binding:"required"`
}

type UpdateItemReq struct {
	ID   uint `json:"id" form:"id" binding:"required,gt=0"`
	Item Item `json:"item" form:"item" binding:"required"`
}

type SearchItemsReq struct {
	SearchReq
	CategoryID    uint   `json:"category_id" form:"category_id" binding:"required,gt=0"`
	Name          string `json:"name" form:"name"`
	TagIDs        []uint `json:"tag_ids" form:"tag_ids"`
	CollectionIDs []uint `json:"collection_ids" form:"collection_ids"`
}
