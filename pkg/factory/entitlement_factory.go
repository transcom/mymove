package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildEntitlement creates an Entitlement
// Does not create other models
// If Orders customization is provided - it will use the grade to calculate weight allotment
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildEntitlement(db *pop.Connection, customs []Customization, traits []Trait) models.Entitlement {
	customs = setupCustomizations(customs, traits)

	// Find Entitlement Customization and extract the custom Entitlement
	var cEntitlement models.Entitlement
	if result := findValidCustomization(customs, Entitlement); result != nil {
		cEntitlement = result.Model.(models.Entitlement)
		if result.LinkOnly {
			return cEntitlement
		}
	}
	// At this point, Entitlement may exist or be blank. Depending on if customization was provided.

	// Find an Orders customization if available - to extract the grade
	var grade *internalmessages.OrderPayGrade
	defaultGrade := models.ServiceMemberGradeE1
	var order models.Order
	result := findValidCustomization(customs, Order)
	if result != nil {
		order = result.Model.(models.Order)
		grade = order.Grade // extract grade
	}
	if grade == nil || *grade == "" {
		grade = &defaultGrade
	}

	dependents := 0
	storageInTransit := 90
	rmeWeight := 1000
	ocie := true
	proGearWeight := 2000
	proGearWeightSpouse := 500

	// Create default Entitlement
	entitlement := models.Entitlement{
		DependentsAuthorized:                         setBoolPtr(cEntitlement.DependentsAuthorized, true),
		TotalDependents:                              &dependents,
		NonTemporaryStorage:                          setBoolPtr(cEntitlement.NonTemporaryStorage, true),
		PrivatelyOwnedVehicle:                        setBoolPtr(cEntitlement.PrivatelyOwnedVehicle, true),
		StorageInTransit:                             &storageInTransit,
		ProGearWeight:                                proGearWeight,
		ProGearWeightSpouse:                          proGearWeightSpouse,
		RequiredMedicalEquipmentWeight:               rmeWeight,
		OrganizationalClothingAndIndividualEquipment: ocie,
	}
	// Set default calculated values
	entitlement.SetWeightAllotment(string(*grade))
	entitlement.DBAuthorizedWeight = entitlement.AuthorizedWeight()

	// Overwrite default values with those from custom Entitlement
	testdatagen.MergeModels(&entitlement, cEntitlement)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &entitlement)
	}

	return entitlement
}

func BuildPayGrade(db *pop.Connection, customs []Customization, traits []Trait) models.PayGrade {
	customs = setupCustomizations(customs, traits)

	// Find Pay Grade Customization and extract the custom Pay Grade
	var cPayGrade models.PayGrade
	if result := findValidCustomization(customs, PayGrade); result != nil {
		cPayGrade = result.Model.(models.PayGrade)
		if result.LinkOnly {
			return cPayGrade
		}
	}

	// Check if the Grade already exists
	var existingPayGrade models.PayGrade
	if db != nil {
		err := db.Where("grade = ?", cPayGrade.Grade).First(&existingPayGrade)
		if err == nil {
			return existingPayGrade
		}
	}

	// Create default Pay Grade
	payGrade := models.PayGrade{
		Grade:            "E-5",
		GradeDescription: models.StringPointer("Enlisted Grade E-5"),
	}

	// Overwrite default values with those from custom Pay Grade
	testdatagen.MergeModels(&payGrade, cPayGrade)

	if db != nil {
		mustCreate(db, &payGrade)
	}

	return payGrade
}

func BuildHHGAllowance(db *pop.Connection, customs []Customization, traits []Trait) models.HHGAllowance {
	customs = setupCustomizations(customs, traits)

	// Find HHG Allowance Customization and extract the custom HHG Allowance
	var cHHGAllowance models.HHGAllowance
	if result := findValidCustomization(customs, HHGAllowance); result != nil {
		cHHGAllowance = result.Model.(models.HHGAllowance)
		if result.LinkOnly {
			return cHHGAllowance
		}
	}

	// Check if Allowance with this Grade already exists
	var existingHHGAllowance models.HHGAllowance
	if db != nil {
		err := db.Where("pay_grade_id = ?", cHHGAllowance.PayGradeID).First(&existingHHGAllowance)
		if err == nil {
			return existingHHGAllowance
		}
	}

	// Create a default HHG Allowance with default pay grade
	payGrade := BuildPayGrade(db, customs, traits)
	defaultWeightData := getDefaultWeightData(payGrade.Grade)

	hhgAllowance := models.HHGAllowance{
		PayGradeID:                    payGrade.ID,
		PayGrade:                      payGrade,
		TotalWeightSelf:               defaultWeightData.TotalWeightSelf,
		TotalWeightSelfPlusDependents: defaultWeightData.TotalWeightSelfPlusDependents,
		ProGearWeight:                 defaultWeightData.ProGearWeight,
		ProGearWeightSpouse:           defaultWeightData.ProGearWeightSpouse,
	}

	// Overwrite default values with those from custom HHG Allowance
	testdatagen.MergeModels(&hhgAllowance, cHHGAllowance)

	if db != nil {
		mustCreate(db, &hhgAllowance)
	}

	return hhgAllowance
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
