package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildStorageFacility creates a single StorageFacility.
// Also creates, if not provided
// - Address
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildStorageFacility(db *pop.Connection, customs []Customization, traits []Trait) models.StorageFacility {
	customs = setupCustomizations(customs, traits)

	// Find StorageFacility assertion and convert to models.StorageFacility
	var cStorageFacility models.StorageFacility
	if result := findValidCustomization(customs, StorageFacility); result != nil {
		cStorageFacility = result.Model.(models.StorageFacility)
		if result.LinkOnly {
			return cStorageFacility
		}
	}

	// Find/create the address model
	address := BuildAddress(db, customs, traits)

	// Create StorageFacility
	StorageFacility := models.StorageFacility{
		FacilityName: "Storage R Us",
		Address:      address,
		AddressID:    address.ID,
		LotNumber:    models.StringPointer("1234"),
		Phone:        models.StringPointer("555-555-5555"),
		Email:        models.StringPointer("storage@email.com"),
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&StorageFacility, cStorageFacility)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &StorageFacility)
	}
	return StorageFacility
}

// BuildDefaultStorageFacility creates one with a phoneline hooked up.
func BuildDefaultStorageFacility(db *pop.Connection) models.StorageFacility {
	return BuildStorageFacility(db, nil, nil)
}

// ------------------------
//        TRAITS
// ------------------------

// GetTraitOfficeUserEmail helps comply with the uniqueness constraint on emails
func GetTraitStorageFacilityKKFA() []Customization {
	return []Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}
}
