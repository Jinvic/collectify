package model

import "gorm.io/gorm"

const (
	UserRoleAdmin = iota + 1 // 管理员
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Role     int    `gorm:"not null;default:1" json:"role"`
}
