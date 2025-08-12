package e

import "errors"

var (
	ErrInvalidParams = errors.New("参数错误")
	ErrInternal      = errors.New("内部错误")

	ErrDuplicated = errors.New("目标已存在")
	ErrNotFound   = errors.New("目标不存在")
)
