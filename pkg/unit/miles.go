package unit

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Miles represents mile value in int
type Miles int

// String gives a string representation of Miles
func (miles Miles) String() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", int(miles))
}

// Int gives an int representation of Miles
func (miles Miles) Int() int {
	return int(miles)
}

// Float64 gives a float representation of Miles
func (miles Miles) Float64() float64 {
	return float64(miles)
}
