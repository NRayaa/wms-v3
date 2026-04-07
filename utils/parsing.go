package utils

import (
	"strconv"
)

// ParseIntDefault tries to parse string to int, returns 0 if fail
func ParseIntDefault(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// ParseFloatDefault tries to parse string to float64, returns 0 if fail
func ParseFloatDefault(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
