package langchainx

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/tmc/langchaingo/llms/openai"
)

type langChainX struct {
}

const APIToken = `Your OpenAI-Token`
const OpenAIMode = `gpt-3.5-turbo-16k`
const SocksURL = `socks5://127.0.0.1:10553`

func langchainSample01() {
	proxyHttpClient, err := newSocksHTTPProxy(SocksURL)
	opts := []openai.Option{
		openai.WithToken(APIToken),
		openai.WithModel(OpenAIMode),
		openai.WithHTTPClient(proxyHttpClient),
	}
	llm, err := openai.New(opts...)
	if err != nil {
		log.Fatal(err)
	}
	prompt := "What would be a good company name for a company that makes colorful socks?"
	ctx := context.Background()
	completion, err := llm.Call(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("chan1:%s", completion)

	completion, err = llm.Call(ctx, "请基于给定的单词，输出有创意的中文名称公司（5个）:"+completion)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("chan2:%s", completion)
}

type DebugTransport struct {
	Transport http.RoundTripper
}

func newSocksHTTPProxy(socksURL string) (*http.Client, error) {
	// 设置代理地址
	proxyURL, err := url.Parse(socksURL)
	if err != nil {
		return nil, err
	}

	// 创建一个自定义的Transport，并设置代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 0,
	}

	return client, nil
}
