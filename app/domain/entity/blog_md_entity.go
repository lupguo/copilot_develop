package entity

import (
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
	// OpenAIMaxTokenSize OpenAI最大token阈值限制
	OpenAIMaxTokenSize = 1600

	// WeightHigh 默认文章权重
	WeightHigh    = 110
	WeightDefault = 100
	WeightLow     = 50
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
	Draft             bool      `yaml:"draft"`
	Keywords          string    `yaml:"keywords,omitempty"`
	Description       string    `yaml:"description,omitempty"`
	Summary           string    `yaml:"summary,omitempty"`
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
	// 名称过滤
	if filepath.Base(path) == "_index.md" {
		return nil, errors.New("_index.md file cannot be replaced")
	}

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

// OverMaxToken 检测MD是否OK
func (md *BlogMD) OverMaxToken() error {
	// token大小检测
	if length := len(md.MDContent); length > OpenAIMaxTokenSize {
		return errors.Errorf("md content[%d] over max token size[%d]", length, OpenAIMaxTokenSize)
	}

	return nil
}

// IsIndexMD 是否Index MD文件
func (md *BlogMD) IsIndexMD() bool {
	// 名称过滤
	if filepath.Base(md.Filepath) == "_index.md" {
		return true
	}

	return false
}

// IsOverMaxTokenSize 是否超过了Token阈值
func (md *BlogMD) IsOverMaxTokenSize() bool {
	// token大小检测
	if length := md.ContentLength(); length > OpenAIMaxTokenSize {
		return true
	}

	return false
}

// IsDraft 是否为草稿文件
func (md *BlogMD) IsDraft() bool {
	if md.MDHeader != nil && md.MDHeader.Draft {
		return true
	}

	// 如果长度小于10 默认也为草稿
	if md.ContentLength() < 10 {
		return true
	}

	return false
}

// ContentLength MD内容长度
func (md *BlogMD) ContentLength() int {
	return len(md.MDContent)
}

func (md *BlogMD) WordCount() int {
	return len(md.MDContent)
}

// CalcArticleWeight 计算新的文章权重
// 文章权重 = 默认文章权重(100) +/- 时间权重 +/- 字数权重(没100长度权重）
func (md *BlogMD) CalcArticleWeight() int {
	// 内容长度
	if md.ContentLength() == 0 {
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
