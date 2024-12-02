package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"
)

func RuneToBytes(ch rune) []byte {
	var buf [utf8.UTFMax]byte        // 最大4バイトのバッファを定義
	n := utf8.EncodeRune(buf[:], ch) // バッファのスライスとしてエンコード
	return buf[:n]                   // エンコードされたバイト列を返す
}

// RunesToBytes は、ルーンのスライスをバイトスライスに変換します。
func RunesToBytes(runes []rune) []byte {
	buf := make([]byte, 0, len(runes)*utf8.UTFMax)
	for _, r := range runes {
		n := utf8.EncodeRune(buf[len(buf):cap(buf)], r)
		buf = buf[:len(buf)+n]
	}
	return buf
}

// RunesToBytes は、ルーンのスライスをバイトスライスに変換します。
/*
func RunesToBytes(runes []rune) []byte {
	var buf []byte
	for _, r := range runes {
		buf = append(buf, []byte(string(r))...)
	}
	return buf
}
*/

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

/*
func ReverseUTF8Bytes(bytes []byte) []byte {
	r := []rune(string(bytes))
	slices.Reverse(r)
	return RunesToBytes(r)
}
*/

// ReverseUTF8Bytes は、バイトスライスを逆順にする関数です。
func ReverseUTF8Bytes(bytes []byte) []byte {
	j := len(bytes)
	results := make([]byte, j)
	for i := 0; i < len(bytes); {
		_, size := utf8.DecodeRune(bytes[i:])
		copy(results[j-size:], bytes[i:i+size])
		j -= size
		i += size
	}
	return results
}
