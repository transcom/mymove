package unit

import (
	"math"
	"strconv"
)

// Cents represents a value in hundreths of US dollars (aka cents).
type Cents int

// Multiply returns the value of self multiplied by multiplier
func (c Cents) Multiply(i int) Cents {
	return Cents(i * c.Int())
}

// AddCents returns the value of self added with the parameter
func (c Cents) AddCents(a Cents) Cents {
	return Cents(c.Int() + a.Int())
}

// MultiplyFloat64 returns the value of self multiplied by multiplier
func (c Cents) MultiplyFloat64(f float64) Cents {
	return Cents(math.Round(float64(c.Int()) * f))
}

// Multiply returns the value of self multiplied by multiplier
func (c Cents) String() string {
	return strconv.Itoa(int(c))
}

// Int returns the value of self as an int
func (c Cents) Int() int {
	return int(c)
}
