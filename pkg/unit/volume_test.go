package unit

import (
	"testing"
)

func Test_CubicFeetStringConversion(t *testing.T) {
	cubicFeet := CubicFeet(10.0)
	result := cubicFeet.String()

	expected := "10.00"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}
}

func Test_ToCubicFeet(t *testing.T) {
	// This test case covers an implausibly large size for a crate (100 ft cube) to make sure no numerical issues pop up
	feet := 100
	thou := feet * 12 * 1000
	cubicThou := CubicThousandthInch(thou * thou * thou)

	result := cubicThou.ToCubicFeet()
	expected := CubicFeet(float64(feet * feet * feet))

	if result != expected {
		t.Errorf("wrong cubic thousandth inch to cubic feet conversion: expected %s, got %s", expected, result)
	}
}
