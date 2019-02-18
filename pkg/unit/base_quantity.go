package unit

import (
	"fmt"
	"math"
	"strconv"
)

// BaseQuantity represents a value in the 10,000ths of a unit measurement for shipment line items.
// Eg. 10000 BQ = 1.0000 unit or 1 BQ = .0001
// A unit of measurement can be below or more:
// BW = Net Billing Weight
// CF = Cubic Foot
// EA = Each
// FR = Flat Rate
// FP = Fuel Percentage
// NR = Container
// MV = Monetary Value
// TD = Days
// TH = Hours
type BaseQuantity int

// BaseQuantityFromInt creates a BaseQuantity for a provided int
func BaseQuantityFromInt(i int) BaseQuantity {
	return BaseQuantity(i * 10000)
}

// String returns the value of self as string
func (bq BaseQuantity) String() string {
	return strconv.Itoa(int(bq))
}

// ToUnitFloatString returns a unit string representation of this value
// Eg. 10000 BQ -> "1.0000"
func (bq BaseQuantity) ToUnitFloatString() string {
	d := float64(bq) / 10000.0
	s := fmt.Sprintf("%.4f", d)
	return s
}

// ToUnitDollarString returns a dollar string representation of this value
func (bq BaseQuantity) ToUnitDollarString() string {
	// drop the numbers after two decimal points to make sure Sprintf doesn't round up
	// then convert to dollars
	d := math.Floor(float64(bq)/100.0) / 100.0
	s := fmt.Sprintf("$%.2f", d)
	return s
}

// ToUnitInt returns a unit int representation of this value
// Eg. 10000 BQ -> 1
func (bq BaseQuantity) ToUnitInt() int {
	d := float64(bq) / 10000.0
	return int(d)
}

// ToUnitFloat returns a unit float representation of this value
// Eg. 10000 BQ -> 1.0000
func (bq BaseQuantity) ToUnitFloat() float64 {
	d := float64(bq) / 10000.0
	return d
}

// IntToBaseQuantity returns a unit in BaseQuantity format
func IntToBaseQuantity(quantity *int64) *BaseQuantity {
	if quantity == nil {
		return nil
	}
	formattedValue := BaseQuantity(int(*quantity))
	return &formattedValue
}
