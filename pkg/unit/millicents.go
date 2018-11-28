package unit

import (
	"fmt"
	"math"
)

// ToMillicents converts cents to millicents
func (c Cents) ToMillicents() Millicents {
	return Millicents(c.Int() * 1000)
}

// Int returns the value of self as an int
func (c Millicents) Int() int {
	return int(c)
}

// MultiplyFloat64 returns the value of self multiplied by multiplier
func (c Millicents) MultiplyFloat64(f float64) Millicents {
	return Millicents(math.Round(float64(c.Int()) * f))
}

// ToDollarString returns a dollar string representation of this value
func (c Millicents) ToDollarString() string {
	d := float64(c) / 100000.0
	s := fmt.Sprintf("$%.2f", d)
	return s
}
