package openaix

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
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
	PredefinedPrompts []openai.ChatCompletionMessage `yaml:"predefined_prompts"` // 预先定义的提示内容（例如定义AI角色）
}

var defaultPromptSetting map[string]*Prompt

// ParseAppPromptConfig 解析App提示词配置文件
func ParseAppPromptConfig(filename string) (error, map[string]*Prompt) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "read prompt yaml config got err"), nil
	}

	// 解析prompt app config
	cfg := AppPromptConfig{}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return errors.Wrap(err, "parse app prompt yaml config got err"), nil
	}

	// 转成map
	defaultPromptSetting = make(map[string]*Prompt)
	for _, prompt := range cfg.AppPrompts {
		defaultPromptSetting[prompt.Name] = &prompt
	}

	return nil, defaultPromptSetting
}

// GetPrompt 获取指定的Prompt配置信息
func GetPrompt(key string) (*Prompt, error) {
	if v, ok := defaultPromptSetting[key]; ok {
		return v, nil
	}
	return nil, errors.Errorf("prompt[%s] not found", key)
}
