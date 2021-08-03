package unit

import (
	"testing"
)

func Test_CubicFeetStringConversion(t *testing.T) {
	cubicFeet := CubicFeet(10.0)
	result := cubicFeet.String()

	expected := "10.00"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}
