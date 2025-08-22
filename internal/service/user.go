package service

import (
	"collectify/internal/config"
	model "collectify/internal/model/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成 token
func GenerateToken(user model.User, expire time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	}

	// 如果过期时间大于 0，则设置过期时间
	if expire > 0 {
		claims["exp"] = time.Now().Add(expire).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().Auth.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
