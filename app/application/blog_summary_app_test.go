package application

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
)

func TestNewBlogMD(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.BlogMD
		wantErr bool
	}{
		{"t1", args{filename: "./test.md"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blogMd, err := entity.NewBlogMD(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBlogMD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("blogMD: %+v", blogMd)
		})
	}
}

func TestYamlHead(t *testing.T) {
	content := `---
title: "tcpdump的使用"
date: 2017-08-21
description: tcpdump - dump traffic on a network
weight: 100
type: posts
tags:
  - os
  - linux
categories:
  - network
---

>  tcpdump - dump traffic on a network
参考: https://danielmiessler.com/study/tcpdump/
## 1. tcpdump概述

...
`
	// 通过正则提取YamlHeader和MDContent部分
	re := regexp.MustCompile(`(?ms)(?:-{3}(.*)-{3})(.*)`)
	matchs := re.FindStringSubmatch(content)
	for i, s := range matchs {
		if i == 0 {
			continue
		}
		t.Logf("%d=>%+v", i, s)
	}

	if len(matchs) < 2 {
		t.Fatal("matches len less 2")
	}

	head, body := matchs[1], matchs[2]

	// yaml 解析head
	type YamlHeader struct {
		Title       string   `yaml:"title"`
		Date        string   `yaml:"date"`
		Description string   `yaml:"description"`
		Weight      int      `yaml:"weight"`
		Type        string   `yaml:"type"`
		Tags        []string `yaml:"tags"`
		Categories  []string `yaml:"categories"`
	}

	yamlHead := YamlHeader{}
	err := yaml.Unmarshal([]byte(head), &yamlHead)
	if err != nil {
		t.Fatalf("yaml.Unmarshal got err: %s", err)
	}
	t.Logf("yaml: %+v", shim.ToJsonString(yamlHead, true))

	// body -> api
	_ = body
}

func Test_truncateFileAndWriteNewContent(t *testing.T) {
	type args struct {
		md           *entity.BlogMD
		newYamHeader []byte
	}
	tempMdFile := filepath.Join(os.TempDir(), "01.md")
	err := os.WriteFile(tempMdFile, []byte("temp md content"), 0644)
	if err != nil {
		t.Fatal("write temp file got err", err)
	}
	t.Logf("tempMdDir: %s,  tempMdFile: %s", os.TempDir(), tempMdFile)
	defer os.Remove(tempMdFile)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"c1", args{
			md: &entity.BlogMD{
				Filepath:  tempMdFile,
				MDContent: "Blog Content",
			},
			newYamHeader: []byte("Title: NewTitle"),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeNewYamlHeader(tt.args.md, tt.args.newYamHeader); (err != nil) != tt.wantErr {
				t.Errorf("writeNewYamlHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTruncateTempFile(t *testing.T) {
	err := os.Truncate("/var/folders/n4/2qd9ttqd0b927rbmvy6fn7980000gn/T/test_md3347475352/01.md", 0)
	if err != nil {
		t.Errorf("os.turncate got err: %s", err)
	}
}

func TestOpenFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "exist_*.md")
	assert.NoError(t, err)
	t.Log("tempFile", tempFile.Name())

	// 关闭和清理
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// 打开一个不存在的文件
	noExistFileName := fmt.Sprintf("%s_%s", tempFile.Name(), "not-exist")
	_, err = os.Open(noExistFileName)
	if err != nil {
		t.Logf("os.Open not exist file[%s] got err: %s", noExistFileName, err)
	}
	assert.Error(t, err)
}

// mock出一个srv
type mockAISrv struct {
	mock.Mock
}

func (m *mockAISrv) BlogSummary(ctx context.Context, content string) (summary *entity.BlogSummary, err error) {
	args := m.Called(ctx, content)
	return args[0].(*entity.BlogSummary), args.Error(1)
}

func TestBlogSummaryApp_ReplaceKeywordsAndSummary(t *testing.T) {
	// 创建并打开一个临时文件
	tempFile := filepath.Join(os.TempDir(), "01.md")
	file, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	// 模拟写入一些md内容
	fmt.Fprint(file, `---
title: 苹果Wiki
keywords: Apple,智能手机,乔布斯
---
Phone 是苹果公司生产的一系列智能手机，使用苹果自己的 iOS 移动操作系统。第一代iPhone由时任苹果CEO史蒂夫·乔布斯于2007年1月9日发布。此后，苹果每年都会发布新的iPhone型号和iOS更新。截至 2018 年 11 月 1 日，iPhone 销量已超过 22 亿部。截至2022年，iPhone占据全球智能手机市场份额的15.6%。

iPhone 是第一款使用多点触控技术的手机。 [4] 自 iPhone 推出以来，它获得了更大的屏幕尺寸、视频录制、防水和许多辅助功能。在 iPhone 8 和 8 Plus 之前，iPhone 的前面板上只有一个带有 Touch ID 指纹传感器的按钮。自 iPhone X 以来，iPhone 机型已改用近乎无边框的前屏设计，配备 Face ID 面部识别功能，并可通过手势激活应用程序切换。 Touch ID 仍用于廉价 iPhone SE 系列。

iPhone 是与 Android 并列的世界上最大的两个智能手机平台之一，并且在奢侈品市场中占有很大的份额。 iPhone为苹果公司带来了巨额利润，使其成为全球最有价值的上市公司之一。第一代iPhone被形容为手机行业的一场“革命”，后续机型也获得好评。 [5] iPhone 被誉为普及了智能手机和平板电脑，并为智能手机应用程序（或“应用程序经济”）创造了一个巨大的市场。截至 2017 年 1 月，Apple 的 App Store 包含超过 220 万个 iPhone 应用程序。
`)
	file.Close()

	//
	mockAISrv := new(mockAISrv)
	ctx := context.Background()
	mockAISrv.On("BlogSummary", ctx, mock.Anything).Return(&entity.BlogSummary{
		Keywords:    "Mock Keyword1, Mock Keyword2",
		Summary:     "Mock summary...",
		Description: "Mock Description...",
	}, nil)

	type args struct {
		ctx          context.Context
		blogFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t1", args{
			ctx:          ctx,
			blogFilePath: tempFile,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &BlogSummaryApp{
				srv: mockAISrv,
			}
			if err := app.ReplaceKeywordsAndSummary(tt.args.ctx, tt.args.blogFilePath); (err != nil) != tt.wantErr {
				t.Errorf("ReplaceKeywordsAndSummary() error = %v, wantErr %v", err, tt.wantErr)
			}

			c, err := os.ReadFile(tempFile)
			if err != nil {
				t.Errorf("open temp file to assert got err: %s", err)
			}

			t.Logf("content: %s", c)
		})
	}
}
