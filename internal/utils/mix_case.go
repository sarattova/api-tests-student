package utils

import (
	"strings"
	"unicode"
)

func MixCase(text string) string {
	var result strings.Builder

	for index, letter := range []rune(text) {
		if index%2 == 0 {
			result.WriteRune(unicode.ToUpper(letter))
		} else {
			result.WriteRune(unicode.ToLower(letter))
		}
	}

	return result.String()
}
