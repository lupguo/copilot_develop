package service

import (
	"context"
	"encoding/json"

	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/lupguo/copilot_develop/app/domain/repos"
	"github.com/lupguo/copilot_develop/app/infras/config"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

// IServicesSummaryAI AI汇总服务接口
type IServicesSummaryAI interface {
	// BlogSummary 摘要总结+关键字
	BlogSummary(ctx context.Context, content string) (summary *entity.BlogSummary, err error)
}

// AIService AI汇总服务
type AIService struct {
	Infra        repos.IReposOpenAI
	appPromptMap config.AppPromptMap
}

// NewAIService 底层的SummaryAI服务
func NewAIService(infra repos.IReposOpenAI, appPromptMap config.AppPromptMap) *AIService {
	return &AIService{Infra: infra, appPromptMap: appPromptMap}
}

// BlogSummary 内容摘要+关键字总结
func (srv *AIService) BlogSummary(ctx context.Context, content string) (summary *entity.BlogSummary, err error) {
	promptKey := config.PromptKeySummaryBlog
	prompt, ok := srv.appPromptMap[promptKey]
	if !ok {
		return nil, errors.Errorf("prompt config key[%s] not exist", promptKey)
	}

	// 请求头
	userMsg := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}}
	req := &openai.ChatCompletionRequest{
		Model:     prompt.AIMode,
		Messages:  append(prompt.PredefinedPrompts, userMsg...),
		MaxTokens: 4000,
	}

	// 请求OpenAI获取内容
	resp, err := srv.Infra.DoAIChatCompletionRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	summary = &entity.BlogSummary{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), summary); err != nil {
		return nil, errors.Wrap(err, "the blog summary received response from AI proxy, attempted to unmarshal resp content but got an error")
	}

	// 取值
	return summary, nil
}
