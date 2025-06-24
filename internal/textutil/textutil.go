// Package textutil provides utilities for text manipulation.
package textutil

import (
	"strings"
	"unicode"
)

// Truncate truncates the input string to the specified maximum length.
// If the string is longer than maxLen, it will be truncated at the last
// complete word that fits within maxLen and an ellipsis will be appended.
// If maxLen is less than or equal to 0, returns an empty string.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	// Remove any leading/trailing whitespace
	s = strings.TrimSpace(s)

	if len(s) <= maxLen {
		return s
	}

	// Truncate to maxLen
	truncated := s[:maxLen]

	// Find the last space to avoid cutting mid-word
	lastSpace := strings.LastIndexFunc(truncated, unicode.IsSpace)
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return strings.TrimSpace(truncated) + "..."
}
