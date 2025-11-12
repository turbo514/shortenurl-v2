package config

import (
	"fmt"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"

	"github.com/spf13/viper"
)

type Config struct {
	Server            commonconfig.ApiGatewayConfig       `mapstructure:"api-gateway"`
	ServiceInfo       commonconfig.ServiceInfo            `mapstructure:"service-info"`
	LinkService       commonconfig.LinkServiceConfig      `mapstructure:"link-service"`
	TenantService     commonconfig.TenantServiceConfig    `mapstructure:"tenant-service"`
	AnalyticsService  commonconfig.AnalyticsServiceConfig `mapstructure:"analytics-service"`
	Jaeger            commonconfig.JaegerConfig           `mapstructure:"jaeger"`
	Redis             commonconfig.RedisConfig            `mapstructure:"redis"`
	GlobalRateLimiter commonconfig.RateLimiterConfig      `mapstructure:"global-rate-limiter"`
	LocalRateLimiter  commonconfig.RateLimiterConfig      `mapstructure:"local-rate-limiter"`
}

// 直接依赖viper,可能不是个好做法(?)
func NewConfig(v *viper.Viper) (*Config, error) {
	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("无法将配置反序列化到结构体中: %w", err)
	}
	return cfg, nil
}
