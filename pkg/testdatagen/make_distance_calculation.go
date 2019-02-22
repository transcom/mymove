package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDistanceCalculation creates a single DistanceCalculation
func MakeDistanceCalculation(db *pop.Connection, assertions Assertions) models.DistanceCalculation {
	originAddress := assertions.DistanceCalculation.OriginAddress
	if isZeroUUID(assertions.DistanceCalculation.OriginAddressID) {
		originAddress = MakeAddress(db, assertions)
	}

	destinationAddress := assertions.DistanceCalculation.DestinationAddress
	if isZeroUUID(assertions.DistanceCalculation.DestinationAddressID) {
		destinationAddress = MakeAddress(db, assertions)
	}

	distanceCalculation := models.DistanceCalculation{
		OriginAddress:        originAddress,
		OriginAddressID:      originAddress.ID,
		DestinationAddress:   destinationAddress,
		DestinationAddressID: destinationAddress.ID,
		DistanceMiles:        1044,
	}

	mergeModels(&distanceCalculation, assertions.DistanceCalculation)

	mustCreate(db, &distanceCalculation)

	return distanceCalculation
}

// MakeDefaultDistanceCalculation returns a DistanceCalculation with default values
func MakeDefaultDistanceCalculation(db *pop.Connection) models.DistanceCalculation {
	return MakeDistanceCalculation(db, Assertions{})
}
