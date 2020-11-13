package random

import (
	"testing"
)

func TestGetRandomInt(t *testing.T) {
	intMax := 4
	for i := 0; i < intMax*20; i++ {
		result, err := GetRandomInt(intMax)
		if err != nil {
			t.Fatalf("could not get random integer: %v", err)
		}

		if result >= intMax {
			t.Errorf("random number %d utside of expected max %d", result, intMax)
		}
	}
}

func TestGetRandomIntAddend(t *testing.T) {
	max := 5
	min := 2
	diff := max - min
	for i := 0; i < max*20; i++ {
		result, err := GetRandomIntAddend(min, max)
		if err != nil {
			t.Fatalf("could not get random integer: %v", err)
		}

		if result >= diff {
			t.Errorf("random number %d outside of expected max %d", result, max)
		}
	}
}
