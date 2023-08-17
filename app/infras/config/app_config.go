package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// AppConfig 应用配置
type AppConfig struct {
	App struct {
		AuthToken string `yaml:"auth_token"`
		SocksURL  string `yaml:"socks_url"`
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

// GetAppRoot 获取项目根目录
func GetAppRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	index := strings.Index(dir, "app")

	if index > 0 {
		return dir[:index]
	}
	return dir
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	return filepath.Join(GetAppRoot(), "conf")
}

// GetDataPath 获取数据存储路径
func GetDataPath() string {
	return filepath.Join(GetAppRoot(), "data")
}
