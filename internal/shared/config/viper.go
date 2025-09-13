package viper

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// TODO: 添加自动监测变化并重新读取功能

func NewViper(commonConfigName, commonConfigPath, serviceConfigName, serviceConfigPath string) (*viper.Viper, error) {
	v := viper.New()

	// FIXME: 改成更优雅的读取方式

	// 读取公共配置文件
	if commonConfigName != "" && commonConfigPath != "" {
		v.AddConfigPath(commonConfigPath)
		v.SetConfigName(commonConfigName)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("无法读取公共配置文件: %w", err)
			}
		}
	}

	// 读取专属配置文件
	if serviceConfigName != "" && serviceConfigPath != "" {
		v.AddConfigPath(serviceConfigPath)
		v.SetConfigName(serviceConfigName)
		if err := v.ReadInConfig(); err != nil {
			// 专属配置文件必须存在
			return nil, fmt.Errorf("无法读取专属配置文件: %w", err)
		}
	}

	// 启用环境变量,优先级最高
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
