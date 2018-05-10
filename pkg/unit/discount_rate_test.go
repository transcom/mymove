package unit

import (
	"testing"
)

func TestRateCreateFromPercent(t *testing.T) {
	rate := NewDiscountRateFromPercent(50.5)

	expected := DiscountRate(.505)
	if rate != expected {
		t.Errorf("wrong rate returned: expected %v, got %v", expected, rate)
	}
}

func TestApply(t *testing.T) {
	cents := Cents(1234567) // $12,345.67
	rate := NewDiscountRateFromPercent(65.7)

	result := rate.Apply(cents)

	expected := Cents(423456) // $4,234.56
	if result != expected {
		t.Errorf("wrong number of Cents: expected %d, got %d", expected, result)
	}
}
