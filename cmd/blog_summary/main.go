package main

import (
	"context"
	"time"

	"github.com/lupguo/copilot_develop/app/application"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/lupguo/copilot_develop/app/infras/dbs"
	"github.com/lupguo/copilot_develop/app/infras/openaix"
	"github.com/lupguo/copilot_develop/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	configFile string // 应用配置文件
	blogPath   string // blog路径
)

func init() {
	pflag.StringVar(&configFile, "conf", "./config.yaml", "Path to the app YAML config file")
	pflag.StringVar(&blogPath, "blog_path", "/private/data/www/tkstorm.com/content/", "The path of Blog content AI Summary")
}

// Blog总结基本流程
// 1. 获取指定目录的所有文件内容，返回文件的绝对路径集合
// 2. 并行化读取文件内容，通过OpenAI提取文件内容摘要、关键字信息，对原MD进行替换
func main() {
	pflag.Parse()
	// client config
	if err := config.ParseConfig(configFile); err != nil {
		log.Fatalf("parse config got err: %s", err)
	}

	start := time.Now()
	app, err := buildBlogSummaryApp()
	if err != nil {
		log.Fatalf("init blog summary got err: %s", err)
	}

	err = app.UpdateBlogHeaderYaml(context.Background(), blogPath)
	if err != nil {
		log.Fatalf("update blog summary content got err: %s", err)
	}

	log.Infof("update blog summary using time: %s", time.Since(start))
}

func buildBlogSummaryApp() (*application.BlogSummaryApp, error) {
	// sqlite infra
	sqliteDbInfra, err := dbs.NewBlogSummarySqliteInfra(config.GetDBFilePath())
	if err != nil {
		return nil, errors.Wrap(err, "new sqlite infra got err")
	}

	// openAI infra
	openAIInfra, err := openaix.NewOpenAIHttpProxyClient()
	if err != nil {
		return nil, errors.Wrap(err, "new open ai http proxy client got err")
	}

	// blog summary app
	blogSummaryApp := application.NewBlogSummaryApp(
		service.NewAIService(openAIInfra),
		sqliteDbInfra,
	)
	return blogSummaryApp, nil
}
