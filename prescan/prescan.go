package prescan

import (
	"bytes"
	"io"
	"os"
	"unicode/utf8"
)

const sniffSize = 64 * 1024

type NewlineType int

const (
	LF NewlineType = iota
	CRLF
	CR
	Mixed
)

type FileInfo struct {
	Encoding string
	HasBOM   bool
	Newline  NewlineType
	IsBinary bool
}

func Analyze(path string) (*FileInfo, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	buf := make([]byte, sniffSize)
	n, err := fp.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = buf[:n]

	info := &FileInfo{}

	// --- BOM 判定 ---
	switch {
	case bytes.HasPrefix(buf, []byte{0xEF, 0xBB, 0xBF}):
		info.Encoding = "utf-8"
		info.HasBOM = true
	case bytes.HasPrefix(buf, []byte{0xFF, 0xFE}):
		info.Encoding = "utf-16le"
		info.HasBOM = true
	case bytes.HasPrefix(buf, []byte{0xFE, 0xFF}):
		info.Encoding = "utf-16be"
		info.HasBOM = true
	default:
		info.Encoding = detectUTF8(buf)
	}

	// --- 改行判定 ---
	info.Newline = detectNewline(buf)

	// --- バイナリ判定 ---
	info.IsBinary = detectBinary(buf)

	return info, nil
}

func detectUTF8(b []byte) string {
	if utf8.Valid(b) {
		return "utf-8"
	}
	return "unknown"
}

func detectNewline(b []byte) NewlineType {
	var lf, crlf, cr int

	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			if i > 0 && b[i-1] == '\r' {
				crlf++
			} else {
				lf++
			}
		} else if b[i] == '\r' {
			if i+1 >= len(b) || b[i+1] != '\n' {
				cr++
			}
		}
	}

	max := lf
	result := LF

	if crlf > max {
		max = crlf
		result = CRLF
	}
	if cr > max {
		max = cr
		result = CR
	}

	if (lf > 0 && crlf > 0) ||
		(lf > 0 && cr > 0) ||
		(crlf > 0 && cr > 0) {
		return Mixed
	}

	return result
}

func detectBinary(b []byte) bool {
	for _, c := range b {
		if c == 0 {
			return true
		}
	}
	return false
}
