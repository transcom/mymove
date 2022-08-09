package route

import (
	"testing"
)

func (suite *GHCTestSuite) TestMockDTODZip5Distance() {
	tests := []struct {
		name             string
		pickupZip        string
		destinationZip   string
		expectedDistance int
		shouldError      bool
		errorMessage     string
	}{
		{"zip5 distance", "30907", "29212", 69, false, ""},
		{"zip5 distance reversed", "29212", "30907", 69, false, ""},
		{"zip5 distance min", "30900", "30899", minZip5Distance, false, ""},
		{"zip5 distance max", "00000", "99999", maxZip5Distance, false, ""},
		{"zip3 distance", "30907", "30901", 7, false, ""},
		{"zip3 distance reversed", "30901", "30907", 7, false, ""},
		{"zip3 distance min", "30907", "30907", minZip3Distance, false, ""},
		{"zip3 distance max", "30900", "30999", maxZip3Distance, false, ""},
		{"too short pickup zip", "3090", "29212", 0, true, "pickup zip must be at least 5 digits"},
		{"too short destination zip", "30907", "2921", 0, true, "destination zip must be at least 5 digits"},
		{"invalid pickup zip", "3090x", "29212", 0, true, "pickup zip could not be converted to an integer"},
		{"invalid destination zip", "30907", "2921x", 0, true, "destination zip could not be converted to an integer"},
	}

	for _, test := range tests {
		suite.T().Run("fake call to DTOD: "+test.name, func(t *testing.T) {
			dtod := NewMockDTODZip5Distance()
			distance, err := dtod.DTODZip5Distance(suite.AppContextForTest(), test.pickupZip, test.destinationZip)

			if test.shouldError {
				suite.Error(err)
				suite.Contains(err.Error(), test.errorMessage)
			} else {
				suite.NoError(err)
			}

			suite.Equal(test.expectedDistance, distance)
		})
	}
}
