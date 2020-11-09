package testdatagen

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeEntitlement creates a single Entitlement
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	truePtr := true
	dependents := 1
	storageInTransit := 2
	grade := assertions.Order.Grade

	if grade == nil || *grade == "" {
		grade = stringPointer("E_1")
	}

	entitlement := models.Entitlement{
		DependentsAuthorized:  &truePtr,
		TotalDependents:       &dependents,
		NonTemporaryStorage:   &truePtr,
		PrivatelyOwnedVehicle: &truePtr,
		StorageInTransit:      &storageInTransit,
	}
	entitlement.SetWeightAllotment(*grade)

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement, assertions.Stub)

	return entitlement
}
