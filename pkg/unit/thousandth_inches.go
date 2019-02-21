package unit

import "math"

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// DimensionToCubicFeet - converts dimensions to cubic feet
func DimensionToCubicFeet(length, width, height ThousandthInches) BaseQuantity {
	const cubicInchesPerCubicFoot = 1728
	var l = float64(length) / 1000
	var w = float64(width) / 1000
	var h = float64(height) / 1000
	var cubicFeet = l * w * h / cubicInchesPerCubicFoot
	formatCubicFeet := math.Floor(cubicFeet*100) / 100
	t := BaseQuantityFromFloat(float32(formatCubicFeet))
	return t
}
