package unit

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// DimensionToCubicFeet - converts dimensions to cubic feet
func DimensionToCubicFeet(length, width, height ThousandthInches) BaseQuantity {
	var l = float64(length) / 1000
	var w = float64(width) / 1000
	var h = float64(height) / 1000
	var cubicFeet = l * w * h / 1728
	cubicFeetAsBaseQuantity := BaseQuantityFromFloat(cubicFeet)
	return cubicFeetAsBaseQuantity
}
