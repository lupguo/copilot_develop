package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// AppConfig 应用配置
type AppConfig struct {
	App struct {
		AuthToken       string `yaml:"auth_token"`
		SocksURL        string `yaml:"socks_url"`
		BlogStoragePath string `yaml:"blog_storage_path"`
	} `yaml:"app"`
}

// ParseAppConfig 解析配置文件
func ParseAppConfig(filename string) (*AppConfig, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "parse config filename: %s got err", filename)
	}

	// 解析成结构体
	cfg := AppConfig{}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal app config got err")
	}

	return &cfg, nil
}

// InitAppConfig 初始AppConfig配置信息
func InitAppConfig(filename string) error {
	cfg, err := ParseAppConfig(filename)
	if err != nil {
		return err
	}
	defaultAppConfig = cfg
	return nil
}

var defaultAppConfig *AppConfig

// GetAppConfig 获取AppConfig配置信息
func GetAppConfig() *AppConfig {
	return defaultAppConfig
}
