package dbs

import (
	"context"
	"time"

	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BlogSummarySqliteInfra struct {
	db *gorm.DB
}

func NewBlogSummarySqliteInfra(sqlDBFile string) (*BlogSummarySqliteInfra, error) {
	db, err := gorm.Open(sqlite.Open(sqlDBFile), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite.Open(%s) got err", sqlDBFile)
	}
	return &BlogSummarySqliteInfra{db: db}, nil
}

// SelBlogMDRecord 查询BlogMD记录
func (infra *BlogSummarySqliteInfra) SelBlogMDRecord(ctx context.Context, path string) (*entity.BlogArticle, error) {
	var record entity.BlogArticle
	err := infra.db.Debug().
		First(&record, "path=?", path).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db sql[SelBlogMDRecord] got err")
	}

	return &record, nil
}

// InitBlogSummaryDB 初始化
func (infra *BlogSummarySqliteInfra) InitBlogSummaryDB(ctx context.Context) error {
	if err := infra.db.AutoMigrate(&entity.BlogArticle{}); err != nil {
		return errors.Wrap(err, "db sql[InitBlogSummaryDB] got err")
	}
	return nil
}

// CleanAllBlogSummaryDB 清理整个DB记录
func (infra *BlogSummarySqliteInfra) CleanAllBlogSummaryDB(ctx context.Context) error {
	if err := infra.db.Delete(&entity.BlogArticle{}).Error; err != nil {
		return errors.Wrap(err, "db sql[CleanAllBlogSummaryDB] got err")
	}
	return nil
}

// AddBlogMDRecord 新增DB记录
func (infra *BlogSummarySqliteInfra) AddBlogMDRecord(ctx context.Context, md *entity.BlogMD) error {
	header := md.MDHeader
	err := infra.db.Debug().
		Create(&entity.BlogArticle{
			CreatedAt:   time.Now().Format(shim.StdDateTimeLayout),
			UpdatedAt:   time.Now().Format(shim.StdDateTimeLayout),
			Date:        header.Date,
			Path:        md.Filepath,
			ShortMark:   header.ShortMark,
			Title:       header.Title,
			Categories:  shim.ToJsonString(header.Categories, false),
			Tags:        shim.ToJsonString(header.Tags, false),
			Draft:       header.Draft,
			Weight:      header.Weight,
			WordCount:   header.WordCounts,
			Keywords:    header.Keywords,
			Summary:     header.Summary,
			Description: header.Description,
			Aliases:     shim.ToJsonString(header.Aliases, false),
		}).Error
	if err != nil {
		return errors.Wrap(err, "db sql[AddBlogMDRecord] got err")
	}

	return nil
}

// UpdateBlogMDRecord  更新Blog Md记录
func (infra *BlogSummarySqliteInfra) UpdateBlogMDRecord(ctx context.Context, md *entity.BlogMD) error {
	header := md.MDHeader
	err := infra.db.Debug().
		Where("path=?", md.Filepath).
		Updates(&entity.BlogArticle{
			UpdatedAt:   time.Now().Format(shim.StdDateTimeLayout),
			Date:        header.Date,
			Path:        md.Filepath,
			ShortMark:   header.ShortMark,
			Title:       header.Title,
			Keywords:    header.Keywords,
			Summary:     header.Summary,
			Description: header.Description,
			Categories:  shim.ToJsonString(header.Categories, false),
			Tags:        shim.ToJsonString(header.Tags, false),
			Draft:       header.Draft,
			Weight:      header.Weight,
			WordCount:   header.WordCounts,
			Aliases:     shim.ToJsonString(header.Aliases, false),
		}).Error
	if err != nil {
		return errors.Wrap(err, "db sql[AddBlogMDRecord] got err")
	}

	return nil
}

// ReplaceBlogMDRecord  当文档不存在时候新增，存在时候更新md内容
func (infra *BlogSummarySqliteInfra) ReplaceBlogMDRecord(ctx context.Context, md *entity.BlogMD) error {
	// 查询是否存在
	mdRecord, err := infra.SelBlogMDRecord(ctx, md.Filepath)

	switch {
	case err != nil:
		return errors.Wrap(err, "replace blog md record, sel got err")
	case mdRecord == nil:
		return infra.AddBlogMDRecord(ctx, md)
	default:
		return infra.UpdateBlogMDRecord(ctx, md)
	}
}
