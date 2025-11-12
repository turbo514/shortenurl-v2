package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
)

type Config struct {
	Server      commonconfig.TenantServiceConfig `mapstructure:"tenant-service"`
	Mysql       commonconfig.TenantMysqlConfig   `mapstructure:"tenant-mysql"`
	ServiceInfo commonconfig.ServiceInfo         `mapstructure:"service-info"`
	Jaeger      commonconfig.JaegerConfig        `mapstructure:"jaeger"`
	Prometheus  commonconfig.PrometheusConfig    `mapstructure:"prometheus"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}
