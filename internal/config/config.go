package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

// 应用配置结构体
type Config struct {
	Database   ConfigDatabase   `env:",init"`
	Server     ConfigServer     `env:",init"`
	RecycleBin ConfigRecycleBin `env:",init"`
}

// 数据库配置
type ConfigDatabase struct {
	Type string `env:"DATABASE_TYPE" envDefault:"sqlite"`
	DSN  string `env:"DATABASE_DSN" envDefault:"collectify.db"`
}

// 服务器配置
type ConfigServer struct {
	Port int    `env:"SERVER_PORT" envDefault:"8080"`
	Mode string `env:"SERVER_MODE" envDefault:"release"`
}

// 回收站配置
type ConfigRecycleBin struct {
	Enable bool `env:"RECYCLEBIN_ENABLE" envDefault:"false"`
}

var config = &Config{}

func InitConfig() (*Config, error) {

	opt := env.Options{
		Prefix: "COLLECTIFY_",
	}

	cfg, err := env.ParseAsWithOptions[Config](opt)
	if err != nil {
		return nil, err
	}

	config = &cfg
	log.Println("Config loaded:", config)

	return config, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	return config
}

// SetConfig 设置配置，仅用于测试
func SetConfig(cfg *Config) {
	config = cfg
}
