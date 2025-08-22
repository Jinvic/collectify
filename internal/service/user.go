package service

import (
	"collectify/internal/config"
	model "collectify/internal/model/db"

	"github.com/golang-jwt/jwt/v5"
)

func Login(user model.User) (string, error) {

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		// "exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // 30天过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().Auth.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
