package entity

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hold7techs/go-shim/shim"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// OpenAIMinTokenSize 最小请求OpenAI的rune字符大小，太少没有必要请求OpneAI
	OpenAIMinTokenSize = 1000

	// OpenAIMediumTokenSize 中等长度内容，
	OpenAIMediumTokenSize = 5000

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
	MiniData  *MiniData   `json:"mini_data"` // 精简内容
}

// MiniData 精简后的内容, 参考: https://platform.openai.com/tokenizer
type MiniData struct {
	MiniContent   string `json:"mini_content"`   // 精简化后的内容
	MinLevel      int    `json:"min_level"`      // 精简化后的等级
	MinWordsCount int    `json:"min_word_count"` // 精简后内容
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

// --- head yaml --- ... content //
var blogMdRegex = regexp.MustCompile("(?sm)^---\n(.*?)\n---\n+(.*)$")

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

	// 解析YamlHeader
	header := &YamlHeader{}
	err = yaml.Unmarshal([]byte(match[1]), header)
	if err != nil {
		return nil, errors.Wrap(err, "yaml unmarshal got err")
	}

	// 返回初始的BlogMD实例
	md := &BlogMD{
		Filepath:  path,
		MDHeader:  header,
		MDContent: match[2],
	}

	// MD信息更新
	header.WordCounts = wordsCount(md.MDContent)
	header.Draft = md.IsDraft()                                                           // 是否手稿
	header.Weight = md.CalcArticleWeight()                                                // 文章权重
	header.ShortMark = md.ShortMark()                                                     // 文章签名
	header.Categories = shim.ProcessStringsSlice(header.Categories, nil, strings.ToLower) // 文章分类统一转小写
	header.Tags = shim.ProcessStringsSlice(header.Tags, nil, strings.ToLower)             // 文章标签统一转小写

	// 精简token size
	minContent, minLevel := minimiseContent(header.WordCounts, md.MDContent)
	miniData := &MiniData{
		MiniContent:   minContent,
		MinLevel:      minLevel,
		MinWordsCount: wordsCount(minContent),
	}
	md.MiniData = miniData

	return md, nil
}

// IsIndexMD 是否Index MD文件
func (md *BlogMD) IsIndexMD() bool {
	return filepath.Base(md.Filepath) == "_index.md"
}

// IsMinContentTooLong 文章内容单词太多， 请求OpenAI大概率也行不通
func (md *BlogMD) IsMinContentTooLong(limitSize int) bool {
	return md.MiniData.MinWordsCount > limitSize
}

// IsContentWordsTooSmall 内容单词太少，不去请求OpenAI
func (md *BlogMD) IsContentWordsTooSmall() bool {
	return md.MDHeader.WordCounts < ArticleDraftMinLength
}

// IsDraft 是否为草稿文件
func (md *BlogMD) IsDraft() bool {
	// 已标记为草稿
	if md.MDHeader.Draft == true {
		return true
	}

	// 如果长度小于draft min length 长度的， 也当为草稿
	if md.MDHeader.WordCounts < ArticleDraftMinLength {
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

// mdCodeRegex markdown中的代码正则
var mdCodeRegex = regexp.MustCompile("(?ms)```.*```")
var mdListRightRegex = regexp.MustCompile(`[:：].+`)
var mdListAllRegex = regexp.MustCompile(`[-1-9].+`)
var mdReturnRegex = regexp.MustCompile(`\n+`)

// minimiseContent 通过正则替换，将Blog内容缩小，降低OpenAI的Token使用量
//  1. 剔除```符号```内的内容
//  2. 删除List列表后内容，即仅保留List列表的头部
func minimiseContent(wordsCount int, content string) (miniContent string, miniLevel int) {
	// 基于MD的原始内容长度判断
	switch {
	case wordsCount < OpenAIMinTokenSize: // 小于1000，移除code代码
		contRemoveCode := mdCodeRegex.ReplaceAllString(content, "")
		return contRemoveCode, 0
	case wordsCount < OpenAIMediumTokenSize: // 小于5000, 移除code、list右侧内容
		contRemoveCode := mdCodeRegex.ReplaceAllString(content, "")
		contRemoveListRight := mdListRightRegex.ReplaceAllString(contRemoveCode, "")
		contRemoveReturn := mdReturnRegex.ReplaceAllString(contRemoveListRight, "\n")
		return contRemoveReturn, 1
	default: // 移除code+list全部内容
		contRemoveCode := mdCodeRegex.ReplaceAllString(content, "")
		contRemoveCodeAndList := mdListAllRegex.ReplaceAllString(contRemoveCode, "")
		contRemoveReturn := mdReturnRegex.ReplaceAllString(contRemoveCodeAndList, "\n")
		return contRemoveReturn, 2
	}
}

// 使用正则表达式匹配单词
var wordsRegex = regexp.MustCompile(`(\p{Han}|\b\w+\b)`)

// WordsCount 原始内容字符统计（含代码内容部分）
func wordsCount(content string) int {
	matches := wordsRegex.FindAllString(content, -1)
	return len(matches)
}

// CalcArticleWeight 计算新的文章权重
// 文章权重 = 默认文章权重(100) +/- 时间权重 +/- 字数权重(没100长度权重）
func (md *BlogMD) CalcArticleWeight() int {
	// 内容长度
	if md.MDHeader.WordCounts < ArticleDraftMinLength {
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
