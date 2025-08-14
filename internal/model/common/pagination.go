package common

type Pagination struct {
	Disable bool `json:"disable" form:"disable"` // 是否禁用分页
	Page    int  `json:"page" form:"page"`       // 当前页码
	Size    int  `json:"size" form:"size"`       // 每页数量
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
