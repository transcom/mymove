package unit

import (
	"testing"
)

func TestDimensionToCubicFeet(t *testing.T) {
	smallCrate := 10000
	mediumCrate := 36000
	largeCrate := 60000
	xLarge := 1000000

	var testData = []struct {
		length   ThousandthInches
		width    ThousandthInches
		height   ThousandthInches
		expected string
	}{
		{ThousandthInches(25000), ThousandthInches(25000), ThousandthInches(25000), "90400"},
		{ThousandthInches(-25000), ThousandthInches(25000), ThousandthInches(25000), "-90500"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(25000), "90400"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(-25000), "-90500"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(0), "0"},
		{ThousandthInches(smallCrate), ThousandthInches(smallCrate), ThousandthInches(smallCrate), "5700"},
		{ThousandthInches(smallCrate + 2000), ThousandthInches(smallCrate - 2000), ThousandthInches(smallCrate - 5000), "2700"},
		{ThousandthInches(mediumCrate), ThousandthInches(mediumCrate), ThousandthInches(mediumCrate), "270000"},
		{ThousandthInches(mediumCrate - 6000), ThousandthInches(mediumCrate + 8000), ThousandthInches(mediumCrate - 1000), "267300"},
		{ThousandthInches(largeCrate), ThousandthInches(largeCrate), ThousandthInches(largeCrate), "1250000"},
		{ThousandthInches(largeCrate + 8000), ThousandthInches(largeCrate - 10000), ThousandthInches(largeCrate), "1180500"},
		{ThousandthInches(xLarge), ThousandthInches(xLarge), ThousandthInches(xLarge), "5787037037"},
		{ThousandthInches(xLarge + 60000), ThousandthInches(xLarge - 50000), ThousandthInches(xLarge - 16543), "5731141197"},
	}

	for _, data := range testData {
		cubicFeet := DimensionToCubicFeet(ThousandthInches(data.length), ThousandthInches(data.width), ThousandthInches(data.height))
		result := cubicFeet.String()

		if result != data.expected {
			t.Errorf("volume calculation failed: expected %s, got %s", data.expected, result)
		}
	}
}
