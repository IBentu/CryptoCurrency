package main

import (
	"time"
)

// Abs returns the absolute value for an integer
func Abs(value int) int {
	if value < 0 {
		value = -value
	}
	return value
}

// GetCurrentMillis returns the current time in millisecs
func GetCurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
