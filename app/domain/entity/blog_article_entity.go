package entity

// ArticleSummary 博客总结
type ArticleSummary struct {
	Keywords    string `json:"keywords,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description"`
}

// BlogArticle DB更新记录
type BlogArticle struct {
	ID          uint   `gorm:"id"`
	CreatedAt   string `gorm:"created_at"`
	UpdatedAt   string `gorm:"updated_at"`
	DeletedAt   string `gorm:"index"`
	Path        string `gorm:"path"`        // 文章本地存储路径(目前作为唯一的标识)
	ShortMark   string `gorm:"short_mark"`  // 文章短标记，文章创建后自动生成，基于文章标题做短hash，支持后续软链接快速检索到文章
	Title       string `gorm:"title"`       // 文章标题
	Categories  string `gorm:"categories"`  // 文章类型
	Tags        string `gorm:"tags"`        // 文章标签
	Draft       bool   `gorm:"draft"`       // 是否手稿
	Weight      int    `gorm:"weight"`      // 文章权重
	WordCount   int    `gorm:"word_count"`  // 文章内容长度，可以用于指导hugo文章的基本情况(新增和更新时候都会用到)
	Keywords    string `gorm:"keywords"`    // 文章关键字
	Summary     string `gorm:"summary"`     // 文章摘要
	Description string `gorm:"description"` // 文章描述
	Aliases     string `gorm:"aliases"`     // 软连
}

func (t BlogArticle) TableName() string {
	return "blog_articles"
}
