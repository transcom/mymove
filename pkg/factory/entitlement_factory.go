package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildEntitlement creates an Entitlement
// Does not create other models
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

	// TBD: Find the Orders customization
	// Get the grade from it.
	var grade = models.StringPointer("")
	if grade == nil || *grade == "" {
		grade = models.StringPointer("E_1")
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
	entitlement.SetWeightAllotment(*grade)
	entitlement.DBAuthorizedWeight = entitlement.AuthorizedWeight()

	// Overwrite default values with those from custom Entitlement
	testdatagen.MergeModels(&entitlement, cEntitlement)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &entitlement)
	}

	return entitlement
}

func GetTraitEntitlementsEnabled() []Customization {
	return []Customization{
		{
			Model: models.Entitlement{
				DependentsAuthorized:  models.BoolPointer(true),
				NonTemporaryStorage:   models.BoolPointer(true),
				PrivatelyOwnedVehicle: models.BoolPointer(true),
				TotalDependents:       models.IntPointer(1),
			},
		},
	}
}

func setBoolPtr(customBoolPtr *bool, defaultBool bool) *bool {
	result := models.BoolPointer(defaultBool)
	if customBoolPtr != nil {
		result = customBoolPtr
	}
	return result
}
