package application

import (
	"regexp"
	"testing"

	"github.com/hold7techs/go-shim/shim"
	"github.com/lupguo/copilot_develop/app/domain/entity"
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
