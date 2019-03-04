package unit

import (
	"math"
	"testing"
)

func TestDimensionToCubicFeet(t *testing.T) {
	smallCrate := 10
	mediumCrate := 36
	largeCrate := 60
	xLarge := 100

	var testData = []struct {
		length   ThousandthInches
		width    ThousandthInches
		height   ThousandthInches
		expected float64
	}{
		{IntToThousandthInches(25), IntToThousandthInches(25), IntToThousandthInches(25), 9.04},
		{IntToThousandthInches(-25), IntToThousandthInches(25), IntToThousandthInches(25), -1.00},
		{IntToThousandthInches(smallCrate), IntToThousandthInches(smallCrate), IntToThousandthInches(smallCrate), 0.57},
		{IntToThousandthInches(smallCrate + 2), IntToThousandthInches(smallCrate - 2), IntToThousandthInches(smallCrate - 5), 0.27},
		{IntToThousandthInches(mediumCrate), IntToThousandthInches(mediumCrate), IntToThousandthInches(mediumCrate), 27.00},
		{IntToThousandthInches(mediumCrate - 6), IntToThousandthInches(mediumCrate + 8), IntToThousandthInches(mediumCrate - 1), 26.73},
		{IntToThousandthInches(largeCrate), IntToThousandthInches(largeCrate), IntToThousandthInches(largeCrate), 125.00},
		{IntToThousandthInches(largeCrate + 8), IntToThousandthInches(largeCrate - 10), IntToThousandthInches(largeCrate), 118.05},
		{IntToThousandthInches(xLarge), IntToThousandthInches(xLarge), IntToThousandthInches(xLarge), 578.70},
		{IntToThousandthInches(xLarge) + 16.50*1000, IntToThousandthInches(xLarge) - 15.55*1000, IntToThousandthInches(xLarge) - 27.65*1000, 411.92},
	}

	for _, data := range testData {
		cubicFeet, _ := DimensionToCubicFeet(data.length, data.width, data.height)
		result := cubicFeet
		resultToTwoDecimals := math.Floor(result*100) / 100

		if resultToTwoDecimals != data.expected {
			t.Errorf("volume calculation failed: expected %f, got %.2f", data.expected, resultToTwoDecimals)
		}
	}
}
