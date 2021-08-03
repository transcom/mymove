package unit

import (
	"fmt"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// String converts a CubicFeet value into a string
func (c CubicFeet) String() string {
	return fmt.Sprintf("%.2f", c)
}
