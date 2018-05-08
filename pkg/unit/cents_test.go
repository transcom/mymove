package unit

import (
	"testing"
)

func TestCentsMultiply(t *testing.T) {
	cents := Cents(25)
	result := cents.Multiply(5)

	expected := Cents(125)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %d, got %d", expected, result)
	}
}

func TestCentsAddCents(t *testing.T) {
	cents := Cents(25)
	result := cents.AddCents(5)
	expected := Cents(30)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %d, got %d", expected, result)
	}

	cents = Cents(-5)
	result = cents.AddCents(5)
	expected = Cents(0)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %d, got %d", expected, result)
	}

}

func TestCentsMultiplyFloat64(t *testing.T) {
	cents := Cents(2500)
	result := cents.MultiplyFloat64(0.333)

	expected := Cents(833)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %d, got %d", expected, result)
	}
}

func TestCentsToDollarString(t *testing.T) {
	cents := Cents(1)
	result := cents.ToDollarString()
	expected := "$0.01"
	if result != expected {
		t.Errorf("wrong conversion of Cents: expected %s, got %s", expected, result)
	}

	cents = Cents(100)
	result = cents.ToDollarString()
	expected = "$1.00"
	if result != expected {
		t.Errorf("wrong conversion of Cents: expected %s, got %s", expected, result)
	}

	cents = Cents(10099)
	result = cents.ToDollarString()
	expected = "$100.99"
	if result != expected {
		t.Errorf("wrong conversion of Cents: expected %s, got %s", expected, result)
	}
}
