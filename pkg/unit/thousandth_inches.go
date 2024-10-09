package unit

const (
	thousandthInchPerFoot = 12000
	thousandthInchPerInch = 1000
)

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// Int32Ptr returns the int32 representation of an int type.
func (t ThousandthInches) Int32Ptr() *int32 {
	// #nosec G115: it is unrealistic that an imperial measurement will exceed int32 limits
	val := int32(t)
	return &val
}

// ToFeet returns feet for this value
func (t ThousandthInches) ToFeet() float64 {
	feet := float64(t) / thousandthInchPerFoot

	return feet
}

// ToInches returns inches for this value
func (t ThousandthInches) ToInches() float64 {
	inches := float64(t) / thousandthInchPerInch

	return inches
}
