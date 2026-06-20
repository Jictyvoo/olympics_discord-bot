// Package strutil provides small, project-agnostic string helpers.
package strutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// diacriticFolder decomposes runes and drops the nonspacing marks, folding
// accented letters to their base form (e.g. "fútbol" -> "futbol").
var diacriticFolder = transform.Chain(
	norm.NFD,
	runes.Remove(runes.In(unicode.Mn)),
	norm.NFC,
)

// FoldDiacritics strips diacritics from s, folding accented letters to their
// base form. Casing is preserved; lowercase at the call site if needed.
func FoldDiacritics(s string) string {
	folded, _, err := transform.String(diacriticFolder, s)
	if err != nil {
		return s
	}
	return folded
}

// NormalizeAlfaNum folds diacritics, drops non-alphanumeric runes and lowercases s
// (e.g. "Volley-ball" -> "volleyball").
func NormalizeAlfaNum(s string) string {
	var b strings.Builder
	for _, r := range FoldDiacritics(s) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

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
