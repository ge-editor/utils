package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/width"
)

func StringWidth(s string) int {
	w := 0
	for _, ch := range s {
		w += RuneWidth(ch)
	}
	return w
}

// Guess the width of the rune for console screen.
//
// 種類の説明:
//
//   - width.Neutral: 幅が中立的である文字（例：アルファベット、数字、記号）は幅が1です。
//   - width.EastAsianWide: 全角文字（例：漢字、ひらがな、カタカナ）は幅が2です。
//   - width.EastAsianNarrow: 半角文字（例：英数字、一部の記号）は幅が1です。
//   - width.EastAsianAmbiguous: 幅が曖昧な文字（例：一部の漢字、一部のひらがな、一部のカタカナ）は幅が2です。
//   - width.EastAsianFullwidth: 全角文字（例：漢字、ひらがな、カタカナ）は幅が2です。
//   - width.EastAsianHalfwidth: 半角文字（例：英数字、一部の記号）は幅が1です。
func RuneWidth(ch rune) int {
	// return runewidth.RuneWidth(ch)
	p := width.LookupRune(ch)
	switch p.Kind() {
	case width.Neutral:
		return 1
	case width.EastAsianWide:
		// return runewidth.RuneWidth(ch)
		return 2
	case width.EastAsianNarrow:
		return 1
	case width.EastAsianAmbiguous:
		// return runewidth.RuneWidth(ch)
		return 2
	case width.EastAsianFullwidth:
		return 2
	case width.EastAsianHalfwidth:
		return 1
	default:
		return 1
		// return runewidth.RuneWidth(ch)
	}

	// width.EastAsianAmbiguous.String()
	// width.EastAsianAmbiguous == width.Kind().str
	//return utf8.RuneCount([]byte{byte(r)})
	// return utf8.RuneCountInString(string(r))
	// github.com/mattn/go-runewidth
	/*
		switch r {
		case '…', '○', '→', '―':
			return 2
		default:
			return runewidth.RuneWidth(r)
		}
	*/
}

// Determine the character type and return the result string
func WidthKindString(ch rune) string {
	p := width.LookupRune(ch)
	switch p.Kind() {
	case width.Neutral:
		return "Neutral"
	case width.EastAsianWide:
		return "EastAsianWide"
	case width.EastAsianNarrow:
		return "EastAsianNarrow"
	case width.EastAsianAmbiguous:
		return "EastAsianAmbiguous"
	case width.EastAsianFullwidth:
		return "EastAsianFullwidth"
	case width.EastAsianHalfwidth:
		return "EastAsianHalfwidth"
	default:
		return "Unknown"
	}
}

// Check if a string contains a specific character
// case insensitive
func ContainsAllCharacters(str, characters string) bool {
	return ContainsAllCharactersCaseSensitive(strings.ToLower(str), strings.ToLower(characters))
}

// Check if a string contains a specific character
// case sensitive
func ContainsAllCharactersCaseSensitive(str, characters string) bool {
	for _, char := range characters {
		if !strings.ContainsRune(str, char) {
			return false
		}
	}
	return true
}

// Remove half-width symbols from string
func RemoveSymbols(s string) string {
	var result strings.Builder

	for _, r := range s {
		if !unicode.IsSymbol(r) && !unicode.IsPunct(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}
