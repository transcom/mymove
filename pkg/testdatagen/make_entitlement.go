package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// makeEntitlement creates a single Entitlement
func makeEntitlement(db *pop.Connection, assertions Assertions) models.Entitlement {
	truePtr := true
	dependents := 1
	storageInTransit := 90
	rmeWeight := 1000
	ocie := true
	grade := assertions.Order.Grade
	ordersType := assertions.Order.OrdersType
	proGearWeight := 2000
	proGearWeightSpouse := 500

	if grade == nil || *grade == "" {
		grade = models.ServiceMemberGradeE1.Pointer()
	}

	entitlement := models.Entitlement{
		DependentsAuthorized:                         setDependentsAuthorized(assertions.Entitlement.DependentsAuthorized),
		TotalDependents:                              &dependents,
		NonTemporaryStorage:                          &truePtr,
		PrivatelyOwnedVehicle:                        &truePtr,
		StorageInTransit:                             &storageInTransit,
		ProGearWeight:                                proGearWeight,
		ProGearWeightSpouse:                          proGearWeightSpouse,
		RequiredMedicalEquipmentWeight:               rmeWeight,
		OrganizationalClothingAndIndividualEquipment: ocie,
	}
	entitlement.SetWeightAllotment(string(*grade), ordersType)
	dBAuthorizedWeight := entitlement.AuthorizedWeight()
	entitlement.DBAuthorizedWeight = dBAuthorizedWeight

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement, assertions.Stub)

	return entitlement
}

func setDependentsAuthorized(assertionDependentsAuthorized *bool) *bool {
	dependentsAuthorized := models.BoolPointer(true)
	if assertionDependentsAuthorized != nil {
		dependentsAuthorized = assertionDependentsAuthorized
	}
	return dependentsAuthorized
}
