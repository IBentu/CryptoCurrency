package main

// abs returns the absolute value for an integer
func abs(value int) int {
	if value < 0 {
		value = -value
	}
	return value
}
