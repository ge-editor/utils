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
