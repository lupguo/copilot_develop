package openaix

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/config"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
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

// DoAIChatCompletionRequest 通用的AI ChatCompletion代理请求
func (o *OpenAIHttpProxyClient) DoAIChatCompletionRequest(ctx context.Context, req *openai.ChatCompletionRequest) (response *openai.ChatCompletionResponse, err error) {
	resp, err := o.proxyClient.CreateChatCompletion(ctx, *req)
	if err != nil {
		log.Errorf("DoAIChatCompletionRequest() got error: %v\n", err)
		return nil, errors.Wrap(err, "do AI chat completion request got err")
	}

	// 精简打印请求和响应信息
	// req.Messages[0].Content
	log.Debugf("\nAI REQ:\n%s\nAI RESP:\n%s", shim.ToJsonString(req, true), shim.ToJsonString(resp, true))

	return &resp, nil
}
