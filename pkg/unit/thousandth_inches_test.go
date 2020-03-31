package unit

import (
	"testing"
)

func Test_ThousandthInches(t *testing.T) {
	// Test int -> int64
	thous := ThousandthInches(1000)
	expected := int32(1000)
	result := *thous.Int32Ptr()
	if result != expected {
		t.Errorf("ThousandthInches did not convert properly: expected %d, got %d", expected, result)
	}
}
