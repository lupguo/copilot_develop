package openaix

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hold7techs/go-shim/log"
	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/infras/config"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

const (
	UrlOpenAISummary = "/v3/"
)

// OpenAIHttpProxyClient OpenAI Http代理客户端
type OpenAIHttpProxyClient struct {
	proxyClient  *openai.Client
	appPromptMap config.AppPromptMap
}

// NewOpenAIHttpProxyClient 初始一个OpenAI代理实例
func NewOpenAIHttpProxyClient(cfg *config.AppConfig) (*OpenAIHttpProxyClient, error) {
	// 初始openAI配置
	openaiCfg := openai.DefaultConfig(cfg.App.AuthToken)

	// 设置代理地址
	proxyURL, err := url.Parse(cfg.App.SocksURL)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing proxy URL")
	}

	// 创建一个自定义的Transport，并设置代理
	openaiCfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 0,
	}

	return &OpenAIHttpProxyClient{
		proxyClient: openai.NewClientWithConfig(openaiCfg),
	}, nil
}

//
// // ExtractKeywords 提取content内容的关键字信息
// func (o *OpenAIHttpProxyClient) ExtractKeywords(ctx context.Context, content string) (keywords []string, err error) {
// 	promptKey := "summary-blog"
// 	prompt, ok := o.appPromptMap[promptKey]
// 	if !ok {
// 		return nil, errors.Errorf("prompt config key[%s] not exist", promptKey)
// 	}
//
// 	// 请求头
// 	userMsg := []openai.ChatCompletionMessage{{
// 		Role:    openai.ChatMessageRoleUser,
// 		Content: content,
// 	}}
// 	req := &openai.ChatCompletionRequest{
// 		Model:     prompt.AIMode,
// 		Messages:  append(prompt.PredefinedPrompts, userMsg...),
// 		MaxTokens: 1440,
// 	}
//
// 	// 请求OpenAI获取内容
// 	resp, err := o.DoAIChatCompletionRequest(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_ = resp
//
// 	// 取值
// 	return nil, nil
// }

// DoAIChatCompletionRequest 通用的AI ChatCompletion代理请求
func (o *OpenAIHttpProxyClient) DoAIChatCompletionRequest(ctx context.Context, req *openai.ChatCompletionRequest) (response *openai.ChatCompletionResponse, err error) {
	log.Debugf("AI REQ:\n%s", shim.ToJsonString(req, true))
	resp, err := o.proxyClient.CreateChatCompletion(ctx, *req)
	if err != nil {
		log.Errorf("CreateChatCompletion error: %v\n", err)
		return nil, errors.Wrap(err, "do AI chat completion request got err")
	}
	log.Debugf("AI RESP:\n%s", shim.ToJsonString(resp, true))

	return &resp, nil
}
