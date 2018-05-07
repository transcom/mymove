package unit

// Rate represents a percentage, usually represented as a value between 0%
// and 100%.
type Rate float64

// Float64 returns the Rate's value as a float64 with 1 representing one whole unit.
func (r Rate) Float64() float64 {
	return float64(r)
}

// Invert returns 1 - self
func (r Rate) Invert() Rate {
	return Rate(1 - r.Float64())
}

// NewRateFromPercent createDecimals a new Rate using a float64 with 100 representing one whole unit.
func NewRateFromPercent(input float64) Rate {
	return Rate(input / 100.0)
}
