package utils

import (
	"bytes"
	"slices"
)

// EnsureSize expands the slice b.rows so that rowIndex becomes addressable.
func EnsureSize[T any](ary *[]T, index int) {
	if index < len(*ary) {
		return
	}

	newSize := index + 1

	if newSize > cap(*ary) {
		*ary = slices.Grow(*ary, newSize-len(*ary))
	}

	*ary = (*ary)[:newSize]
}

func MaxValueIndex(array []int) (maxValueIndex int) {
	i := 0
	v := array[i]
	maxValueIndex = i
	for i = 1; i < len(array); i++ {
		if array[i] > v {
			v = array[i]
			maxValueIndex = i
		}
	}
	return
}

/*
JoinRows を完成させて
最終行を含まない場合は全ての行末に newline が必要。
最終行を含む場合：
[
"aaa" // 改行を追加する必要がある
"bbb" // 改行を追加する必要がある
"ccc"
]

[
"aaa" // 改行を追加する必要がある
"bbb" // 改行を追加する必要がある
"ccc" // 改行を追加する必要がある
""
]

[
"aaa" // 改行を追加する必要がある
"bbb" // 改行を追加する必要がある
"" // 改行を追加する必要がある
""
]

*/

// JoinRows concatenates rows into a single byte slice.
// When includeingLastRow is false, a newline is appended to every row.
// When includeingLastRow is true, a newline is appended to every row
// except the last row.
func JoinRows(source [][]byte, newline []byte, hasFinalRow bool) ([]byte, []int, error) {
	// startBytePos[i] is the starting byte position of row i.
	// The last element represents the end position of the final row.
	startBytePos := make([]int, len(source)+1)

	var buffer bytes.Buffer

	for i, row := range source {
		startBytePos[i] = buffer.Len()

		if _, err := buffer.Write(row); err != nil {
			return nil, nil, err
		}

		// When the last row is included, don't append a newline
		// after the final row.
		appendNewLine := !hasFinalRow || i < len(source)-1

		if appendNewLine {
			if _, err := buffer.Write(newline); err != nil {
				return nil, nil, err
			}
		}
	}

	startBytePos[len(source)] = buffer.Len()

	return buffer.Bytes(), startBytePos, nil
}
