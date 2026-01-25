package utils

import (
	"strings"
)

func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}
