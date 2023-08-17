# AI Copilot Develop - AI 协同开发应用

AI Copilot Develop 为一款软件工程师在MacOS下的协同AI开发软件，旨在帮助研发人员在MacOS下流畅的编程.

## Roadmap

1. [-] 支持blog的内容批量keywords提取、内容summary小结，并填补到Blog中 - 进度85%
2. [-] 默认AI辅助角色支持(eg. 提供命名协助、日期unixtime处理、词条Wikipedia翻译)
3. [-] 提供截图、图床功能、提供图片AI处理
4. [-] 文生图的能力，用于公众号读取

## 灵感来源

1. Alfred
2. PopClip
3. Bob翻译软件
4. OpenAI

## 问题

### 在SummaryBlog过程中遇到的问题

#### OpenAI限频问题: 每分钟只能有18w token
```
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176976 / min. Contact us through our help center at help.openai.com if you continue to have issues.
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176963 / min. Contact us through our help center at help.openai.com if you continue to have issues.
2023/08/17 16:35:46 openai_http.go:85: [ERROR] CreateChatCompletion error: error, status code: 429, message: Rate limit reached for default-gpt-3.5-turbo-16k in organization org-81U5Eh72Xnow6FD9opokNgmo on tokens per min. Limit: 180000 / min. Current: 176830 / min. Contact us through our help center at help.openai.com if you continue to have issues.
```

#### 解决方案:
1. 失败重试(间隔一定时间) 
2. 因为OpenAI在响应报文中包含total_tokens，可以按每min统计，超过阈值延缓OpenAI并发请求，待时间窗口到期重置计数器

