package service

import (
	"context"
	"encoding/json"

	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/repos"
	"github.com/lupguo/copilot_develop/config"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

// IServicesSummaryAI AI汇总服务接口
type IServicesSummaryAI interface {
	// SummaryBlogMD 摘要总结+关键字
	SummaryBlogMD(ctx context.Context, md *entity.BlogMD) (summary *entity.ArticleSummary, err error)
}

// AIService AI汇总服务
type AIService struct {
	Infra repos.IReposOpenAI
}

// NewAIService 底层的SummaryAI服务
func NewAIService(infra repos.IReposOpenAI) *AIService {
	return &AIService{
		Infra: infra,
	}
}

// SummaryBlogMD 内容摘要+关键字总结
func (srv *AIService) SummaryBlogMD(ctx context.Context, md *entity.BlogMD) (summary *entity.ArticleSummary, err error) {
	// prompt key
	key := config.PromptKeySummaryBlog
	prompt, err := config.GetPrompt(key)
	if err != nil {
		return nil, errors.Wrap(err, "summary blog cannot found ai prompt key")
	}

	// request OpenAI chat completion
	userMsg := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: md.MiniData.MiniContent,
	}}
	req := &openai.ChatCompletionRequest{
		Model:     prompt.AIMode,
		MaxTokens: prompt.MaxTokens,
		Messages:  append(prompt.PredefinedPrompts, userMsg...),
	}

	// 请求OpenAI获取内容
	resp, err := srv.Infra.DoAIChatCompletionRequest(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "infra do ai chat completion request got err")
	}

	// 响应信息
	summary = &entity.ArticleSummary{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), summary); err != nil {
		return nil, errors.Wrap(err, "the blog summary received response from AI proxy, attempted to unmarshal resp content but got an error")
	}

	// 检测summary结果
	if summary.Summary == "" || summary.Keywords == "" || summary.Description == "" {
		return nil, errors.Wrapf(err, "blog summary empty values, summary: %s\n keywords: %s\n, description: %s\n",
			summary.Summary, summary.Keywords, summary.Description)
	}

	return summary, nil
}
