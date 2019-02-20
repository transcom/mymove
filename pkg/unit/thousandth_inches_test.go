package unit

import (
	"testing"
)

func TestDimensionToCubicFeet(t *testing.T) {
	smallCrate := 10000
	mediumCrate := 36000
	largeCrate := 60000
	xLarge := 1000000

	var flagtests = []struct {
		length   ThousandthInches
		width    ThousandthInches
		height   ThousandthInches
		expected string
	}{
		{ThousandthInches(25000), ThousandthInches(25000), ThousandthInches(25000), "90422"},
		{ThousandthInches(-25000), ThousandthInches(25000), ThousandthInches(25000), "-90422"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(25000), "90422"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(-25000), "-90422"},
		{ThousandthInches(-25000), ThousandthInches(-25000), ThousandthInches(0), "0"},
		{ThousandthInches(smallCrate), ThousandthInches(smallCrate), ThousandthInches(smallCrate), "5787"},
		{ThousandthInches(smallCrate + 2000), ThousandthInches(smallCrate - 2000), ThousandthInches(smallCrate - 5000), "2777"},
		{ThousandthInches(mediumCrate), ThousandthInches(mediumCrate), ThousandthInches(mediumCrate), "270000"},
		{ThousandthInches(mediumCrate - 6000), ThousandthInches(mediumCrate + 8000), ThousandthInches(mediumCrate - 1000), "267361"},
		{ThousandthInches(largeCrate), ThousandthInches(largeCrate), ThousandthInches(largeCrate), "1250000"},
		{ThousandthInches(largeCrate + 8000), ThousandthInches(largeCrate - 10000), ThousandthInches(largeCrate), "1180555"},
		{ThousandthInches(xLarge), ThousandthInches(xLarge), ThousandthInches(xLarge), "5787037037"},
		{ThousandthInches(xLarge + 60000), ThousandthInches(xLarge - 50000), ThousandthInches(xLarge - 16543), "5731141197"},
	}

	for _, flag := range flagtests {
		result := DimensionToCubicFeet(ThousandthInches(flag.length), ThousandthInches(flag.width), ThousandthInches(flag.height)).String()
		if result != flag.expected {
			t.Errorf("volume calculation failed: expected %s, got %s", flag.expected, result)
		}
	}
}
