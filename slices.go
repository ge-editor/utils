package utils

import (
	"bytes"
	"slices"
)

// Add s to slice.
// If s exists in slice, return slice as is.
// Add s to the beginning of slice if addToFront is true.
func AppendIfNotExists(slice []string, s string, addToFront bool) []string {
	if slices.Contains(slice, s) {
		return slice
	}
	if addToFront {
		return append([]string{s}, slice...)
	}
	return append(slice, s)
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

// JoinBytes concatenates a slice of byte slices into a single byte slice without separators.
// It also calculates the starting byte positions of each row within the final concatenated byte slice.
// The function returns the concatenated byte slice, a slice of starting byte positions of rows, and an error if any.
func JoinBytes(source [][]byte) ([]byte, []int, error) {
	// Initialize a slice to store the starting byte positions of rows; one extra element for the initial zero value.
	startBytePos := make([]int, len(source)+1)
	sum := 0 // Tracks the total bytes written so far.

	// Buffer for efficiently concatenating byte slices.
	var buffer bytes.Buffer

	// Iterate through each byte slice in the source.
	for i, p := range source {
		// Write the current byte slice to the buffer and get the number of bytes written.
		pLen, err := buffer.Write(p)
		if err != nil {
			// Return nil results and the error if writing fails.
			return nil, nil, err
		}

		// Update the total bytes written and record the starting position of the next row.
		sum += pLen
		startBytePos[i+1] = sum
	}

	// Return the concatenated byte slice, the slice of starting byte positions, and no error.
	return buffer.Bytes(), startBytePos, nil
}
