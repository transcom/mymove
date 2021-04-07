package unit

import (
	"testing"
)

func TestMillicents_Float64(t *testing.T) {
	millicents := Millicents(250000)
	result := millicents.Float64()

	expected := float64(250000)
	if result != expected {
		t.Errorf("wrong number of Millicents: expected %v, got %v", expected, result)
	}
}

func TestMillicents_MultiplyFloat64(t *testing.T) {
	millicents := Millicents(25)
	result := millicents.MultiplyFloat64(5)

	expected := Millicents(125.00)
	if result != expected {
		t.Errorf("wrong number of Millicents: expected %d, got %d", expected, result)
	}
}

func TestMillicents_ToCents(t *testing.T) {
	// expected to round down
	millicents := Millicents(32125)
	result := millicents.ToCents()

	expected := Cents(32)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %v, got %v", expected, result)
	}

	// Expected to round up
	millicents = Millicents(32725)
	result = millicents.ToCents()

	expected = Cents(33)
	if result != expected {
		t.Errorf("wrong number of Cents: expected %v, got %v", expected, result)
	}
}

func TestMillicents_ToDollarString(t *testing.T) {
	millicents := Millicents(32125)
	result := millicents.ToDollarString()

	expected := "$0.32"
	if result != expected {
		t.Errorf("wrong number of dollars: expected %s, got %s", expected, result)
	}
}

func TestMillicents_ToDollarFloat(t *testing.T) {
	// expected to round down
	millicents := Millicents(32125)
	result := millicents.ToDollarFloat()

	expected := float64(0.32)
	if result != expected {
		t.Errorf("wrong number of dollars: expected %v, got %v", expected, result)
	}

	// Expected to round up
	millicents = Millicents(32725)
	result = millicents.ToDollarFloat()

	expected = float64(0.33)
	if result != expected {
		t.Errorf("wrong number of dollars: expected %v, got %v", expected, result)
	}
}

func TestMillicents_ToDollarFloatNoRound(t *testing.T) {
	millicents := Millicents(32125)
	result := millicents.ToDollarFloatNoRound()

	expected := float64(0.32125)
	if result != expected {
		t.Errorf("wrong number of dollars: expected %v, got %v", expected, result)
	}
}
