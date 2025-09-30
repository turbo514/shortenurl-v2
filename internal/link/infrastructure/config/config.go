package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
)

type Config struct {
	Common         commonconfig.CommonConfig `mapstructure:"common"`
	Server         ServerConfig              `mapstructure:"server"`
	DatabaseConfig DatabaseConfig            `mapstructure:"database"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	Options         string `mapstructure:"options"`
	MigrateFilePath string `mapstructure:"migrate_file_path"`
}
