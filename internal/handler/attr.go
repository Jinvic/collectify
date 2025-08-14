package handler

import (
	"collectify/internal/model/common"
	"collectify/internal/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 从路径中获取ID并验证
func GetID(c *gin.Context, idKey string) (uint, error) {
	idStr := c.Param(idKey)
	id := cast.ToUint(idStr)
	if id <= 0 {
		return 0, e.ErrInvalidParams
	}
	return id, nil
}

// 从查询参数中获取分页参数
func GetPagination(c *gin.Context) (common.Pagination, error) {
	page := cast.ToInt(c.Query("page"))
	pageSize := cast.ToInt(c.Query("page_size"))
	noPaging := cast.ToBool(c.Query("no_paging"))

	if page <= 0 || pageSize <= 0 {
		return common.Pagination{}, e.ErrInvalidParams
	}

	return common.Pagination{
		Page:    page,
		Size:    pageSize,
		Disable: noPaging,
	}, nil
}
