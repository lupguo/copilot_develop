package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

const (
	PromptKeySummaryBlog = "summary-blog"
)

// AppPromptConfig 提示词配置
type AppPromptConfig struct {
	AppPrompts []Prompt `yaml:"app_prompt"`
}

// Prompt 提示词
type Prompt struct {
	Name              string                         `yaml:"name"`
	AIMode            string                         `yaml:"ai_mode"`
	MaxTokens         int                            `yaml:"max_tokens"`
	PredefinedPrompts []openai.ChatCompletionMessage `yaml:"predefined_prompts"`
}

// AppPromptMap name到提示词的映射
type AppPromptMap map[string]*Prompt

// ParseAppPromptConfig 解析App提示词配置文件
func ParseAppPromptConfig(filename string) (AppPromptMap, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read prompt yaml config got err")
	}

	// 解析prompt app config
	cfg := AppPromptConfig{}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "parse app prompt yaml config got err")
	}

	// 转成map
	m := AppPromptMap{}
	for _, prompt := range cfg.AppPrompts {
		m[prompt.Name] = &prompt
	}

	return m, nil
}
