package ref

import (
	"fmt"
	"regexp"
	"testing"
)

func TestExpand(t *testing.T) {
	content := []byte(`
	# comment line
	option1: value1
	option2: value2

	# another comment line
	option3: value3
`)

	// Regex pattern captures "key: value" pair from the content.
	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)$`)

	// Template to convert "key: value" to "key=value" by
	// referencing the values captured by the regex pattern.
	template := []byte("$key=$value\n")

	result := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(content, -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		result = pattern.Expand(result, template, content, submatches)
	}
	fmt.Println(string(result))
}

func Test_YamlHeaderGet(t *testing.T) {
	str := "---header---content1---content2"

	// (.*?)：使用非贪婪模式匹配任意字符，直到遇到下一个换行符。
	re := regexp.MustCompile(`---(.*?)---(.*)`)
	matches := re.FindStringSubmatch(str)

	if len(matches) >= 2 {
		header := matches[1]
		content := matches[2]
		t.Log("Header:", header)
		t.Log("Content:", content)
	}
}
