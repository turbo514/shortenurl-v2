package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Services struct {
		Link          string `mapstructure:"link"`
		LinkPort      int    `mapstructure:"link_port"`
		Tenant        string `mapstructure:"tenant"`
		TenantPort    int    `mapstructure:"tenant_port"`
		Analytics     string `mapstructure:"analytics"`
		AnalyticsPort int    `mapstructure:"analytics_port"`
	} `mapstructure:"services"`
}

// 直接依赖viper,可能不是个好做法(?)
func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}
