package service

import (
	"context"
	"encoding/json"

	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/repos"
	"github.com/lupguo/copilot_develop/app/infras/openaix"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

const (
	PromptKeySummaryBlog = "summary-blog"
)

// IServicesSummaryAI AI汇总服务接口
type IServicesSummaryAI interface {
	// SummaryBlogMD 摘要总结+关键字
	SummaryBlogMD(ctx context.Context, md *entity.BlogMD) (summary *entity.ArticleSummary, err error)
}

// AIService AI汇总服务
type AIService struct {
	infra     repos.IReposOpenAI
	promptCfg map[string]*openaix.Prompt
}

// NewAIService 底层的SummaryAI服务
func NewAIService(infra repos.IReposOpenAI, promptFile string) (*AIService, error) {

	// 解析prom
	err, promptCfg := openaix.ParseAppPromptConfig(promptFile)
	if err != nil {
		return nil, errors.Wrap(err, "parse app prompt config got err")
	}

	return &AIService{
		infra:     infra,
		promptCfg: promptCfg,
	}, nil
}

// SummaryBlogMD 内容摘要+关键字总结
func (srv *AIService) SummaryBlogMD(ctx context.Context, md *entity.BlogMD) (summary *entity.ArticleSummary, err error) {
	// 获取指定key的提示词
	prompt, err := openaix.GetPrompt("summary-blog")
	if err != nil {
		return nil, errors.Wrap(err, "summary blog cannot found ai prompt key")
	}

	// 组装请求内容消息
	userMsg := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: md.MiniData.MiniContent,
	}}
	req := &openai.ChatCompletionRequest{
		Model:     prompt.AIMode,
		MaxTokens: prompt.MaxTokens,
		Messages:  append(prompt.PredefinedPrompts, userMsg...),
	}

	// 请求OpenAI获取响应
	resp, err := srv.infra.DoAIChatCompletionRequest(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "infra do ai chat completion request got err")
	}

	// 解析响应信息
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
