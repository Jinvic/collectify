package conn

import (
	"collectify/internal/config"
	model "collectify/internal/model/db"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(cfg *config.Config) error {

	// 暂时只支持sqlite
	switch cfg.Database.Type {
	case "sqlite":
		var err error
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if cfg.Server.Mode == "debug" {
		db = db.Debug()
	}

	err := db.AutoMigrate(
		&model.Category{},
		&model.Collection{},
		&model.Field{},
		&model.Item{},
		&model.Tag{},
		&model.ItemFieldValue{},
		&model.User{},
	)
	if err != nil {
		return err
	}

	err = initAdminUser()
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}

// 当用户表为空时，创建一个管理员用户
func initAdminUser() error {
	var count int64
	err := db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	user := &model.User{
		Username: "admin",
		Password: "admin",
		Role:     model.UserRoleAdmin,
	}

	return db.Create(user).Error
}
