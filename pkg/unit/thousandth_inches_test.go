package unit

import (
	"testing"
)

func Test_ThousandthInches(t *testing.T) {
	// Test int -> int32
	thous := ThousandthInches(1000)
	expected := int32(1000)
	result := *thous.Int32Ptr()
	if result != expected {
		t.Errorf("ThousandthInches did not convert properly: expected %d, got %d", expected, result)
	}
}

func TestToFeet(t *testing.T) {
	thous := ThousandthInches(5 * 12000)
	expected := float64((5 * 12000) / 12000)
	result := thous.ToFeet()
	if result != expected {
		t.Errorf("ThousandthInches did not convert properly to feet: expected %f, got %f", expected, result)
	}
}

func TestToInches(t *testing.T) {
	thous := ThousandthInches(12 * 1000)
	expected := float64((12 * 1000) / 1000)
	result := thous.ToInches()
	if result != expected {
		t.Errorf("ThousandthInches did not convert properly to inches: expected %f, got %f", expected, result)
	}
}
