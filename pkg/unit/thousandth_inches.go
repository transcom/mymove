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
