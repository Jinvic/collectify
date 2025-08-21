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
	)
	return err
}

func GetDB() *gorm.DB {
	return db
}
