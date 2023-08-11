package entity

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/hold7techs/go-shim/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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
	Tags              []string  `yaml:"tags,omitempty"`
	Categories        []string  `yaml:"categories,omitempty"`
	Keywords          string    `yaml:"keywords,omitempty"`
	Description       string    `yaml:"description,omitempty"`
	Summary           string    `yaml:"summary,omitempty"`
	SummaryUpdateTime time.Time `yaml:"summary_update_time,omitempty"`
}

var blogMdRegex = regexp.MustCompile("(?sm)^---\n(.*)\n---\n(.*)$")

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
