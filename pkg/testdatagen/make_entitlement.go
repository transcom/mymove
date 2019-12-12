package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeEntitlement creates a single Entitlement
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	truePtr := true
	dependents := 1
	proGearWeight := 100
	proGearWeightSpouse := 200
	storageInTransit := 2

	entitlement := models.Entitlement{
		DependentsAuthorized:  &truePtr,
		TotalDependents:       &dependents,
		NonTemporaryStorage:   &truePtr,
		PrivatelyOwnedVehicle: &truePtr,
		ProGearWeight:         &proGearWeight,
		ProGearWeightSpouse:   &proGearWeightSpouse,
		StorageInTransit:      &storageInTransit,
	}

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement)

	return entitlement
}
