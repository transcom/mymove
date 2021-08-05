package unit

import (
	"fmt"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// CubicThousandthInch represents values in cubic thousandths of an inch
type CubicThousandthInch int

// String converts a CubicFeet value into a string
func (c CubicFeet) String() string {
	return fmt.Sprintf("%.2f", c)
}

// ToCubicFeet converts cubic thousandths of an inch to cubic feet
func (cti CubicThousandthInch) ToCubicFeet() CubicFeet {
	cubicThousandthInchPerCubicFoot := thousandthInchPerFoot * thousandthInchPerFoot * thousandthInchPerFoot
	return CubicFeet(float64(cti) / float64(cubicThousandthInchPerCubicFoot))
}
