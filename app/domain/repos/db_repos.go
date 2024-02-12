package repos

import (
	"context"

	"github.com/lupguo/copilot_develop/app/domain/entity"
)

type IReposSQLiteBlogSummary interface {

	// InitBlogSummaryDB 初始化blog_summary.db sqlite数据库
	InitBlogSummaryDB(ctx context.Context) error

	// SelBlogMDRecord 查询是否有Path处理的记录
	SelBlogMDRecord(ctx context.Context, path string) (*entity.BlogArticle, error)

	// CleanAllBlogSummaryDB 清理整个数据库记录
	CleanAllBlogSummaryDB(ctx context.Context) error

	// AddBlogMDRecord 新增一条已处理的MD记录
	AddBlogMDRecord(ctx context.Context, md *entity.BlogMD) error

	// UpdateBlogMDRecord 更新BlogMD信息
	UpdateBlogMDRecord(ctx context.Context, md *entity.BlogMD) error

	// ReplaceBlogMDRecord 当文档不存在时候新增，存在时候更新md内容
	ReplaceBlogMDRecord(ctx context.Context, md *entity.BlogMD) error
}
