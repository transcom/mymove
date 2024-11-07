package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildUBAllowance creates a UB allowance
// Does not create other models
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUBAllowance(db *pop.Connection, customs []Customization, traits []Trait) models.UBAllowances {
	customs = setupCustomizations(customs, traits)

	// Find UBAllowances Customization and extract the custom UBAllowances
	var cUBAllowance models.UBAllowances
	if result := findValidCustomization(customs, UBAllowance); result != nil {
		cUBAllowance = result.Model.(models.UBAllowances)
		if result.LinkOnly {
			return cUBAllowance
		}
	}
	ubAllowanceValue := 2000

	ubAllowance := models.UBAllowances{
		BranchOfService: models.AffiliationAIRFORCE.String(),
		OrderPayGrade:   string(models.ServiceMemberGradeE1),
		OrdersType:      string(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION),
		HasDependents:   true,
		AccompaniedTour: true,
		UBAllowance:     ubAllowanceValue,
	}

	// Overwrite default values with those from custom UB allowance
	testdatagen.MergeModels(&ubAllowance, cUBAllowance)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &ubAllowance)
	}

	return ubAllowance
}
