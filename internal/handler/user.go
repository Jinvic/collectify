package handler

import (
	"collectify/internal/config"
	"collectify/internal/conn"
	model "collectify/internal/model/db"
	"collectify/internal/model/define"
	e "collectify/internal/pkg/e"
	"collectify/internal/service"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
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
