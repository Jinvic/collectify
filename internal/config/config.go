package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

// 应用配置结构体
type Config struct {
	Database   ConfigDatabase   `env:",init"`
	Server     ConfigServer     `env:",init"`
	RecycleBin ConfigRecycleBin `env:",init"`
	Auth       ConfigAuth       `env:",init"`
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

type ConfigAuth struct {
	Enable    bool   `env:"AUTH_ENABLE" envDefault:"false"`
	JwtSecret string `env:"AUTH_JWT_SECRET"`
	ExpireDay int    `env:"AUTH_EXPIRE_DAY" envDefault:"15"`
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

	// 如果启用认证，且未设置 JWT 密钥
	if config.Auth.Enable && config.Auth.JwtSecret == "" {
		if config.Server.Mode == "debug" { // 开发模式下，生成一个随机密钥
			secret, err := GenerateJWTSecret()
			if err != nil {
				return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
			}
			config.Auth.JwtSecret = secret
		} else { // 生产模式下，必须设置 JWT 密钥
			return nil, errors.New("JWT_SECRET is required")
		}
	}

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

// GenerateJWTSecret 生成 32 字节（256 位）的随机密钥，并以 Base64 编码返回
func GenerateJWTSecret() (string, error) {
	secret := make([]byte, 32)

	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(secret)

	return encoded, nil
}
