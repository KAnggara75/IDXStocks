package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-04-17", "2024-04-17"},
		{"17-04-2024", "2024-04-17"},
		{"17/04/2024", "2024-04-17"},
		{"17 Apr 2024", "2024-04-17"},
		{"17 April 2024", "2024-04-17"},
		{"17-Apr-2024", "2024-04-17"},
		{"2024 Apr 17", "2024-04-17"},
		{"17 Agu 2024", "2024-08-17"},
		{"17 agustus 2024", "2024-08-17"},
		{"", ""},
		{"random string", "random string"},
		{"17-04", "17-04"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := NormalizeDate(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
