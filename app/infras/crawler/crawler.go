package crawler

import (
	"context"
	"io"
)

// 最大深度
var bfsMaxDepth int

// PageUrlRegex 针对pageLink页面，爬取Content和Url内容
type PageUrlRegex struct {
	Content string // 页面内容url
	Url     string // 页面内的URL正则
	// MaxDepth int    // 页面爬取最大深度
}

// PageUrlRegexConfig 同类型的页面url的提取规则
type PageUrlRegexConfig map[string]PageUrlRegex

// ICrawler 爬虫客户端接口
type ICrawler interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context, stop chan struct{}) error

	// GetPageUrlDoc 解析Doc，得到文本内容和urls信息
	// 1. 设定如何提取文本内容？
	// 2. 设定提取URLs规则？
	GetPageUrlDoc(ctx context.Context, pageUrl string, regex *PageUrlRegex) (doc io.Reader, urls []string, err error)

	// UnderstandDocByAI 通过AI解析
	// UnderstandDocByAI(ctx context.Context, url string, regexCfg *PageUrlRegex) (content []byte, urls []string, err error)

	// FoundURLsFromUrlBFS 通过BFS找到URL
	FoundURLsFromUrlBFS(ctx context.Context, url string, maxDepth int) (urls []string, err error)
}

type Crawler struct {
	client ICrawler
}
