package e

import (
	"errors"
	"fmt"
)

type EInterface interface {
	error

	Error() string
	Wrap(err error) EInterface
}

type EStruct struct {
	err error
}

func (e EStruct) Error() string {
	return e.err.Error()
}
func (e EStruct) Wrap(err error) EInterface {
	return EStruct{
		err: fmt.Errorf("%s: %w", e.Error(), err),
	}
}

var (
	ErrInvalidParams = EStruct{
		err: errors.New("参数错误"),
	}
	ErrInternal = EStruct{
		err: errors.New("内部错误"),
	}

	ErrDuplicated = EStruct{
		err: errors.New("目标已存在"),
	}
	ErrNotFound = EStruct{
		err: errors.New("目标不存在"),
	}

	ErrUnauthorized = EStruct{
		err: errors.New("未授权"),
	}
	ErrUserNotFound = EStruct{
		err: errors.New("用户不存在"),
	}
	ErrUserInvalidPassword = EStruct{
		err: errors.New("用户名或密码错误"),
	}
)
