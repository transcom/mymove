package utils

import (
	"strings"
)

// checks if string is null, empty, or whitespace
func IsNullOrWhiteSpace(s *string) bool {
	return s == nil || len(strings.TrimSpace(*s)) == 0
}
