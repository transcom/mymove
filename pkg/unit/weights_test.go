package unit

import (
	"testing"
)

func Test_PoundsToCWT(t *testing.T) {
	// Test 1050lbs -> 11cwt
	pounds := Pound(1050)
	expected := CWT(11)
	result := pounds.ToCWT()
	if result != expected {
		t.Errorf("pounds did not convert properly: expected %d, got %d", expected, result)
	}

	// Test 5lbs -> 0cwt
	pounds = Pound(10)
	expected = CWT(0)
	result = pounds.ToCWT()
	if result != expected {
		t.Errorf("pounds did not convert properly: expected %d, got %d", expected, result)
	}

	// Test 49lbs -> 0cwt
	pounds = Pound(49)
	expected = CWT(0)
	result = pounds.ToCWT()
	if result != expected {
		t.Errorf("pounds did not convert properly: expected %d, got %d", expected, result)
	}
}

func Test_CWTToPounds(t *testing.T) {
	// Test 20cwt -> 2000lbs
	cwt := CWT(20)
	expected := Pound(2000)
	result := cwt.ToPounds()
	if result != expected {
		t.Errorf("cwt did not convert properly: expected %d, got %d", expected, result)
	}

	// Test 0cwt -> 0lbs
	cwt = CWT(0)
	expected = Pound(0)
	result = cwt.ToPounds()
	if result != expected {
		t.Errorf("cwt did not convert properly: expected %d, got %d", expected, result)
	}
}

func Test_CWTInt(t *testing.T) {
	// Test 10cwt -> 10
	cwt := CWT(10)
	expected := 10
	result := cwt.Int()
	if result != expected {
		t.Errorf("cwt did not convert properly: expected %d, got %d", expected, result)
	}

	// Test -5cwt -> -5
	cwt = CWT(-5)
	expected = -5
	result = cwt.Int()
	if result != expected {
		t.Errorf("cwt did not convert properly: expected %d, got %d", expected, result)
	}
}

func Test_PoundInt(t *testing.T) {
	// Test 10lbs -> 10
	lbs := Pound(10)
	expected := 10
	result := lbs.Int()
	if result != expected {
		t.Errorf("pound did not convert properly: expected %d, got %d", expected, result)
	}

	// Test -5cwt -> -5
	lbs = Pound(-5)
	expected = -5
	result = lbs.Int()
	if result != expected {
		t.Errorf("pound did not convert properly: expected %d, got %d", expected, result)
	}
}
