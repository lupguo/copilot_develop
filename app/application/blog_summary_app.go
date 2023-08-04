package application

import (
	"bytes"
	"context"
	"os"

	"github.com/hold7techs/go-shim/log"
	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/lupguo/copilot_develop/app/infras/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// BlogSummaryApp Blog的汇总App
type BlogSummaryApp struct {
	srv service.IServicesSummaryAI
}

// UpdateBlogSummaryContent 更新Blog的汇总信息
func (app *BlogSummaryApp) UpdateBlogSummaryContent() error {
	// 查询目录下所有的markdown目录 -> slice内 []*BlogMD
	var blogStorageRoot string
	if config.GetAppConfig() != nil {
		blogStorageRoot = config.GetAppConfig().App.BlogStoragePath
	}
	blogFilePaths, err := shim.FindFilePaths(blogStorageRoot, "*.md")
	if err != nil {
		return errors.Wrapf(err, "shim find file paths root [%s] got err", blogStorageRoot)
	}

	// 通过正则提取md的主题内容
	ctx := context.Background()
	for _, blogFilePath := range blogFilePaths {
		// 替换每个md的汇总信息、关键字、描述信息
		if err := app.ReplaceKeywordsAndSummary(ctx, blogFilePath); err != nil {
			return err
		}
	}

	return nil
}

// ReplaceKeywordsAndSummary 将keywords, summary 填充到原有的blog文章内
func (app *BlogSummaryApp) ReplaceKeywordsAndSummary(ctx context.Context, blogFilePath string) error {
	// 初始每个md
	md, err := entity.NewBlogMD(blogFilePath)
	if err != nil {
		return errors.Wrap(err, "entity new blog md got err in replace")
	}

	// cont 提前摘要
	summary, err := app.srv.BlogSummary(ctx, md.MDContent)
	if err != nil {
		return errors.Wrap(err, "using open ai extract summary for blogContent got err")
	}

	// 汇总、关键字、描述
	md.MDHeader.Summary = summary.Summary
	md.MDHeader.Keywords = summary.Keywords
	md.MDHeader.Description = summary.Description

	// 将调整后的md更新回去
	newMDHeaderYamlStr, err := yaml.Marshal(md.MDHeader)
	if err != nil {
		return errors.Wrapf(err, "marsh file[%s] yaml head got err", blogFilePath)
	}
	log.Debugf("newMDHeaderYamlStr: %s", newMDHeaderYamlStr)

	// 清空原文件
	mdFile, err := os.Open(md.Filepath)
	if err != nil {
		return errors.Wrapf(err, "open blog file[%s] got err", md.Filepath)
	}
	err = mdFile.Truncate(0)
	if err != nil {
		return errors.Wrapf(err, "truncate blog file[%s] content got err", md.Filepath)
	}

	// 写入新文件
	buf := bytes.Buffer{}
	buf.Write([]byte("---\n"))
	buf.Write(newMDHeaderYamlStr)
	buf.Write([]byte("---\n"))
	buf.WriteString(md.MDContent)
	if _, err := mdFile.Write(buf.Bytes()); err != nil {
		return errors.Wrapf(err, "replace write into blog file[%s] got err", md.Filepath)
	}

	return nil
}
