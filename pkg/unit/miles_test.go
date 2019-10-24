package unit

import (
	"testing"
)

func TestMilesString(t *testing.T) {
	miles := Miles(2500)
	result := miles.String()

	expected := "2,500"
	if result != expected {
		t.Errorf("wrong number of Miles: expected %s, got %s", expected, result)
	}
}

func TestMilesInt(t *testing.T) {
	miles := Miles(2500)
	result := miles.Int()

	expected := 2500
	if result != expected {
		t.Errorf("wrong number of Miles: expected %v, got %v", expected, result)
	}
}

func TestMilesFloat(t *testing.T) {
	miles := Miles(2500)
	result := miles.Float64()

	expected := float64(2500)
	if result != expected {
		t.Errorf("wrong number of Miles: expected %v, got %v", expected, result)
	}
}
