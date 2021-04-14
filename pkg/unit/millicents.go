package unit

import (
	"fmt"
	"math"
)

// Millicents represents hundredthousandths of US dollars (1000 millicents/ cent)
type Millicents int

// Int64 returns the value of self as an int64
func (m Millicents) Int64() int64 {
	return int64(m)
}

// Int returns the value of self as an int
func (m Millicents) Int() int {
	return int(m)
}

// Float64 returns the value of self as a float64
func (m Millicents) Float64() float64 {
	return float64(m)
}

// MultiplyFloat64 returns the value of self multiplied by multiplier
func (m Millicents) MultiplyFloat64(f float64) Millicents {
	return Millicents(math.Round(float64(m.Int()) * f))
}

// ToCents returns a Cents representation of this value
func (m Millicents) ToCents() Cents {
	return Cents(math.Round(float64(m) / 1000))
}

// ToDollarString returns a dollar string representation of this value
func (m Millicents) ToDollarString() string {
	d := float64(m) / 100000.0
	s := fmt.Sprintf("$%.2f", d)
	return s
}

// ToDollarFloat returns a dollar representation of this value (rounded to nearest 2 decimals)
func (m Millicents) ToDollarFloat() float64 {
	// rounds to nearest cent
	d := math.Round(float64(m) / 1000)
	// convert cents to dollars
	d = d / 100
	return d
}

// ToDollarFloatNoRound returns a dollar representation of this value (no rounding, so partial cents possible)
func (m Millicents) ToDollarFloatNoRound() float64 {
	return float64(m) / 100000.0
}
