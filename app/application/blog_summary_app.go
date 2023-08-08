package application

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/hold7techs/go-shim/log"
	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// BlogSummaryApp Blog的汇总App
type BlogSummaryApp struct {
	srv service.IServicesSummaryAI
}

// 博客地址
const blogStorageRoot = "/data/www/tkstorm.com/content/posts"

// UpdateBlogSummaryContent 更新Blog的汇总信息
func (app *BlogSummaryApp) UpdateBlogSummaryContent() error {
	// 查询目录下所有的markdown目录 -> slice内 []*BlogMD
	blogFilePaths, err := shim.FindFilePaths(blogStorageRoot, "*.md")
	if err != nil {
		return errors.Wrapf(err, "shim find file paths root [%s] got err", blogStorageRoot)
	}

	// 通过正则提取md的主题内容
	ctx := context.Background()
	for _, blogFilePath := range blogFilePaths {
		// 替换每个md的汇总信息、关键字、描述信息
		if err := app.ReplaceKeywordsAndSummary(ctx, blogFilePath); err != nil {
			log.Errorf("replace blog md file[%s] got err: %s", blogFilePath, err)
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

	// content 提前摘要
	summary, err := app.srv.BlogSummary(ctx, md.MDContent)
	if err != nil {
		return errors.Wrap(err, "using open ai extract summary for blogContent got err")
	} else if summary == nil {
		return errors.New("app srv blog summary got nil val")
	}

	// 检测summary结果
	if summary.Summary == "" || summary.Keywords == "" || summary.Description == "" {
		return errors.Wrapf(err, "blog summary empty values, summary: %s\n keywords: %s\n, description: %s\n",
			summary.Summary, summary.Keywords, summary.Description)
	}

	// 汇总、关键字、描述，将调整后的md更新回去
	md.MDHeader.Summary = summary.Summary
	md.MDHeader.Keywords = summary.Keywords
	md.MDHeader.Description = summary.Description
	headerStr, err := yaml.Marshal(md.MDHeader)
	if err != nil {
		return errors.Wrapf(err, "marsh file[%s] yaml head got err", blogFilePath)
	}
	log.Debugf("newMDHeaderStr: %s", headerStr)

	// 清空并写入新的文件内容
	if err := writeNewYamlHeader(md, headerStr); err != nil {
		return errors.Wrapf(err, "replace write into blog file[%s] got err", md.Filepath)
	}

	return nil
}

func writeNewYamlHeader(md *entity.BlogMD, headerStr []byte) error {
	// 清理+重写如新的文件内容
	mdFile, err := os.OpenFile(md.Filepath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "open blog file[%s] got err", md.Filepath)
	}
	defer mdFile.Close()
	w := bufio.NewWriter(mdFile)
	if _, err = fmt.Fprintf(w, "---\n%s---\n\n%s", headerStr, md.MDContent); err != nil {
		return errors.Wrapf(err, "write into blog file[%s] with new yaml header got err", md.Filepath)
	}
	log.Debugf("buffed bytes[%s]: %d", md.Filepath, w.Buffered())
	if err := w.Flush(); err != nil {
		return errors.Wrapf(err, "flush failed")
	}
	return nil
}
