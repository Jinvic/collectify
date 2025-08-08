package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Database ConfigDatabase `mapstructure:"database"`
	Server   ConfigServer   `mapstructure:"server"`
}

type ConfigDatabase struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}

type ConfigServer struct {
	Port int `mapstructure:"port"`
}

var defaultConfig = Config{
	Database: ConfigDatabase{
		Type: "sqlite",
		DSN:  "collectify.db",
	},
	Server: ConfigServer{
		Port: 8080,
	},
}

var config *Config

func InitConfig() error {
	viper.SetConfigName("config")            // 读取名为config的配置文件，没有设置特定的文件后缀名
	viper.SetConfigType("yaml")              // 当没有设置特定的文件后缀名时，必须要指定文件类型
	viper.AddConfigPath("./")                // 在当前文件夹下寻找
	viper.AddConfigPath("$HOME/.collectify") // 在用户目录下寻找

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在或无法读取，则加载默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 使用默认配置
			setDefaultFromStruct(defaultConfig)
		} else {
			return fmt.Errorf("fatal error config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}

func GetConfig() *Config {
	if config == nil {
		panic("config not loaded, call InitConfig() first")
	}
	return config
}

// setDefaultFromStruct 将结构体中的默认值设置到 Viper 中
func setDefaultFromStruct(cfg interface{}) {
	val := reflect.ValueOf(cfg).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := typ.Field(i)
		tagValue := typeField.Tag.Get("mapstructure")

		if valueField.Kind() == reflect.Struct {
			setDefaultFromStruct(valueField.Addr().Interface())
		} else {
			viper.SetDefault(tagValue, valueField.Interface())
		}
	}
}
