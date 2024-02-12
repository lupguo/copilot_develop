package config

import (
	"testing"
)

func TestGetAppRoot(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"t1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAppRoot()
			t.Logf("got app root: %s", got)
		})
	}
}
