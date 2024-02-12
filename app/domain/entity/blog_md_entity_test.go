package entity

import (
	"testing"

	"github.com/hold7techs/go-shim/shim"
)

func TestMinimiseContent(t *testing.T) {
	tests := []struct {
		name        string
		mdContent   string
		mdWordCount int
		expected    string
	}{
		{
			name:        "t1",
			mdContent:   "This is a test\n```\ncode block\n```\nThis is another test",
			mdWordCount: OpenAIMinTokenSize - 100,
			expected:    "This is a test\n\nThis is another test",
		},
		{
			name: "t3",
			mdContent: `This is a test
1. item1: xxx
2. item2: yyy
This is another test`,
			mdWordCount: OpenAIMediumTokenSize - 100,
			expected:    "This is a test\n1. item1\n2. item2\nThis is another test",
		},
		{
			name: "t4",
			mdContent: `This is a test
- List item 1: xxx
- List item 2: yyy
This is another test`,
			mdWordCount: OpenAIMaxTokenSize,
			expected:    "This is a test\nThis is another test",
		},
		{
			name:        "t6",
			mdContent:   "This is a test\n```\ncode block\n```\n- List item 1\n- List item 2\n: Right info\nThis is another test",
			mdWordCount: OpenAIMediumTokenSize,
			expected:    "This is a test\n: Right info\nThis is another test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			miniContent, miniLevel := minimiseContent(tt.mdWordCount, tt.mdContent)
			if miniContent != tt.expected {
				t.Errorf("got %s, but want %s, miniLevel=%v", miniContent, tt.expected, miniLevel)
			}
		})
	}
}
func TestWordsCount(t *testing.T) {
	tests := []struct {
		name          string
		mdContent     string
		expectedCount int
	}{
		{name: "t1", mdContent: "", expectedCount: 0},
		{name: "t2", mdContent: "This is a test", expectedCount: 4},
		{name: "t3", mdContent: "这是一个测试", expectedCount: 6},
		{name: "t4", mdContent: "This is a test 这是一个测试", expectedCount: 10},
		{name: "t5", mdContent: "This is a test\n```\ncode block\n```\nThis is another test", expectedCount: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wordsCount(tt.mdContent)
			if result != tt.expectedCount {
				t.Errorf("got %d, but want %d", result, tt.expectedCount)
			}
		})
	}
}
func TestNewBlogMD(t *testing.T) {
	path := `/private/data/www/tkstorm.com/content/posts/architecture/design-principle/how-to-deal-with-technial-debt.md`
	md, err := NewBlogMD(path)
	if err != nil {
		t.Error(err)
	}

	t.Logf("md=%v", shim.ToJsonString(md, true))
}
