package repos

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// IReposOpenAI 负责和底层OpenAI交互的接口
type IReposOpenAI interface {
	// // SummaryBlogMD 摘要总结
	// SummaryBlogMD(ctx context.Context, content string) (summary string, err error)
	//
	// // ExtractKeywords 内容关键字提炼
	// ExtractKeywords(ctx context.Context, content string) (keywords []string, err error)

	// DoAIChatCompletionRequest 通用的AI ChatCompletion代理请求
	DoAIChatCompletionRequest(ctx context.Context, req *openai.ChatCompletionRequest) (response *openai.ChatCompletionResponse, err error)
}
