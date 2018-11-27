package unit

import (
	"fmt"
	"math"
	"strconv"
)

// Cents represents a value in hundreths of US dollars (aka cents).
type Cents int

// Millicents represents hundredthousandths of US dollars (1000 millicents/ cent)
type Millicents int

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

func (c Cents) String() string {
	return strconv.Itoa(int(c))
}

// Int returns the value of self as an int
func (c Cents) Int() int {
	return int(c)
}

// Int64 returns the value of self as an int
func (c Cents) Int64() int64 {
	return int64(c)
}

// ToDollarString returns a dollar string representation of this value
func (c Cents) ToDollarString() string {
	d := float64(c) / 100.0
	s := fmt.Sprintf("$%.2f", d)
	return s
}

// ToDollarFloat returns a dollar string representation of this value
func (c Cents) ToDollarFloat() float64 {
	d := float64(c) / 100.0
	return d
}

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
