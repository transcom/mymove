package unit

import (
	"testing"
)

func TestRateCreateFromPercent(t *testing.T) {
	rate := NewRateFromPercent(50.5)

	expected := Rate(.505)
	if rate != expected {
		t.Errorf("wrong rate returned: expected %v, got %v", expected, rate)
	}
}

func TestRateInvert(t *testing.T) {
	rate := Rate(.75).Invert()

	expected := Rate(.25)
	if rate != expected {
		t.Errorf("wrong rate returned: expected %v, got %v", expected, rate)
	}
}
