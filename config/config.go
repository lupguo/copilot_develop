package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type OpenAIProxyConfig struct {
	AuthToken string `yaml:"auth_token"`
	SocksURL  string `yaml:"socks_url"`
}

type BlogSummaryConfig struct {
	AIPromptFile string `yaml:"ai_prompt_file"` // blog summary prompt配置
	SQLiteDBFile string `yaml:"sqlite_db_file"` // blog sqlite db存储
}

// Config 应用配置
type Config struct {
	App *AppConfig `yaml:"app"`
}

type AppConfig struct {
	RootPath    string             `yaml:"root_path"` // 根目录
	OpenAIProxy *OpenAIProxyConfig `yaml:"openai_proxy"`
	BlogSummary *BlogSummaryConfig `yaml:"blog_summary"`
}

var (
	defaultConfig *Config
	appConfig     *AppConfig
)

// ParseConfig 解析配置文件
func ParseConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "parse config filename: %s got err", filename)
	}

	// 解析成结构体
	if err = yaml.Unmarshal(file, &defaultConfig); err != nil {
		return errors.Wrap(err, "yaml unmarshal app config got err")
	}

	// 解析检测
	appConfig = defaultConfig.App
	switch {
	case appConfig.RootPath == "":
		return errors.New("empty root_path config")
	case appConfig.OpenAIProxy == nil:
		return errors.New("empty openai_proxy config")
	case appConfig.BlogSummary == nil:
		return errors.New("empty blog_summary config")
	}

	// // prompt parse
	// if err = ParseAppPromptConfig(GetPromptConfigPath()); err != nil {
	// 	return errors.Wrapf(err, "parse app prompt config got err")
	// }

	return nil
}

// GetPromptConfigPath 获取配置文件路径
func GetPromptConfigPath() string {
	return filepath.Join(appConfig.RootPath, appConfig.BlogSummary.AIPromptFile)
}

// GetDBFilePath 获取数据存储路径
func GetDBFilePath() string {
	return filepath.Join(appConfig.RootPath, appConfig.BlogSummary.SQLiteDBFile)
}

// GetOpenAIProxy 底层OpenAI Http Proxy配置
func GetOpenAIProxy() *OpenAIProxyConfig {
	return appConfig.OpenAIProxy
}
