package model

type Pagination struct {
	Page int `json:"page" form:"page"` // 当前页码
	Size int `json:"size" form:"size"` // 每页数量
}

func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Size < 1 || p.Size > 100 {
		p.Size = 20 // 默认每页 20 条
	}
	return p.Size
}
