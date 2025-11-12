package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
)

type Config struct {
	Prometheus  commonconfig.PrometheusConfig          `mapstructure:"prometheus"`
	Redis       commonconfig.RedisConfig               `mapstructure:"redis"`
	RabbitMq    commonconfig.RabbitMqConfig            `mapstructure:"rabbitmq"`
	Server      commonconfig.AnalyticsServiceConfig    `mapstructure:"analytics-service"`
	Clickhouse  commonconfig.AnalyticsClickHouseConfig `mapstructure:"analytics-clickhouse"`
	Mysql       commonconfig.LinkMysqlConfig           `mapstructure:"link-mysql"`
	ServiceInfo commonconfig.ServiceInfo               `mapstructure:"service-info"`
	Jaeger      commonconfig.JaegerConfig              `mapstructure:"jaeger"`
}

func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}
