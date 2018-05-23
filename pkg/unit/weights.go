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
func (cwt CWT) ToPounds() Pound {
	return Pound(cwt * 100)
}

// String gives a string representation of CWT
func (cwt CWT) String() string {
	return fmt.Sprintf("%d CWT", int(cwt))
}

// Int returns an integer representation of this weight
func (cwt CWT) Int() int {
	return int(cwt)
}

// ToCWT returns the weight of this in CWT, rounded to the nearest integer
func (pounds Pound) ToCWT() CWT {
	return CWT(math.Round(float64(pounds) / 100.0))
}

// Int returns an integer representation of this weight
func (pounds Pound) Int() int {
	return int(pounds)
}

// Float64 returns a float representation of this weight
func (pounds Pound) Float64() float64 {
	return float64(pounds)
}
