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
func (b *BlogSummarySqliteInfra) SelBlogMDRecord(ctx context.Context, path string) (*entity.BlogArticle, error) {
	var record entity.BlogArticle
	err := b.db.First(&record, "path=?", path).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db sql[SelBlogMDRecord] got err")
	}

	return &record, nil
}

// InitBlogSummaryDB 初始化
func (b *BlogSummarySqliteInfra) InitBlogSummaryDB(ctx context.Context) error {
	if err := b.db.AutoMigrate(&entity.BlogArticle{}); err != nil {
		return errors.Wrap(err, "db sql[InitBlogSummaryDB] got err")
	}
	return nil
}

// CleanAllBlogSummaryDB 清理整个DB记录
func (b *BlogSummarySqliteInfra) CleanAllBlogSummaryDB(ctx context.Context) error {
	if err := b.db.Delete(&entity.BlogArticle{}).Error; err != nil {
		return errors.Wrap(err, "db sql[CleanAllBlogSummaryDB] got err")
	}
	return nil
}

// AddBlogMDRecord 新增DB记录
func (b *BlogSummarySqliteInfra) AddBlogMDRecord(ctx context.Context, md *entity.BlogMD) error {
	header := md.MDHeader
	err := b.db.Create(&entity.BlogArticle{
		CreatedAt:   time.Now().Format(shim.StdDateTimeLayout),
		Path:        md.Filepath,
		Title:       header.Title,
		Keywords:    header.Keywords,
		Summary:     header.Summary,
		Description: header.Description,
		Headers:     header.String(),
	}).Error
	if err != nil {
		return errors.Wrap(err, "db sql[AddBlogMDRecord] got err")
	}

	return nil
}

// UpdateBlogMDRecord  更新Blog Md记录
func (b *BlogSummarySqliteInfra) UpdateBlogMDRecord(ctx context.Context, md *entity.BlogMD) error {
	err := b.db.Where("path=?", md.Filepath).Updates(&entity.BlogArticle{
		Draft:  md.MDHeader.Draft,
		Weight: md.MDHeader.Weight,
	}).Error

	if err != nil {
		return errors.Wrap(err, "db sql[AddBlogMDRecord] got err")
	}

	return nil
}
