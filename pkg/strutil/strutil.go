// Package strutil provides small, project-agnostic string helpers.
package strutil

import (
	"strings"
	"unicode"
)

// EqualAlfaNum reports whether a and b are equal once non-alphanumeric runes
// are stripped, comparing case-insensitively.
func EqualAlfaNum(a string, b string) bool {
	var buffers struct{ a, b strings.Builder }
	if len(b) > len(a) {
		a, b = b, a
	}
	var index uint
	for index = 0; index < uint(min(len(a), len(b))); index++ {
		characterA, characterB := rune(a[index]), rune(b[index])
		if unicode.IsLetter(characterA) || unicode.IsNumber(characterA) {
			buffers.a.WriteRune(characterA)
		}
		if unicode.IsLetter(characterB) || unicode.IsNumber(characterB) {
			buffers.b.WriteRune(characterB)
		}
	}

	for index < uint(len(a)) {
		if characterA := rune(a[index]); unicode.IsLetter(characterA) ||
			unicode.IsNumber(characterA) {
			buffers.a.WriteRune(characterA)
		}
		index++
	}

	return strings.EqualFold(buffers.a.String(), buffers.b.String())
}
