package unit

// ThousandthInches Inches represents a value in thousandth of an inch Eg. 1.00 inch = 1000 ThousandthInches
type ThousandthInches int

// Int32Ptr returns the int32 representation of an int type.
func (t ThousandthInches) Int32Ptr() *int32 {
	val := int32(t)
	return &val
}
