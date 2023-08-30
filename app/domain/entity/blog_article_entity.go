package entity

// ArticleSummary 博客总结
type ArticleSummary struct {
	Keywords    string `json:"keywords,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description"`
}

// BlogArticle DB更新记录
type BlogArticle struct {
	ID            uint   `gorm:"id"`
	CreatedAt     string `gorm:"created_at"`
	UpdatedAt     string `gorm:"updated_at"`
	DeletedAt     string `gorm:"index"`
	Path          string `gorm:"path"`
	Title         string `gorm:"title"`
	Draft         bool   `gorm:"draft"`
	Weight        int    `gorm:"weight"`
	ContentLength int    `gorm:"content_length"` // 内容长度，可以用于指导hugo文章的基本情况(新增和更新时候都会用到)
	ShortMark     string `gorm:"short_mark"`     // 短标记，文章创建后自动生成，支持后续软链接快速检索到文章
	Keywords      string `gorm:"keywords"`
	Summary       string `gorm:"summary"`
	Description   string `gorm:"description"`
	Headers       string `gorm:"headers"`
}

func (t BlogArticle) TableName() string {
	return "blog_articles"
}
