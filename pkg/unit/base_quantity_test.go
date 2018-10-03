package unit

import (
	"testing"
)

func TestStringConversion(t *testing.T) {
	BaseQuantity := BaseQuantity(10000)
	result := BaseQuantity.String()

	expected := "10000"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}

func TestToUnitString(t *testing.T) {
	BaseQuantity := BaseQuantity(19999999)
	result := BaseQuantity.ToUnitFloatString()

	expected := "1999.9999"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}

func TestToUnitFloat64(t *testing.T) {
	BaseQuantity := BaseQuantity(199999)
	result := BaseQuantity.ToUnitFloat()

	expected := float64(19.9999)
	if result != expected {
		t.Errorf("wrong number of BaseQuantity: expected %.4f, got %.4f", expected, result)
	}
}

func TestToUnitInt(t *testing.T) {
	BaseQuantity := BaseQuantity(19999)
	result := BaseQuantity.ToUnitInt()

	expected := 1
	if result != expected {
		t.Errorf("wrong number of BaseQuantity: expected %d, got %d", expected, result)
	}
}

func TestToUnitDollarString(t *testing.T) {
	BaseQuantity := BaseQuantity(19999999)
	result := BaseQuantity.ToUnitDollarString()

	expected := "$1999.99"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}
