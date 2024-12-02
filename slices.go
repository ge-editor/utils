package utils

// Add s to slice.
// If s exists in slice, return slice as is.
// Add s to the beginning of slice if addToFront is true.
func AppendIfNotExists(slice []string, s string, addToFront bool) []string {
	if Contains(slice, s) {
		return slice
	}
	if addToFront {
		return append([]string{s}, slice...)
	}
	return append(slice, s)
}

// Check whether s is contained within the slice.
func Contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func MoveElement(slice []string, index int, moveToStart bool) []string {
	if index < 0 || index >= len(slice) {
		// If the index is outside the range of the slice, return the original slice
		return slice
	}

	// Temporarily save the element to be moved
	element := slice[index]

	// Delete the source index
	slice = append(slice[:index], slice[index+1:]...)

	// Determine the destination index and insert the element
	if moveToStart {
		slice = append([]string{element}, slice...)
	} else {
		slice = append(slice, element)
	}

	return slice
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

// CustomAppend は、スライスに新しい要素を追加し、必要に応じて新しいメモリを確保するカスタム関数です。
func CustomAppend[T any](slice []T, elems ...T) []T {
	// 現在のスライスの長さと容量を取得
	currentLen := len(slice)
	currentCap := cap(slice)

	// 追加する要素の数
	requiredLen := currentLen + len(elems)

	// 必要な容量を確保する
	// if requiredLen > currentCap {
	// 新しい容量を決定（倍にするなどの戦略）
	newCap := currentCap * 2
	if newCap < requiredLen {
		newCap = requiredLen
	}

	// 新しいメモリ領域を確保し、既存の要素をコピー
	newSlice := make([]T, requiredLen, newCap)
	// if newSlice == nil {
	// 	panic("failed to allocate memory")
	// }
	copy(newSlice, slice)
	slice = newSlice
	// } else {
	//     // 容量が十分な場合は長さだけ増やす
	//     slice = slice[:requiredLen]
	// }

	// 追加する要素を新しいスライスにコピー
	copy(slice[currentLen:], elems)

	return slice
}
