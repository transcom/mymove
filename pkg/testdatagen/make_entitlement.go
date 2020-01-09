package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeEntitlement creates a single Entitlement
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	truePtr := true
	dependents := 1
	storageInTransit := 2
	grade := models.ServiceMemberRank(assertions.MoveOrder.Grade)

	if grade == "" {
		grade = models.ServiceMemberRankE1
	}

	allotment := models.GetWeightAllotment(grade)
	entitlement := models.Entitlement{
		DependentsAuthorized:  &truePtr,
		TotalDependents:       &dependents,
		NonTemporaryStorage:   &truePtr,
		PrivatelyOwnedVehicle: &truePtr,
		StorageInTransit:      &storageInTransit,
		WeightAllotment:       &allotment,
	}

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement)

	return entitlement
}
