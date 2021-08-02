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

func Test_CubicFeetFromStringSuccess(t *testing.T) {
	result, err := CubicFeetFromString("11.5")
	if err != nil {
		t.Error(err)
	}

	expected := CubicFeet(11.5)
	if result != expected {
		t.Errorf("Incorrectly parsed CubicFeet string: expected %f, got %f", expected, result)
	}
}

func Test_CubicFeetFromStringFailure(t *testing.T) {
	result, err := CubicFeetFromString("11.5")
	if err != nil {
		t.Error(err)
	}

	expected := CubicFeet(11.5)
	if result != expected {
		t.Errorf("Incorrectly parsed CubicFeet string: expected %f, got %f", expected, result)
	}
}
