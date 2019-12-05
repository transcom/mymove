package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

//  MakeEntitlement creates a single Entitlement
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	entitlement := models.Entitlement{
		DependentsAuthorized:  true,
		TotalDependents:       1,
		NonTemporaryStorage:   true,
		PrivatelyOwnedVehicle: true,
		ProGearWeight:         100,
		ProGearWeightSpouse:   200,
		StorageInTransit:      2,
	}
	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement)

	return entitlement
}
