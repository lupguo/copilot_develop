package service

import (
	"testing"

	"github.com/lupguo/copilot_develop/app/domain/entity"
)

func Test_minimiseContent(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"t1", args{"Hello```World```"}, "Hello"},
		{"t2", args{"Hello```World```\nSection2 ```Keywords```"}, "Hello\nSection2 "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entity.MinimiseContent(tt.args.content); got != tt.want {
				t.Errorf("minimiseContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
