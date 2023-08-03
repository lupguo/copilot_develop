package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/hold7techs/go-shim/log"
	"gopkg.in/yaml.v3"
)

var blogPath = `/private/data/www/tkstorm.com/content/posts`

// BlogMD Blog文档内容
type BlogMD struct {
	Filename   string      `json:"filename,omitempty"`
	YamlHeader *YamlHeader `json:"yaml_header"`
	MDContent  string      `json:"md_content,omitempty"`
}

// YamlHeader YamlHeader内容
type YamlHeader struct {
	Title       string   `json:"title,omitempty"`
	Date        string   `json:"date,omitempty"`
	Description string   `json:"description,omitempty"`
	Weight      int      `json:"weight,omitempty"`
	Type        string   `json:"type,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Summary     string   `json:"summary,omitempty" json:"summary,omitempty"`
}

//
// title: "tcpdump的使用"
// date: 2017-08-21
// description: tcpdump - dump traffic on a network
// weight: 100
// type: posts
// tags:
//   - os
//   - linux
// categories:
//   - network

// NewBlogMD 通过文件filename 初始化一个Blog MD内容
func NewBlogMD(filename string) (*BlogMD, error) {
	// 检测文件是否存在，不存在报错
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist")
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 通过正则提取YamlHeader和MDContent部分
	re := regexp.MustCompile(`(?m)(---\n(.*)\n---)(.*)`)
	match := re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return nil, fmt.Errorf("YamlHeader not found")
	}
	log.Infof("matches: %v", match)

	// 解析YamlHeader
	yamlHeader := &YamlHeader{}
	err = yaml.Unmarshal([]byte(match[1]), yamlHeader)
	if err != nil {
		return nil, err
	}

	// 返回初始的BlogMD实例
	return &BlogMD{
		Filename:   filename,
		YamlHeader: yamlHeader,
		MDContent:  match[2],
	}, nil
}

//

func main() {
	// 查询目录下所有的markdown目录 -> slice内 []*BlogMD
	blogFiles, err := findBlogFiles(blogPath)
	if err != nil {
		panic(err)
	}

	// 通过正则提取md的主题内容
	for _, blogFile := range blogFiles {
		cont, err := getBlogContent(blogFile)
		if err != nil {
			panic(err)
		}

		// cont 提前摘要
		keywords, summary, err := OpenAIContentSummary(cont)
		if err != nil {
			panic(err)
			// return errors.Wrap(err, "OpenAIContentSummary got err")
		}

		// 将keywords, summary 填充到原有的blog文章内
		replaceKeywordsAndSummary(blogFile, keywords, summary)
	}

}

// 将keywords, summary 填充到原有的blog文章内
func replaceKeywordsAndSummary(file string, keywords, summary []byte) {

}

// OpenAIContentSummary 通过OpenAI提前内容摘要
func OpenAIContentSummary(cont []byte) (keywords []byte, summary []byte, err error) {
	time.Sleep(time.Second)
	return []byte(`keywords-mock1, keywords-mock2`), []byte(`summary mock`), nil
}

// 读取文件内容，通过正则筛选内容
func getBlogContent(file string) ([]byte, error) {
	return nil, nil
}

// 查询目录下所有的markdown目录 -> slice内 []*BlogMD
func findBlogFiles(path string) ([]string, error) {

	return nil, nil
}
