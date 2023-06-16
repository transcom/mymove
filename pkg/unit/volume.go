package unit

import (
	"fmt"
	"math"
	"math/big"
)

// CubicFeet represents cubic feet
type CubicFeet float64

// CubicThousandthInch represents values in cubic thousandths of an inch
type CubicThousandthInch int

// truncateFloat truncates a float to the given number of decimal places
func truncateFloat(f float64) float64 {
	// use big package to create a new float and get a minimum precision value before rounding would occur
	value := big.NewFloat(f)
	minPrec := value.MinPrec()

	// 52 is the MinPrec return for 2 decimal places, so we're checking if our values is equal to or less than that
	if minPrec <= 52 {
		return f
	}
	// if we have more than 2 decimal places, we need to truncate
	shift := math.Pow(10, float64(2))
	return math.Floor(f*shift) / shift
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
