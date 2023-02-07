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

	minZip3Distance = 2
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
	// DTOD apparently returns -1 for the distance when there are errors, so we simulate that here
	// (we shouldn't be using the distance value on error, but this at least keeps things in sync).

	// Get first 5 digits of zip codes
	if len(pickupZip) < 5 {
		return -1, fmt.Errorf("pickup zip must be at least 5 digits")
	}
	pickupZip5 := pickupZip[0:5]
	pickupZip3 := pickupZip5[0:3]

	if len(destinationZip) < 5 {
		return -1, fmt.Errorf("destination zip must be at least 5 digits")
	}
	destinationZip5 := destinationZip[0:5]
	destinationZip3 := destinationZip5[0:3]

	// According to https://www.unitedstateszipcodes.org/80901/ this
	// zip is a PO Box only zip. From the slack thread
	// https://ustcdp3.slack.com/archives/C0250Q8469K/p1672933844613089
	if pickupZip5 == "80901" || destinationZip5 == "80901" {
		// in dtod_zip5_distance.go, this is the error if DTOD returns
		// a distance less than or equal to 0. The DTOD service
		// returns 0 for PO Box zips
		return 0, fmt.Errorf("invalid distance using pickup %s and destination %s", pickupZip5, destinationZip5)
	}

	// invalid zip from
	// https://gist.github.com/lsl/98eb26082f71ce5d4f39eb348401b28b
	if pickupZip5 == "11111" || destinationZip5 == "11111" {
		// in dtod_zip5_distance.go, this is the error if DTOD returns
		// a distance less than or equal to 0. The DTOD service
		// returns 0 for invalid zips
		return 0, fmt.Errorf("invalid distance using pickup %s and destination %s", pickupZip5, destinationZip5)
	}

	// Convert zip codes to integers to help with creating a deterministic distance
	pickupZip5Int, err := strconv.Atoi(pickupZip5)
	if err != nil {
		return -1, fmt.Errorf("pickup zip could not be converted to an integer")
	}
	destinationZip5Int, err := strconv.Atoi(destinationZip5)
	if err != nil {
		return -1, fmt.Errorf("destination zip could not be converted to an integer")
	}

	// If the zip5 values are the same, just return a distance of zero like DTOD
	if pickupZip5 == destinationZip5 {
		return 0, nil
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
