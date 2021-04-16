package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeEntitlement creates a single Entitlement
func MakeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	truePtr := true
	dependents := 1
	storageInTransit := 2
	rmeWeight := 1000
	ocie := true
	grade := assertions.Order.Grade

	if grade == nil || *grade == "" {
		grade = stringPointer("E_1")
	}

	entitlement := models.Entitlement{
		DependentsAuthorized:                         setDependentsAuthorized(assertions.Entitlement.DependentsAuthorized),
		TotalDependents:                              &dependents,
		NonTemporaryStorage:                          &truePtr,
		PrivatelyOwnedVehicle:                        &truePtr,
		StorageInTransit:                             &storageInTransit,
		RequiredMedicalEquipmentWeight:               rmeWeight,
		OrganizationalClothingAndIndividualEquipment: ocie,
	}
	entitlement.SetWeightAllotment(*grade)
	dBAuthorizedWeight := entitlement.AuthorizedWeight()
	entitlement.DBAuthorizedWeight = dBAuthorizedWeight

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement, assertions.Stub)

	return entitlement
}

func setDependentsAuthorized(assertionDependentsAuthorized *bool) *bool {
	dependentsAuthorized := swag.Bool(true)
	if assertionDependentsAuthorized != nil {
		dependentsAuthorized = assertionDependentsAuthorized
	}
	return dependentsAuthorized
}
