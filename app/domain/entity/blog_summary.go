package entity

// BlogSummary 博客总结
type BlogSummary struct {
	Keywords    string `json:"keywords,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description"`
}
