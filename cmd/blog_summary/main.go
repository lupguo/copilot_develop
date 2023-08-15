package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/hold7techs/go-shim/log"
	"github.com/lupguo/copilot_develop/app/application"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/lupguo/copilot_develop/app/infras/config"
	"github.com/lupguo/copilot_develop/app/infras/openaix"
)

// Blog总结基本流程
// 1. 获取指定目录的所有文件内容，返回文件的绝对路径集合
// 2. 并行化读取文件内容，通过OpenAI提取文件内容摘要、关键字信息，对原MD进行替换
func main() {
	var (
		blogStoragePath     = `/private/data/www/tkstorm.com/content/posts/application/ai/chatgpt/cloudflare-warp.md`
		appYamlConfigFile   = filepath.Join(config.GetConfigPath(), "app_dev.yaml")
		appPromptConfigFile = filepath.Join(config.GetConfigPath(), "prompt.yaml")
	)

	start := time.Now()
	app, err := initializeBlogSummaryApp(appYamlConfigFile, appPromptConfigFile)
	if err != nil {
		log.Fatalf("init blog summary got err: %s", err)
	}

	err = app.UpdateBlogSummaryContent(context.Background(), blogStoragePath)
	if err != nil {
		log.Fatalf("update blog summary content got err: %s", err)
	}

	log.Debugf("update blog summary using time: %s", time.Since(start))
}

func initializeBlogSummaryApp(appYamlConfigFile, appPromptConfigFile string) (*application.BlogSummaryApp, error) {
	// client config
	appCfg, err := config.ParseAppConfig(appYamlConfigFile)
	if err != nil {
		return nil, err
	}

	// app prompt
	appPromptCfg, err := config.ParseAppPromptConfig(appPromptConfigFile)
	if err != nil {
		return nil, err
	}

	// openAI infra
	openAIInfra, err := openaix.NewOpenAIHttpProxyClient(appCfg)
	if err != nil {
		return nil, err
	}

	// service
	aiService := service.NewAIService(openAIInfra, appPromptCfg)

	// app
	blogSummaryApp := application.NewBlogSummaryApp(aiService)

	return blogSummaryApp, nil
}
