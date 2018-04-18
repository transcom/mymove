package unit

import (
	"fmt"
	"math"
)

// CWT represents a value that is a multiple of 100 pounds
type CWT int

// Pound represents a value that is a multiple of 1 pound
type Pound int

// ToPounds returns the weight of this CWT in pounds
func (cwt CWT) ToPounds() (pounds Pound) {
	pounds = Pound(cwt * 100)

	return pounds
}

// ToCWT returns the weight of this in CWT, rounded to the nearest integer
func (pounds Pound) ToCWT() (cwt CWT) {
	cwt = CWT(math.Round(float64(pounds) / 100.0))

	return cwt
}

// String gives a string representation of CWT
func (cwt CWT) String() string {
	return fmt.Sprintf("%d CWT", int(cwt))
}
