package unit

// Dollars represents dollar value in float64
type Dollars float64

// ToMillicents converts dollars to millicents
func (d Dollars) ToMillicents() Millicents {
	return Millicents(d * 100 * 1000)
}
