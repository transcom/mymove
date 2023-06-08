package unit

import (
	"testing"
)

func TestDollarsToMillicents(t *testing.T) {
	dollars := Dollars(1)
	result := dollars.ToMillicents()

	expected := Millicents(100000)
	if result != expected {
		t.Errorf("wrong number of Miles: expected %v, got %v", expected, result)
	}
}
