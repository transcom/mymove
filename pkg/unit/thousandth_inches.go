package unit

import (
	"fmt"
	"math"
)

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// ToInchString returns a inch string of this value
func (ti ThousandthInches) ToInchString() string {
	d := float64(ti) / 1000.0
	s := fmt.Sprintf("$%.2f", d)
	return s
}

// ToInchFloat returns a inch parsed to float64
func (ti ThousandthInches) ToInchFloat() float64 {
	// rounds to nearest inch
	return math.Round(float64(ti) / 1000)
}

// DimensionToCubicFeet - converts dimensions to cubic feet
func DimensionToCubicFeet(length, width, height ThousandthInches) BaseQuantity {
	var l = float64(length) / 1000
	var w = float64(width) / 1000
	var h = float64(height) / 1000
	var cubicFeet = l * w * h / 1728
	cubicFeetAsBaseQuantity := BaseQuantityFromFloat(cubicFeet)
	return cubicFeetAsBaseQuantity
}
