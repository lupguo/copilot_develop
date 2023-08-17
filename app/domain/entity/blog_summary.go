package entity

// BlogSummary 博客总结
type BlogSummary struct {
	Keywords    string `json:"keywords,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description"`
}

// BlogSummaryUpdatedRecord DB更新记录
type BlogSummaryUpdatedRecord struct {
	ID          uint   `gorm:"primarykey"`
	CreatedAt   string `gorm:"created_at"`
	UpdatedAt   string `gorm:"updated_at"`
	DeletedAt   string `gorm:"index"`
	Path        string `gorm:"path"`
	Title       string `gorm:"title"`
	Keywords    string `gorm:"keywords"`
	Summary     string `gorm:"summary"`
	Description string `gorm:"description"`
	Headers     string `gorm:"headers"`
}

func (t BlogSummaryUpdatedRecord) TableName() string {
	return "already_updated_blogs"
}
