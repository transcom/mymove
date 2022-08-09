package route

import (
	"fmt"
	"math"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
)

const (
	minZip5Distance = 10
	maxZip5Distance = 3500
	zip5Range       = 99999

	minZip3Distance = 1
	maxZip3Distance = 100
	zip3Range       = 99
)

type mockDTODZip5DistanceInfo struct{}

// NewMockDTODZip5Distance is a mock (but deterministic) implementation of a DTOD call to calculate mileage
func NewMockDTODZip5Distance() DTODPlannerMileage {
	return &mockDTODZip5DistanceInfo{}
}

// DTODZip5Distance returns a deterministic (but fake) distance between a pickup and destination zip
func (m *mockDTODZip5DistanceInfo) DTODZip5Distance(appCtx appcontext.AppContext, pickupZip string, destinationZip string) (int, error) {
	// Get first 5 digits of zip codes
	if len(pickupZip) < 5 {
		return 0, fmt.Errorf("pickup zip must be at least 5 digits")
	}
	pickupZip5 := pickupZip[0:5]
	pickupZip3 := pickupZip5[0:3]

	if len(destinationZip) < 5 {
		return 0, fmt.Errorf("destination zip must be at least 5 digits")
	}
	destinationZip5 := destinationZip[0:5]
	destinationZip3 := destinationZip5[0:3]

	// Convert zip codes to integers to help with creating a deterministic distance
	pickupZip5Int, err := strconv.Atoi(pickupZip5)
	if err != nil {
		return 0, fmt.Errorf("pickup zip could not be converted to an integer")
	}
	destinationZip5Int, err := strconv.Atoi(destinationZip5)
	if err != nil {
		return 0, fmt.Errorf("destination zip could not be converted to an integer")
	}

	// For zips where the first 3 digits match, pick a deterministic distance between MinZip3Distance and
	// MazZip3Distance miles. Otherwise, pick a deterministic distance between MinZip5Distance and
	// MaxZip5Distance miles.
	var mileage int
	if pickupZip3 == destinationZip3 {
		mileage = mapZipsToDistance(pickupZip5Int, destinationZip5Int, zip3Range, minZip3Distance, maxZip3Distance)
	} else {
		mileage = mapZipsToDistance(pickupZip5Int, destinationZip5Int, zip5Range, minZip5Distance, maxZip5Distance)
	}

	return mileage, nil
}

// Maps the difference between two zips to a number in the min/max range
func mapZipsToDistance(pickupZip int, destinationZip int, zipRange int, min int, max int) int {
	zipDiff := math.Abs(float64(pickupZip - destinationZip))
	targetRange := float64(max - min)
	distance := ((zipDiff * targetRange) / float64(zipRange)) + float64(min)
	return int(math.Floor(distance))
}
