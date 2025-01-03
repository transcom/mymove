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

	weightData := getDefaultWeightData(string(*grade))
	allotment := models.WeightAllotment{
		TotalWeightSelf:               weightData.TotalWeightSelf,
		TotalWeightSelfPlusDependents: weightData.TotalWeightSelfPlusDependents,
		ProGearWeight:                 weightData.ProGearWeight,
		ProGearWeightSpouse:           weightData.ProGearWeightSpouse,
	}
	entitlement.WeightAllotted = &allotment
	dBAuthorizedWeight := entitlement.AuthorizedWeight()
	entitlement.DBAuthorizedWeight = dBAuthorizedWeight

	// Overwrite values with those from assertions
	mergeModels(&entitlement, assertions.Entitlement)

	mustCreate(db, &entitlement, assertions.Stub)

	return entitlement
}

// Helper function to retrieve default weight data by grade
func getDefaultWeightData(grade string) struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
} {
	if data, ok := knownAllowances[grade]; ok {
		return data
	}
	return knownAllowances["EMPTY"] // Default to EMPTY if grade not found. This is just dummy default data
}

// Default allowances CAO December 2024
// Note that the testdatagen package has its own default allowance
var knownAllowances = map[string]struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
}{
	"EMPTY":                    {0, 0, 0, 0},
	"ACADEMY_CADET":            {350, 350, 0, 0},
	"MIDSHIPMAN":               {350, 350, 0, 0},
	"AVIATION_CADET":           {7000, 8000, 2000, 500},
	"E-1":                      {5000, 8000, 2000, 500},
	"E-2":                      {5000, 8000, 2000, 500},
	"E-3":                      {5000, 8000, 2000, 500},
	"E-4":                      {7000, 8000, 2000, 500},
	"E-5":                      {7000, 9000, 2000, 500},
	"E-6":                      {8000, 11000, 2000, 500},
	"E-7":                      {11000, 13000, 2000, 500},
	"E-8":                      {12000, 14000, 2000, 500},
	"E-9":                      {13000, 15000, 2000, 500},
	"E-9SPECIALSENIORENLISTED": {14000, 17000, 2000, 500},
	"O-1ACADEMYGRADUATE":       {10000, 12000, 2000, 500},
	"O-2":                      {12500, 13500, 2000, 500},
	"O-3":                      {13000, 14500, 2000, 500},
	"O-4":                      {14000, 17000, 2000, 500},
	"O-5":                      {16000, 17500, 2000, 500},
	"O-6":                      {18000, 18000, 2000, 500},
	"O-7":                      {18000, 18000, 2000, 500},
	"O-8":                      {18000, 18000, 2000, 500},
	"O-9":                      {18000, 18000, 2000, 500},
	"O-10":                     {18000, 18000, 2000, 500},
	"W-1":                      {10000, 12000, 2000, 500},
	"W-2":                      {12500, 13500, 2000, 500},
	"W-3":                      {13000, 14500, 2000, 500},
	"W-4":                      {14000, 17000, 2000, 500},
	"W-5":                      {16000, 17500, 2000, 500},
	"CIVILIAN_EMPLOYEE":        {18000, 18000, 2000, 500},
}

func setDependentsAuthorized(assertionDependentsAuthorized *bool) *bool {
	dependentsAuthorized := models.BoolPointer(true)
	if assertionDependentsAuthorized != nil {
		dependentsAuthorized = assertionDependentsAuthorized
	}
	return dependentsAuthorized
}
