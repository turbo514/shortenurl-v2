package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
)

type Config struct {
	ServiceInfo commonconfig.ServiceInfo       `mapstructure:"service-info"`
	Mysql       commonconfig.LinkMysqlConfig   `mapstructure:"link-mysql"`
	Jaeger      commonconfig.JaegerConfig      `mapstructure:"jaeger"`
	Redis       commonconfig.RedisConfig       `mapstructure:"redis"`
	Kafka       commonconfig.KafkaConfig       `mapstructure:"kafka"`
	Prometheus  commonconfig.PrometheusConfig  `mapstructure:"prometheus"`
	Server      commonconfig.LinkServiceConfig `mapstructure:"link-service"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}
