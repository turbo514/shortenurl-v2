package commonconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// TODO: 添加自动监测变化并重新读取功能

func NewViper(commonConfigName, commonConfigPath, serviceConfigName, serviceConfigPath string) (*viper.Viper, error) {
	// 1. 创建主Viper实例，用于最终配置和公共配置
	v := viper.New()

	// --- 读取公共配置文件 ---
	if commonConfigName != "" && commonConfigPath != "" {
		// 为公共配置设置路径和名称
		v.AddConfigPath(commonConfigPath)
		v.SetConfigName(commonConfigName)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("无法读取公共配置文件: %w", err)
		}
	}

	// --- 读取专属配置文件，并合并 ---
	if serviceConfigName != "" && serviceConfigPath != "" {
		// **关键步骤：创建一个临时的Viper实例来读取专属配置**
		// 这样它就不会干扰主v实例已加载的配置。
		serviceViper := viper.New()
		serviceViper.AddConfigPath(serviceConfigPath)
		serviceViper.SetConfigName(serviceConfigName)

		if err := serviceViper.ReadInConfig(); err != nil {
			// 专属配置文件必须存在
			return nil, fmt.Errorf("无法读取专属配置文件: %w", err)
		}

		// **将专属配置的内容合并到主配置v中**
		// MergeConfigMap 是进行配置分层合并的正确方法。
		if err := v.MergeConfigMap(serviceViper.AllSettings()); err != nil {
			return nil, fmt.Errorf("无法合并专属配置文件: %w", err)
		}
	}

	// 启用环境变量,优先级最高
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}
