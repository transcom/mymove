package unit

import (
	"fmt"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// CubicThousandthInch represents values in cubic thousandths of an inch
type CubicThousandthInch int

// truncateFloat truncates a float to 2 decimal places
func truncateFloat(f float64) float64 {
	return float64(int(f*100)) / 100
}

// String converts a CubicFeet value into a string
func (c CubicFeet) String() string {
	// truncate to 2 decimal places
	truncatedValue := truncateFloat(float64(c))
	// convert truncatedValue to a string
	return fmt.Sprintf("%.2f", truncatedValue)
}

// ToCubicFeet converts cubic thousandths of an inch to cubic feet
func (cti CubicThousandthInch) ToCubicFeet() CubicFeet {
	cubicThousandthInchPerCubicFoot := thousandthInchPerFoot * thousandthInchPerFoot * thousandthInchPerFoot
	return CubicFeet(float64(cti) / float64(cubicThousandthInchPerCubicFoot))
}
