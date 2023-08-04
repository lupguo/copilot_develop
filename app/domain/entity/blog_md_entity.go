package entity

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hold7techs/go-shim/log"
	"gopkg.in/yaml.v3"
)

// BlogMD Blog文档内容
type BlogMD struct {
	Filepath  string      `json:"filename,omitempty"`
	MDHeader  *YamlHeader `json:"yaml_header"`
	MDContent string      `json:"md_content,omitempty"`
}

// YamlHeader YamlHeader内容
type YamlHeader struct {
	Title       string   `json:"title,omitempty"`
	Date        string   `json:"date,omitempty"`
	Weight      int      `json:"weight,omitempty"`
	Type        string   `json:"type,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Keywords    string   `json:"keywords"`
	Description string   `json:"description,omitempty"`
	Summary     string   `json:"summary,omitempty" json:"summary,omitempty"`
}

// NewBlogMD 通过文件filename 初始化一个Blog MD内容
func NewBlogMD(path string) (*BlogMD, error) {
	// 检测文件是否存在，不存在报错
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist")
	}

	// 读取文件内容
	fileContent, err := os.ReadFile(path)
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
		Filepath:  path,
		MDHeader:  yamlHeader,
		MDContent: match[2],
	}, nil
}
