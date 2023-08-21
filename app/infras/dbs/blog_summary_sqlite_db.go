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

func (b *BlogSummarySqliteInfra) SelSummaryRecord(ctx context.Context, path string) (*entity.BlogSummaryUpdatedRecord, error) {
	var record entity.BlogSummaryUpdatedRecord
	err := b.db.First(&record, "path=?", path).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db sql[SelSummaryRecord] got err")
	}

	return &record, nil
}

func (b *BlogSummarySqliteInfra) InitBlogSummaryDB(ctx context.Context) error {
	if err := b.db.AutoMigrate(&entity.BlogSummaryUpdatedRecord{}); err != nil {
		return errors.Wrap(err, "db sql[InitBlogSummaryDB] got err")
	}
	return nil
}

func (b *BlogSummarySqliteInfra) CleanAllBlogSummaryDB(ctx context.Context) error {
	if err := b.db.Delete(&entity.BlogSummaryUpdatedRecord{}).Error; err != nil {
		return errors.Wrap(err, "db sql[CleanAllBlogSummaryDB] got err")
	}
	return nil
}

func (b *BlogSummarySqliteInfra) AddMDSummaryRecord(ctx context.Context, md *entity.BlogMD) error {
	header := md.MDHeader
	err := b.db.Create(&entity.BlogSummaryUpdatedRecord{
		CreatedAt:   time.Now().Format(shim.StdDateTimeLayout),
		Path:        md.Filepath,
		Title:       header.Title,
		Keywords:    header.Keywords,
		Summary:     header.Summary,
		Description: header.Description,
		Headers:     header.String(),
	}).Error
	if err != nil {
		return errors.Wrap(err, "db sql[AddMDSummaryRecord] got err")
	}

	return nil
}
