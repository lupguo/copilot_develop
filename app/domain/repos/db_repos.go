package repos

import (
	"context"

	"github.com/lupguo/copilot_develop/app/domain/entity"
)

type IReposSQLiteBlogSummary interface {

	// InitBlogSummaryDB 初始化blog_summary.db sqlite数据库
	InitBlogSummaryDB(ctx context.Context) error

	// SelSummaryRecord 查询是否有Path处理的记录
	SelSummaryRecord(ctx context.Context, path string) (*entity.BlogSummaryUpdatedRecord, error)

	// CleanAllBlogSummaryDB 清理整个数据库记录
	CleanAllBlogSummaryDB(ctx context.Context) error

	// AddMDSummaryRecord 新增一条已处理的MD记录
	AddMDSummaryRecord(ctx context.Context, md *entity.BlogMD) error
}
