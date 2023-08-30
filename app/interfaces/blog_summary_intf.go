package interfaces

import (
	"github.com/lupguo/copilot_develop/app/application"
)

// CopilotDevelop 助手
type CopilotDevelop struct {
	blogSummaryApp *application.BlogSummaryApp
}

// UpdateBlogSummary 更新BlogSummary信息
// func (c *CopilotDevelop) UpdateBlogSummary() error {
// 	ctx := context.Background()
// 	// 博客地址
// 	storageRoot := "/data/www/tkstorm.com/content/posts"
//
// 	return c.blogSummaryApp.UpdateBlogSummaryContent(ctx, storageRoot)
// }
