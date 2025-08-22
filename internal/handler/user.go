package handler

import (
	"collectify/internal/config"
	"collectify/internal/conn"
	"collectify/internal/dao"
	model "collectify/internal/model/db"
	"collectify/internal/model/define"
	e "collectify/internal/pkg/e"
	"collectify/internal/service"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 检查认证是否启用
func IsAuthEnabled(c *gin.Context) {
	cfg := config.GetConfig()
	SuccessWithData(c, cfg.Auth.Enable)
}

// 登录
func UserLogin(c *gin.Context) {
	cfg := config.GetConfig()

	// 如果未启用认证，则直接返回成功
	if !cfg.Auth.Enable {
		Success(c)
		return
	}

	var req define.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, e.ErrInvalidParams)
		return
	}

	db := conn.GetDB()
	var user model.User
	err := db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Fail(c, e.ErrUserNotFound)
			return
		}
		Fail(c, err)
		return
	}

	if user.Password != req.Password {
		Fail(c, e.ErrUserInvalidPassword)
		return
	}

	expireTime := time.Duration(cfg.Auth.ExpireDay) * time.Hour * 24
	token, err := service.GenerateToken(user, expireTime)
	if err != nil {
		Fail(c, err)
		return
	}

	SuccessWithData(c, define.LoginResp{
		Token:    token,
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
}

func UserUpdate(c *gin.Context) {
	cfg := config.GetConfig()

	// 如果未启用认证，则直接返回成功
	if !cfg.Auth.Enable {
		Success(c)
		return
	}

	var req define.UpdateUserReq
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, e.ErrInvalidParams)
		return
	}

	// 检查是否重复
	uniqueFields := map[string]interface{}{"username": req.Username}
	filters := []dao.Filter{}
	id, isDeleted, err := dao.DuplicateCheck[model.User](conn.GetDB(), uniqueFields, filters)
	if err != nil {
		Fail(c, err)
		return
	}
	if id != 0 {
		FailWithData(c, e.ErrDuplicated, map[string]interface{}{
			"id":        id,
			"isDeleted": isDeleted,
		})
		return
	}

	uniqueFields = map[string]interface{}{"id": req.ID}
	updateFields := map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
	}

	err = dao.Update[model.User](conn.GetDB(), uniqueFields, updateFields)
	if err != nil {
		Fail(c, err)
		return
	}

	Success(c)
}
