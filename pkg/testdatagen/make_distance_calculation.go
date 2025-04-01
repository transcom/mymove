package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDistanceCalculation creates a single DistanceCalculation
func MakeDistanceCalculation(db *pop.Connection, assertions Assertions) (models.DistanceCalculation, error) {
	originAddress := assertions.DistanceCalculation.OriginAddress
	if isZeroUUID(assertions.DistanceCalculation.OriginAddressID) {
		var err error
		originAddress, err = MakeAddress(db, assertions)
		if err != nil {
			return models.DistanceCalculation{}, err
		}
	}

	destinationAddress := assertions.DistanceCalculation.DestinationAddress
	if isZeroUUID(assertions.DistanceCalculation.DestinationAddressID) {
		var err error
		destinationAddress, err = MakeAddress(db, assertions)
		if err != nil {
			return models.DistanceCalculation{}, err
		}
	}

	distanceCalculation := models.DistanceCalculation{
		OriginAddress:        originAddress,
		OriginAddressID:      originAddress.ID,
		DestinationAddress:   destinationAddress,
		DestinationAddressID: destinationAddress.ID,
		DistanceMiles:        1044,
	}

	mergeModels(&distanceCalculation, assertions.DistanceCalculation)

	mustCreate(db, &distanceCalculation, assertions.Stub)

	return distanceCalculation, nil
}

// MakeDefaultDistanceCalculation returns a DistanceCalculation with default values
func MakeDefaultDistanceCalculation(db *pop.Connection) (models.DistanceCalculation, error) {
	return MakeDistanceCalculation(db, Assertions{})
}
