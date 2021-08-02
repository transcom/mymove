package unit

import (
	"fmt"
	"strconv"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// CubicFeetFromString parses a CubicFeet value from a string
func CubicFeetFromString(s string) (CubicFeet, error) {
	parsed, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return CubicFeet(0), err
	}
	return CubicFeet(parsed), nil
}

// String converts a CubicFeet value into a string
func (c CubicFeet) String() string {
	return fmt.Sprintf("%.2f", c)
}
