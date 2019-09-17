package unit

import (
	"testing"
)

func TestString(t *testing.T) {
	miles := Miles(2500)
	result := miles.String()

	expected := "2,500"
	if result != expected {
		t.Errorf("wrong number of Miles: expected %s, got %s", expected, result)
	}
}
