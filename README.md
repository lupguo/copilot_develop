# AI Copilot Develop - AI 协同开发应用

~~AI Copilot Develop 为一款软件工程师在 MacOS 下的协同 AI 开发软件，旨在帮助研发人员在 MacOS 下流畅的编程.~~

P.S.: 目前仅应用在 BlogAISummary 摘要生成，其他功能还在开发中！

## Roadmap

1. [x] 支持 blog 的内容批量 keywords 提取、内容 summary 小结，并填补到 Blog 中 - 进度 85%
2. [ ] 文生图的能力，用于公众号读取、`wisdom-httpd`使用
3. [ ] ~~默认 AI 辅助角色支持(eg. 提供命名协助、日期 unixtime 处理、词条 Wikipedia 翻译)~~
4. [ ] 提供截图、图床功能、提供图片 AI 处理

## 灵感来源

1. Alfred
2. PopClip
3. Bob 翻译软件
4. OpenAI

## 问题

### 在 SummaryBlog 过程中遇到的问题

1. [x] 存储问题: 采用 SQlite 本地存储
2. [x] 配置问题: 还是采用 `config.Method()`形式，而非依赖注入方式，优点在于更灵活
3. [x] 文章内容过长，导致超过 OpenAI `gpt-3.5-turbo-16k` Token 阈值：
   - 之前统计方法有问题，参考 https://platform.openai.com/tokenizer 可以基于正则 `(\p{Han}|\b\w+\b)`
     匹配汉字和单词，和 Token 统计比较接近
   - 使用正则表达式替换，移除 Code，仅保留文章关键信息用于文章摘要生成
4. [x] pflag 支持，方便`tstorm.com`编写时候，快速补充博文摘要

#### OpenAI 限频问题: 每分钟只能有 18w token

```
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176976 / min. Contact us through our help center at help.openai.com if you continue to have issues.
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176963 / min. Contact us through our help center at help.openai.com if you continue to have issues.
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176830 / min. Contact us through our help center at help.openai.com if you continue to have issues.
```

**解决方案**

1. 失败重试(间隔一定时间)
2. 因为 OpenAI 在响应报文中包含 total_tokens，可以按每 min 统计，超过阈值延缓 OpenAI 并发请求，待时间窗口到期重置计数器
