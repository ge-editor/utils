package utils

// Find overlap between two ranges
// Range is 2 integers greater than or equal to 0
func FindOverlap(range1Start, range1End, range2Start, range2End int) (int, int) {
	start := max(range1Start, range2Start)
	end := min(range1End, range2End)

	// overlap
	if start <= end {
		return start, end
	}

	// not overlap
	return -1, -1
}

func TabWidth(cursorPosition int, tabWidth int) int {
	tabstop := cursorPosition + 1
	quotient := tabstop / tabWidth
	remainder := tabstop % tabWidth
	tabstop = quotient * tabWidth
	if remainder > 0 {
		tabstop += tabWidth
	}
	return tabstop - cursorPosition
}

func Threshold(maxThreshold, screenSize int) int {
	threshold := (screenSize - 1) / 2
	if threshold < maxThreshold {
		return threshold
	}
	return maxThreshold
}

func Swap[T any](a *T, b *T) {
	temp := *a
	*a = *b
	*b = temp
}
