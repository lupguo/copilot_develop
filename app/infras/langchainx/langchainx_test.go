package langchainx

import (
	"testing"
)

func Test_langchainSample01(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			langchainSample01()
		})
	}
}
