package application

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hold7techs/go-shim/log"
	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/repos"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// BlogSummaryApp Blog的汇总App
type BlogSummaryApp struct {
	srv         service.IServicesSummaryAI
	sqliteInfra repos.IReposSQLiteBlogSummary
}

// NewBlogSummaryApp 初始一个BlogSummaryApp
func NewBlogSummaryApp(srv service.IServicesSummaryAI, sqliteInfra repos.IReposSQLiteBlogSummary) *BlogSummaryApp {
	return &BlogSummaryApp{srv: srv, sqliteInfra: sqliteInfra}
}

// UpdateBlogSummaryContent 更新Blog的汇总信息
func (app *BlogSummaryApp) UpdateBlogSummaryContent(ctx context.Context, storageRoot string) error {
	// 查询目录下所有的markdown目录 -> slice内 []*BlogMD
	blogFilePaths, err := shim.FindFilePaths(storageRoot, "*.md")
	if err != nil {
		return errors.Wrapf(err, "shim find file paths root [%s] got err", storageRoot)
	}

	// 基于条件过滤掉"path为_index.md"的
	blogFilePaths = shim.ProcessStringsSlice(blogFilePaths, func(path string) bool {
		return filepath.Base(path) == "_index.md"
	}, nil)

	// 通过正则提取md的主题内容 - 改并发版本
	egp := errgroup.Group{}
	egp.SetLimit(10)
	for _, blogFilePath := range blogFilePaths {
		mdPath := blogFilePath
		egp.Go(func() error {
			// 替换每个md的汇总信息、关键字、描述信息
			if err := app.replaceMdSummary(ctx, mdPath); err != nil {
				return errors.Wrapf(err, "replace blog md file[%s] got err(1): %s", mdPath, err)
			}

			// 替换每个md的weight(文本字数+时间排序)、draft设置
			if err := app.updateMdWeightAndDraft(ctx, mdPath); err != nil {
				return errors.Wrapf(err, "replace blog md file[%s] got err(2): %s", mdPath, err)
			}

			return nil
		})
	}

	if err := egp.Wait(); err != nil {
		log.Errorf("egp got err: %s", err)
		return err
	}

	return nil
}

// replaceMdSummary 将keywords, summary 填充到原有的blog文章内
func (app *BlogSummaryApp) replaceMdSummary(ctx context.Context, mdPath string) error {
	// DB查看是否存在mdPath已Replace过了
	record, err := app.sqliteInfra.SelBlogMDRecord(ctx, mdPath)
	if err != nil {
		return err
	} else if record != nil {
		log.Warnf("md path[%s] already replaced", mdPath)
		return nil
	}

	// 初始每个md
	md, err := entity.NewBlogMD(mdPath)
	if err != nil {
		log.Warnf("entity new blog md [%s] got err: %s", mdPath, err)
		return errors.Wrap(err, "entity new blog md got err in replace")
	}

	// 检测md是否符合replace规则
	if md.IsOverMaxTokenSize() {
		log.Warnf("md[%s] content is over max token size", mdPath)
		return nil
	}

	if md.IsContentTooSmall() {
		log.Warnf("md[%s] content is too small", mdPath)
		return nil
	}

	// content 提前摘要
	summary, err := app.srv.BlogSummary(ctx, md)
	if err != nil {
		return errors.Wrap(err, "srv blog summary got err")
	}

	// 汇总、关键字、描述，将调整后的md更新回去
	md.MDHeader.Summary = summary.Summary
	md.MDHeader.Keywords = summary.Keywords
	md.MDHeader.Description = summary.Description
	if err := md.ReplaceWithNewYamlHeader(); err != nil {
		return errors.Wrapf(err, "replace write into blog file[%s] got err", md.Filepath)
	}

	// 添加MD Record记录
	if err := app.sqliteInfra.AddBlogMDRecord(ctx, md); err != nil {
		log.Errorf("insert md file[%s] summary record got err: %s", md.Filepath, err)
	}

	return nil
}

// updateMdWeightAndDraft 更新文章权重和Draft信息
// 1. 文章长度过短的，更新处理
func (app *BlogSummaryApp) updateMdWeightAndDraft(ctx context.Context, mdPath string) error {
	// 初始每个md
	md, err := entity.NewBlogMD(mdPath)
	if err != nil {
		log.Errorf("entity new blog md [%s] got err: %s", mdPath, err)
		return errors.Wrap(err, "entity new blog md got err in replace")
	}

	// MD信息更新
	header := md.MDHeader
	header.Draft = md.IsDraft()                                                           // 是否手稿
	header.Weight = md.CalcArticleWeight()                                                // 文章权重
	header.ShortMark = md.ShortMark()                                                     // 文章签名
	header.WordCounts = md.WordCount()                                                    // 文章字符统计
	header.Tags = shim.ProcessStringsSlice(header.Tags, nil, strings.ToLower)             // 文章标签统一转小写
	header.Categories = shim.ProcessStringsSlice(header.Categories, nil, strings.ToLower) // 文章标签统一转小写

	// db信息更新
	if err := app.sqliteInfra.UpdateBlogMDRecord(ctx, md); err != nil {
		return errors.Wrap(err, "update blog weight and draft got err")
	}

	// 用新的yaml header替换
	if err := md.ReplaceWithNewYamlHeader(); err != nil {
		return errors.Wrapf(err, "replace write into blog file[%s] got err", md.Filepath)
	}

	return nil
}
