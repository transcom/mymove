package unit

import (
	"fmt"
	"math"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// CubicThousandthInch represents values in cubic thousandths of an inch
type CubicThousandthInch int

// truncateFloat truncates a float to the given number of decimal places
func truncateFloat(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift) / shift
}

// String converts a CubicFeet value into a string
func (c CubicFeet) String() string {
	// truncate to 2 decimal places
	truncatedValue := truncateFloat(float64(c), 2)
	// convert truncatedValue to a string
	return fmt.Sprintf("%.2f", truncatedValue)
}

// ToCubicFeet converts cubic thousandths of an inch to cubic feet
func (cti CubicThousandthInch) ToCubicFeet() CubicFeet {
	cubicThousandthInchPerCubicFoot := thousandthInchPerFoot * thousandthInchPerFoot * thousandthInchPerFoot
	return CubicFeet(float64(cti) / float64(cubicThousandthInchPerCubicFoot))
}
