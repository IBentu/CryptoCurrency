package main

import (
	"time"
)

// abs returns the absolute value for an integer
func abs(value int) int {
	if value < 0 {
		value = -value
	}
	return value
}

// GetCurrentMillis returns the current time in millisecs
func GetCurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
