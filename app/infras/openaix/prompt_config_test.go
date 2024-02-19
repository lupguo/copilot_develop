package openaix

import (
	"testing"

	"github.com/hold7techs/go-shim/shim"
)

func TestParseAppPromptConfig(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t1", args{"prompt.example.yaml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, m := ParseAppPromptConfig(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAppPromptConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("app prompt config: %v", shim.ToJsonString(m, true))
		})
	}
}
