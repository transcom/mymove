package unit

import (
	"fmt"
)

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// ThousandthIncheToInches - converts from ThousandthInches to inches by dividing the value by 1,000
func ThousandthIncheToInches(thouInch ThousandthInches) float64 {
	return float64(thouInch) / 1000
}

// IntToThousandthInches - converts int to ThousandthInches
func IntToThousandthInches(inches int) ThousandthInches {
	return ThousandthInches(inches * 1000)
}

// DimensionToCubicFeet - converts dimensions to cubic feet
func DimensionToCubicFeet(length, width, height ThousandthInches) (float64, error) {
	if length <= 0 || width <= 0 || height <= 0 {
		errorMessage := fmt.Sprintf("all dimensions length: %d width: %d and height: %d  must be grater than 0", length, width, height)
		return -1, fmt.Errorf(errorMessage)
	}
	const InchesPerCubicFoot = 1728
	var l = ThousandthIncheToInches(length)
	var w = ThousandthIncheToInches(width)
	var h = ThousandthIncheToInches(height)
	cubicFeet := l * w * h / InchesPerCubicFoot
	return cubicFeet, nil
}
