package utils

import (
	"strings"
)

func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

// TruncateString truncates a string to a specified length and appends "..." if truncated.
func TruncateString(s string, length int) string {
	runes := []rune(s) // convert string to rune slice (safe for Unicode)
	if len(runes) <= length {
		return s
	}
	return string(runes[:length])
}
