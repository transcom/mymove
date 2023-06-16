package unit

import (
	"testing"
)

func Test_CubicFeetStringConversion(t *testing.T) {
	cubicFeet := CubicFeet(10.005)
	result := cubicFeet.String()

	expected := "10.00"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}

	cubicFeet = CubicFeet(1117.000000000000000001)
	result = cubicFeet.String()

	expected = "1117.00"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}

	cubicFeet = CubicFeet(9.9)
	result = cubicFeet.String()

	expected = "9.90"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}
	cubicFeet = CubicFeet(9.90000000000000000000000001)
	result = cubicFeet.String()

	expected = "9.90"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}
	cubicFeet = CubicFeet(1.92345)
	result = cubicFeet.String()

	expected = "1.92"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}

	cubicFeet = CubicFeet(1520)
	result = cubicFeet.String()

	expected = "1520.00"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}
	cubicFeet = CubicFeet(1520.123956789012356756856896956533734573585689038905237835)
	result = cubicFeet.String()

	expected = "1520.12"
	if result != expected {
		t.Errorf("wrong string of CubicFeet: expected %s, got %s", expected, result)
	}

	cubicFeet = CubicFeet(523452.55)
	result = cubicFeet.String()

	expected = "523452.55"
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
