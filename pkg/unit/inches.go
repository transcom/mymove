package unit

// Inch represents a value that is a multiple of 1 inch
type Inch int

// InchfromInt returns a value multiplied by 100 so we can represent as an int but still support 2 decimal point precision
func InchfromInt(i int) Inch {
	return Inch(i * 100)
}
