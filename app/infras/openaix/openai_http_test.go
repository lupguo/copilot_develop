package openaix

import (
	"context"
	"testing"

	"github.com/lupguo/copilot_develop/app/infras/config"
	"github.com/sashabaranov/go-openai"
)

func TestOpenAIHttpProxyClient_Summary(t *testing.T) {
	// client
	filename := `../../../conf/app_dev.yaml`
	cfg, err := config.ParseAppConfig(filename)
	if err != nil {
		t.Error(err)
	}

	// client proxy
	aiProxyClient, err := NewOpenAIHttpProxyClient(cfg)
	if err != nil {
		t.Error(err)
	}

	// request
	req := &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo16K,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个Blog文章摘要汇总工具，用200字左右提炼出文章的中心思想，要求言简意赅",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "hello",
			},
		},
		MaxTokens: 1440,
	}
	t.Run("t1", func(t *testing.T) {
		gotSummary, err := aiProxyClient.DoAIChatCompletionRequest(context.Background(), req)
		if err != nil {
			t.Errorf("SummaryBlogMD() error = %v", err)
			return
		}
		t.Logf("summary:\n%+v", gotSummary)
	})
}
