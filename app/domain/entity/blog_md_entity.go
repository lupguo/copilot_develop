package entity

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// OpenAIMinTokenSize 最小请求OpenAI的rune字符大小，太少没有必要请求OpneAI
	OpenAIMinTokenSize = 500
	// OpenAIMaxTokenSize OpenAI最大token阈值限制(按16k的3/4估算)
	OpenAIMaxTokenSize = 12000

	// WeightHigh 默认文章权重
	WeightHigh    = 50
	WeightDefault = 100
	WeightLow     = 200

	// ArticleDraftMinLength 手稿文章字符判断
	ArticleDraftMinLength = 50
)

// BlogMD Blog文档内容
type BlogMD struct {
	Filepath  string      `json:"filename,omitempty"`
	MDHeader  *YamlHeader `json:"yaml_header,omitempty"`
	MDContent string      `json:"md_content,omitempty"`
}

// YamlHeader YamlHeader内容
type YamlHeader struct {
	Title             string    `yaml:"title,omitempty"`
	Date              string    `yaml:"date,omitempty"`
	Weight            int       `yaml:"weight,omitempty"`
	Type              string    `yaml:"type,omitempty"`
	Categories        []string  `yaml:"categories,omitempty"`
	Tags              []string  `yaml:"tags,omitempty"`
	Draft             bool      `yaml:"draft"`                  // 是否手稿
	Keywords          string    `yaml:"keywords,omitempty"`     // 文章关键字
	Description       string    `yaml:"description,omitempty"`  // 文章描述
	Summary           string    `yaml:"summary,omitempty"`      // 文章摘要
	WordCounts        int       `yaml:"words_counts,omitempty"` // 文件字数统计
	ShortMark         string    `yaml:"short_mark,omitempty"`   // 文章短标记
	Aliases           []string  `yaml:"aliases,omitempty"`
	SummaryUpdateTime time.Time `yaml:"summary_update_time,omitempty"`
}

func (y *YamlHeader) String() string {
	marshal, err := json.Marshal(y)
	if err != nil {
		return err.Error()
	}

	return string(marshal)
}

var blogMdRegex = regexp.MustCompile("(?sm)^---\n(.*?)\n---(?:\n+)(.*)$")

// NewBlogMD 通过文件filename 初始化一个Blog MD内容
func NewBlogMD(path string) (*BlogMD, error) {
	// 读取文件内容
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read md file[%s] got err", path)
	}

	// 通过正则提取YamlHeader和MDContent部分
	match := blogMdRegex.FindStringSubmatch(string(fileContent))
	if len(match) != 3 {
		return nil, fmt.Errorf("blog content yaml header not found")
	}
	// log.Infof("matches: %v", match)

	// 解析YamlHeader
	yamlHeader := &YamlHeader{}
	err = yaml.Unmarshal([]byte(match[1]), yamlHeader)
	if err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal got err")
	}

	// 返回初始的BlogMD实例
	md := &BlogMD{
		Filepath:  path,
		MDHeader:  yamlHeader,
		MDContent: match[2],
	}

	return md, nil
}

// IsIndexMD 是否Index MD文件
func (md *BlogMD) IsIndexMD() bool {
	return filepath.Base(md.Filepath) == "_index.md"
}

// IsOverMaxTokenSize 是否超过了OpenAI的Token阈值
func (md *BlogMD) IsOverMaxTokenSize() bool {
	return md.MinimiseContentLength() > OpenAIMaxTokenSize
}

// IsContentTooSmall 内容太少了，不去请求OpenAI
func (md *BlogMD) IsContentTooSmall() bool {
	return md.MinimiseContentLength() < ArticleDraftMinLength
}

// IsDraft 是否为草稿文件
func (md *BlogMD) IsDraft() bool {
	if md.MDHeader != nil && md.MDHeader.Draft {
		return true
	}

	// 如果长度小于10 默认也为草稿
	if md.WordCount() < ArticleDraftMinLength {
		return true
	}

	return false
}

// ShortMark 获取MD的shortMark短标记
func (md *BlogMD) ShortMark() string {
	header := md.MDHeader
	if header.ShortMark == "" && header.Title != "" {
		header.ShortMark = fmt.Sprintf("%x", md5.Sum([]byte(header.Title)))
	}

	return header.ShortMark
}

// MDCodeRegex markdown中的代码正则
var MDCodeRegex = regexp.MustCompile("(?ms)```.*```")

// MinimiseContent 返回处理后的最小化内容
// 通过正则替换content内的代码，降低Token使用量
// 1. 剔除```符号```内的内容
func (md *BlogMD) MinimiseContent() string {
	cont := MDCodeRegex.ReplaceAllString(md.MDContent, "")
	if len([]rune(cont)) < OpenAIMinTokenSize {
		return md.MDContent
	}

	return cont
}

// RawContentLength 原始内容字节长度
func (md *BlogMD) RawContentLength() int {
	return len(md.MDContent)
}

// MinimiseContentLength 内容通过精简替换后的内容长度
func (md *BlogMD) MinimiseContentLength() int {
	return len(md.MinimiseContent())
}

// WordCount 原始内容字符统计（含代码内容部分）
func (md *BlogMD) WordCount() int {
	return len([]rune(md.MDContent))
}

// CalcArticleWeight 计算新的文章权重
// 文章权重 = 默认文章权重(100) +/- 时间权重 +/- 字数权重(没100长度权重）
func (md *BlogMD) CalcArticleWeight() int {

	// 内容长度
	if md.MinimiseContentLength() < ArticleDraftMinLength {
		return WeightLow
	}

	return WeightDefault
}

// ReplaceWithNewYamlHeader 更新成新的MD信息
func (md *BlogMD) ReplaceWithNewYamlHeader() error {
	// 虚拟化处理
	headerStr, err := yaml.Marshal(md.MDHeader)
	if err != nil {
		return errors.Wrapf(err, "marsh file[%s] yaml head got err", md.Filepath)
	}
	// log.Debugf("newMDHeaderStr: %s", headerStr)

	// 清理+重写如新的文件内容
	mdFile, err := os.OpenFile(md.Filepath, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "open blog file[%s] got err", md.Filepath)
	}
	defer mdFile.Close()

	// 重写如file
	if _, err = fmt.Fprintf(mdFile, "---\n%s---\n\n%s", headerStr, md.MDContent); err != nil {
		return errors.Wrapf(err, "write into blog file[%s] with new yaml header got err", md.Filepath)
	}
	return nil
}
