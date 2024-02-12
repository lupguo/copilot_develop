package application

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/repos"
	"github.com/lupguo/copilot_develop/app/domain/service"
	"github.com/lupguo/copilot_develop/internal/intershim"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// BlogSummaryApp Blog的汇总App
type BlogSummaryApp struct {
	aiSrv       service.IServicesSummaryAI
	sqliteInfra repos.IReposSQLiteBlogSummary
}

// NewBlogSummaryApp 初始一个BlogSummaryApp
func NewBlogSummaryApp(aiSrv service.IServicesSummaryAI, sqliteInfra repos.IReposSQLiteBlogSummary) *BlogSummaryApp {
	return &BlogSummaryApp{aiSrv: aiSrv, sqliteInfra: sqliteInfra}
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
			// 1. 更新md的摘要信息、关键字、描述信息
			if err := app.updateBlogSummaryInfos(ctx, mdPath); err != nil {
				return intershim.LogAndWrapf(err, "replace summary for md file[%s]  got err", mdPath)
			}

			// 2. 更新md的基本信息，例如文本字数+时间排序、draft设置
			if err := app.updateBlogBasicInfos(ctx, mdPath); err != nil {
				return intershim.LogAndWrapf(err, "replace md file[%s] got err(2): %s", mdPath, err)
			}

			return nil
		})
	}

	if err := egp.Wait(); err != nil {
		return intershim.LogAndWrapf(err, "egp got err: %s", err)
	}

	return nil
}

// updateBlogSummaryInfos 将keywords, summary 填充到原有的blog文章内
//   - 非Draft文章才可以被生成摘要
//   - 内容太少、太长都不生成摘要
func (app *BlogSummaryApp) updateBlogSummaryInfos(ctx context.Context, mdfile string) error {
	// DB查看是否存在mdPath已Replace过了
	record, err := app.sqliteInfra.SelBlogMDRecord(ctx, mdfile)
	if err != nil { // db error
		return err
	} else if record != nil {
		return errors.Wrapf(err, "md[%s] already replaced", mdfile)
	}

	// 基于本地文件，初始每个md
	md, err := entity.NewBlogMD(mdfile)
	if err != nil {
		return errors.Wrapf(err, "new md [%s] got err", mdfile)
	}

	// 检测md是否为Draft文章，若为Draft文章不做更新
	if md.IsDraft() {
		return errors.Errorf("md[%s] is draft, would not replace md summary, try update draft to false", mdfile)
	}

	// 内容太少了，不做AI生成
	if md.IsContentWordsTooSmall() {
		return errors.Errorf("md[%s] min content is too small", mdfile)
	}

	// 检测md是否符合replace规则
	limitSize := entity.OpenAIMaxTokenSize
	if md.IsMinContentTooLong(limitSize) {
		return errors.Errorf("md[%s] min content is over max token size(%d)", mdfile, limitSize)
	}

	// 使用openAI生成blog文章内容摘要
	summary, err := app.aiSrv.SummaryBlogMD(ctx, md)
	if err != nil {
		return errors.Wrapf(err, "ai srv summary blog md[%s] got err", mdfile)
	}

	// 汇总、关键字、描述，将调整后的md更新回去
	md.MDHeader.Summary = summary.Summary
	md.MDHeader.Keywords = summary.Keywords
	md.MDHeader.Description = summary.Description
	if err = md.ReplaceWithNewYamlHeader(); err != nil {
		return errors.Wrapf(err, "replace write into blog file[%s] got err", md.Filepath)
	}

	// 添加MD Record记录
	if err = app.sqliteInfra.AddBlogMDRecord(ctx, md); err != nil {
		return errors.Wrapf(err, "insert md file[%s] to sqlite db record got err", md.Filepath)
	}

	return nil
}

// updateBlogBasicInfos 更新文章权重和Draft信息
// 1. 文章长度过短的，更新处理
func (app *BlogSummaryApp) updateBlogBasicInfos(ctx context.Context, mdPath string) error {
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
	header.Categories = shim.ProcessStringsSlice(header.Categories, nil, strings.ToLower) // 文章分类统一转小写
	header.Tags = shim.ProcessStringsSlice(header.Tags, nil, strings.ToLower)             // 文章标签统一转小写

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
