package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildSITAddressUpdate creates an SITAddressUpdate
// It builds
//   - MTOServiceItem and associated set relationships
//   - OldAddress
//   - NewAddress
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildSITAddressUpdate(db *pop.Connection, customs []Customization, traits []Trait) models.SITAddressUpdate {
	customs = setupCustomizations(customs, traits)

	// Find SITAddressUpdate Customization and extract the custom SITAddressUpdate
	var cSITAddressUpdate models.SITAddressUpdate
	if result := findValidCustomization(customs, SITAddressUpdate); result != nil {
		cSITAddressUpdate = result.Model.(models.SITAddressUpdate)
		if result.LinkOnly {
			return cSITAddressUpdate
		}
	}

	serviceItem := BuildMTOServiceItem(db, customs, traits)

	tempOldAddressCustoms := customs
	if result := findValidCustomization(customs, Addresses.SITAddressUpdateOldAddress); result != nil {
		tempOldAddressCustoms = convertCustomizationInList(tempOldAddressCustoms, Addresses.SITAddressUpdateOldAddress, Address)
	}
	oldAddress := BuildAddress(db, tempOldAddressCustoms, traits)

	//Make sure new address is different from old if no customizations/traits were passed in
	traits = append(traits, GetTraitAddress2)
	tempNewAddressCustoms := customs
	if result := findValidCustomization(customs, Addresses.SITAddressUpdateNewAddress); result != nil {
		tempNewAddressCustoms = convertCustomizationInList(tempNewAddressCustoms, Addresses.SITAddressUpdateNewAddress, Address)
	}
	newAddress := BuildAddress(db, tempNewAddressCustoms, traits)

	// Create default SITAddressUpdate
	SITAddressUpdate := models.SITAddressUpdate{
		MTOServiceItem:    serviceItem,
		MTOServiceItemID:  serviceItem.ID,
		OldAddress:        oldAddress,
		OldAddressID:      oldAddress.ID,
		NewAddress:        newAddress,
		NewAddressID:      newAddress.ID,
		ContractorRemarks: "contractor remarks",
		Distance:          40,
		Reason:            "new reason",
		Status:            models.SITAddressStatusRequested,
	}

	// Overwrite default values with those from custom SITAddressUpdate
	testdatagen.MergeModels(&SITAddressUpdate, cSITAddressUpdate)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &SITAddressUpdate)
	}

	return SITAddressUpdate
}

// ------------------------
//      TRAITS
// ------------------------

func GetTraitSITAddressUpdateOver50Miles() []Customization {
	return []Customization{
		{
			Model: models.Address{
				City:       "Beverly Hills",
				State:      "CA",
				PostalCode: "90210",
			},
			Type: &Addresses.SITAddressUpdateOldAddress,
		},
		{
			Model: models.Address{
				City:       "San Diego",
				State:      "CA",
				PostalCode: "92114",
			},
			Type: &Addresses.SITAddressUpdateNewAddress,
		},
		{
			Model: models.SITAddressUpdate{
				Distance: 140,
				Status:   models.SITAddressStatusRequested,
			},
		},
	}
}

func GetTraitSITAddressUpdateUnder50Miles() []Customization {
	return []Customization{
		{
			Model: models.Address{
				City:       "Beverly Hills",
				State:      "CA",
				PostalCode: "90210",
			},
			Type: &Addresses.SITAddressUpdateOldAddress,
		},
		{
			Model: models.Address{
				City:       "Long Beach",
				State:      "CA",
				PostalCode: "90802",
			},
			Type: &Addresses.SITAddressUpdateNewAddress,
		},
		{
			Model: models.SITAddressUpdate{
				Distance: 16,
				Status:   models.SITAddressStatusApproved,
			},
		},
	}
}

func GetTraitSITAddressUpdateRejected() []Customization {
	return []Customization{
		{
			Model: models.Address{
				City:       "Beverly Hills",
				State:      "CA",
				PostalCode: "90210",
			},
			Type: &Addresses.SITAddressUpdateOldAddress,
		},
		{
			Model: models.Address{
				City:       "San Diego",
				State:      "CA",
				PostalCode: "92114",
			},
			Type: &Addresses.SITAddressUpdateNewAddress,
		},
		{
			Model: models.SITAddressUpdate{
				Distance: 140,
				Status:   models.SITAddressStatusRejected,
			},
		},
	}
}
