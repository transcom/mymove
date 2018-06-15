package edi

import (
	"math"
	"strconv"
)

// NxToFloat converts strings with the "numeric" EDI attribute type to float64.
// This is a type with an implied decimal.
// N1 (x = 1): 123 --> 12.3
// N2 (x = 2): 123 --> 1.23
func NxToFloat(s string, x int) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f / math.Pow10(x), nil
}

// FloatToNx converts float64 to the Nx string format
func FloatToNx(n float64, x int) string {
	return strconv.FormatFloat(n*math.Pow10(x), 'f', 0, 64)
}
