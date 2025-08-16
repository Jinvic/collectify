package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// 应用配置结构体
type Config struct {
	Database   ConfigDatabase   `mapstructure:"database"`
	Server     ConfigServer     `mapstructure:"server"`
	RecycleBin ConfigRecycleBin `mapstructure:"recyclebin"`
}

// 数据库配置
type ConfigDatabase struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}

// 服务器配置
type ConfigServer struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// 回收站配置
type ConfigRecycleBin struct {
	Enable bool `mapstructure:"enable"`
}

// 默认配置
var defaultConfig = Config{
	Database: ConfigDatabase{
		Type: "sqlite",
		DSN:  "./data/collectify.db", // 默认数据库路径
	},
	Server: ConfigServer{
		Port: 8080,
		Mode: "release", // 默认运行模式
	},
	RecycleBin: ConfigRecycleBin{
		Enable: false, // 默认禁用回收站
	},
}

var config = &Config{}

func InitConfig() (*Config, error) {
	// 1. 手动加载 .env 文件到环境变量
	if err := godotenv.Load(); err != nil {
		// .env 不存在也 OK
		fmt.Println("No .env file found, proceeding with environment variables and defaults.")
	}

	// 2. 设置默认值（带嵌套路径）
	setDefaultFromStruct(&defaultConfig)

	// 创建 data 目录用于持久化存储
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// 3. 启用环境变量读取
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COLLECTIFY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	fmt.Printf("Loaded config: %+v\n", config)
	return config, nil
}

// 获取配置，未初始化时panic
func GetConfig() *Config {
	if config == nil {
		panic("config not loaded, call InitConfig() first")
	}
	return config
}

// 设置配置（仅用于测试）
func SetConfig(cfg *Config) {
	config = cfg
}

// 递归设置默认配置
func setDefaultFromStruct(cfg interface{}, parts ...string) {
	val := reflect.ValueOf(cfg)
	typ := reflect.TypeOf(cfg)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("cfg must be a struct or pointer to struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		tag := structField.Tag.Get("mapstructure")

		if tag == "" || tag == "-" {
			continue
		}

		keyPath := append(parts, tag)
		fullKey := strings.Join(keyPath, ".")

		if field.Kind() == reflect.Struct {
			setDefaultFromStruct(field.Addr().Interface(), keyPath...)
		} else if field.CanInterface() {
			viper.SetDefault(fullKey, field.Interface())
		}
	}
}
