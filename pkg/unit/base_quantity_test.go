package unit

import (
	"testing"
)

func TestStringConversion(t *testing.T) {
	baseQuantity := BaseQuantity(10000)
	result := baseQuantity.String()

	expected := "10000"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}

func TestToUnitString(t *testing.T) {
	baseQuantity := BaseQuantity(19999999)
	result := baseQuantity.ToUnitFloatString()

	expected := "1999.9999"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}

func TestToUnitFloat64(t *testing.T) {
	baseQuantity := BaseQuantity(199999)
	result := baseQuantity.ToUnitFloat()

	expected := float64(19.9999)
	if result != expected {
		t.Errorf("wrong number of BaseQuantity: expected %.4f, got %.4f", expected, result)
	}
}

func TestToUnitInt(t *testing.T) {
	baseQuantity := BaseQuantity(19999)
	result := baseQuantity.ToUnitInt()

	expected := 1
	if result != expected {
		t.Errorf("wrong number of BaseQuantity: expected %d, got %d", expected, result)
	}
}

func TestToUnitDollarString(t *testing.T) {
	baseQuantity := BaseQuantity(19999999)
	result := baseQuantity.ToUnitDollarString()

	expected := "$1999.99"
	if result != expected {
		t.Errorf("wrong string of BaseQuantity: expected %s, got %s", expected, result)
	}
}

func TestBaseQuantityFromInt(t *testing.T) {
	result := BaseQuantityFromInt(123)
	expected := BaseQuantity(1230000)
	if result != expected {
		t.Errorf("wrong BaseQuantity for int: expected %d, got %d", expected, result)
	}
}

func TestBaseQuantityFromFloat(t *testing.T) {
	result := BaseQuantityFromFloat(123)
	expected := BaseQuantity(1230000)
	if result != expected {
		t.Errorf("wrong BaseQuantity for float: expected %d, got %d", expected, result)
	}
}

func TestBaseQuantityFromThousandthInches(t *testing.T) {
	result := BaseQuantityFromThousandthInches(123)
	expected := BaseQuantity(1230)
	if result != expected {
		t.Errorf("wrong BaseQuantity for thousandth inches: expected %d, got %d", expected, result)
	}
}
func TestBaseQuantityFromCents(t *testing.T) {
	result := BaseQuantityFromCents(123)
	expected := BaseQuantity(12300)
	if result != expected {
		t.Errorf("wrong BaseQuantity for cents: expected %d, got %d", expected, result)
	}
}

func TestIntToBaseQuantity(t *testing.T) {
	var number int64 = 123
	result := IntToBaseQuantity(&number)
	expected := BaseQuantity(123)

	if *result != expected {
		t.Errorf("wrong BaseQuantity for int64: expected %d, got %v", expected, result)
	}
}

func TestIntToBaseQuantityNilCheck(t *testing.T) {
	nilCheck := IntToBaseQuantity(nil)
	if nilCheck != nil {
		t.Errorf("value found in int64 to IntToBaseQuantity: expected nil")
	}
}
