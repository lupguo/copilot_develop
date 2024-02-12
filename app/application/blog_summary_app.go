package application

import (
	"context"
	"path/filepath"

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
	return &BlogSummaryApp{
		aiSrv:       aiSrv,
		sqliteInfra: sqliteInfra,
	}
}

// UpdateBlogHeaderYaml 并发更新Blog的汇总信息
func (app *BlogSummaryApp) UpdateBlogHeaderYaml(ctx context.Context, storageRoot string) error {
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
			if err = app.updateBlogYamlHeader(ctx, mdPath); err != nil {
				return intershim.LogAndWrapf(err, "replace summary for md file[%s]  got err", mdPath)
			}
			return nil
		})
	}

	if err := egp.Wait(); err != nil {
		return intershim.LogAndWrapf(err, "egp got err: %s", err)
	}

	return nil
}

// updateBlogYamlHeader 结合DB有替换记录、ForceUpdate是否被设置成true，决策是否需要刷新HeaderYaml头部
func (app *BlogSummaryApp) updateBlogYamlHeader(ctx context.Context, mdfile string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("app panic recover for path[%v]: %v", mdfile, err)
		}
	}()

	// 基于本地文件，初始每个md
	md, err := entity.NewBlogMD(mdfile)
	if err != nil {
		return errors.Wrapf(err, "app new md[%s] got err", mdfile)
	}

	// DB查看是否存在mdPath已Replace过了
	record, err := app.sqliteInfra.SelBlogMDRecord(ctx, mdfile)
	if err != nil { // db error
		return err
	} else if record != nil && md.NeedForceUpdate() == false { // 有记录和无强刷，则直接返回
		return nil
	}

	// 通过AIService更新md内容
	if md.MDHeader.ForceUpdate == entity.UpdateALL {
		if err := app.refreshBlogSummaryAndKeywords(ctx, md); err != nil {
			return errors.Wrapf(err, "app refreash md[%s] blog summary and keywords got err", mdfile)
		}
	}

	// 重置强制更新字段，设置为默认空值
	md.MDHeader.ForceUpdate = ""
	if err = md.ReplaceWithNewYamlHeader(); err != nil {
		return errors.Wrapf(err, "app replace write into blog md[%s] got err", mdfile)
	}

	// 新增或者更改 MD Record记录
	if err = app.sqliteInfra.ReplaceBlogMDRecord(ctx, md); err != nil {
		return errors.Wrapf(err, "app replace md[%s] db's record got err", mdfile)
	}

	return nil
}

// 刷新Blog的Summary和Keywords信息
func (app *BlogSummaryApp) refreshBlogSummaryAndKeywords(ctx context.Context, md *entity.BlogMD) error {
	// 内容太少了，不做AI生成
	if md.IsContentWordsTooSmall() {
		return errors.Errorf("content is too small, needn't request OpenAI")
	}

	// 检测md是否符合replace规则
	limitSize := entity.OpenAIMaxTokenSize
	if md.IsMinContentTooLong(limitSize) {
		return errors.Errorf("min content is over max token size(%d), cannot request OpenAI", limitSize)
	}

	// 使用openAI生成blog文章内容摘要
	summary, err := app.aiSrv.SummaryBlogMD(ctx, md)
	if err != nil {
		return errors.Wrapf(err, "aiSrv summary blog content got err")
	}

	// 汇总、关键字、描述，将调整后的md更新回去
	md.MDHeader.Summary = summary.Summary
	md.MDHeader.Keywords = summary.Keywords
	md.MDHeader.Description = summary.Description

	return nil
}
