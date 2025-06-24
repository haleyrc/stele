package textutil_test

import (
	"testing"

	"github.com/haleyrc/assert"
	"github.com/haleyrc/stele/internal/textutil"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short text",
			input:    "Hello world",
			maxLen:   50,
			expected: "Hello world",
		},
		{
			name:     "exact length",
			input:    "Hello world",
			maxLen:   11,
			expected: "Hello world",
		},
		{
			name:     "truncate at word boundary",
			input:    "The quick brown fox jumps over the lazy dog",
			maxLen:   20,
			expected: "The quick brown fox...",
		},
		{
			name:     "truncate with no spaces",
			input:    "abcdefghijklmnopqrstuvwxyz",
			maxLen:   10,
			expected: "abcdefghij...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "zero max length",
			input:    "Hello world",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "negative max length",
			input:    "Hello world",
			maxLen:   -5,
			expected: "",
		},
		{
			name:     "text with leading/trailing whitespace",
			input:    "  Hello world  ",
			maxLen:   20,
			expected: "Hello world",
		},
		{
			name:     "truncate strips trailing whitespace before ellipsis",
			input:    "Hello world this is a test",
			maxLen:   12,
			expected: "Hello world...",
		},
		{
			name:     "150 char truncation",
			input:    "This is a very long description that contains more than one hundred and fifty characters and should be truncated at the last complete word before reaching the maximum length limit.",
			maxLen:   150,
			expected: "This is a very long description that contains more than one hundred and fifty characters and should be truncated at the last complete word before...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textutil.Truncate(tt.input, tt.maxLen)
			assert.Equal(t, "truncated text", tt.expected, result)
		})
	}
}
