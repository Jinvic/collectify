package middleware

import (
	"collectify/internal/config"
	"collectify/internal/handler"
	"collectify/internal/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthCheck(c *gin.Context) {
	// 如果未启用认证，则直接跳过
	if !config.GetConfig().Auth.Enable {
		c.Next()
		return
	}

	// 获取 token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		handler.Fail(c, e.ErrUnauthorized)
		return
	}

	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Auth.JwtSecret), nil
	})
	if err != nil {
		handler.Fail(c, e.ErrUnauthorized)
		return
	}

	// 验证 token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["id"])
		c.Set("user_username", claims["username"])
		c.Set("user_role", claims["role"])
		c.Next()
		return
	}

	handler.Fail(c, e.ErrUnauthorized)
}
