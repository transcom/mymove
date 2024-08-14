package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
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

	tempOldAddressCustoms := customs
	if result := findValidCustomization(customs, Addresses.SITAddressUpdateOldAddress); result != nil {
		tempOldAddressCustoms = convertCustomizationInList(tempOldAddressCustoms, Addresses.SITAddressUpdateOldAddress, Address)
	}
	oldAddress := BuildAddress(db, tempOldAddressCustoms, traits)

	if db != nil {
		// Now we need a LinkOnly customization for the created address
		// can only do LinkOnly if we have an ID, which we won't have
		// for a stubbed address
		customs = replaceCustomization(customs, Customization{
			Model:    oldAddress,
			LinkOnly: true,
			Type:     &Addresses.SITDestinationFinalAddress,
		})
	}

	//Make sure new address is different from old if no customizations/traits were passed in
	tempTraits := append(traits, GetTraitAddress2)
	tempNewAddressCustoms := customs
	if result := findValidCustomization(customs, Addresses.SITAddressUpdateNewAddress); result != nil {
		tempNewAddressCustoms = convertCustomizationInList(tempNewAddressCustoms, Addresses.SITAddressUpdateNewAddress, Address)
	}
	newAddress := BuildAddress(db, tempNewAddressCustoms, tempTraits)

	serviceItem := BuildMTOServiceItem(db, customs, traits)

	// Create default SITAddressUpdate
	SITAddressUpdate := models.SITAddressUpdate{
		MTOServiceItem:    serviceItem,
		MTOServiceItemID:  serviceItem.ID,
		OldAddress:        oldAddress,
		OldAddressID:      oldAddress.ID,
		NewAddress:        newAddress,
		NewAddressID:      newAddress.ID,
		ContractorRemarks: models.StringPointer("contractor remarks"),
		Distance:          40,
		Status:            models.SITAddressUpdateStatusRequested,
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
				Status:   models.SITAddressUpdateStatusRequested,
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
				Status:   models.SITAddressUpdateStatusApproved,
			},
		},
	}
}

// GetTraitSITAddressUpdateWithMoveSetUp sets up a move with a REQUESTED SITAddressUpdate of over 50 miles
func GetTraitSITAddressUpdateWithMoveSetUp() []Customization {
	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	sitDaysAllowance := 200
	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	threeMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	twoMonthsAgo := threeMonthsAgo.Add(time.Hour * 24 * 30)
	originalPostalCode := "90210"
	reason := "peak season all trucks in use"

	return []Customization{
		{
			Model: models.Address{
				City:       "Beverly Hills",
				State:      "CA",
				PostalCode: originalPostalCode,
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
				Status:   models.SITAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Entitlement{
				DependentsAuthorized: models.BoolPointer(true),
				StorageInTransit:     &sitDaysAllowance,
			},
		},
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:  models.PoundPointer(unit.Pound(1400)),
				PrimeActualWeight:     models.PoundPointer(unit.Pound(2000)),
				ShipmentType:          models.MTOShipmentTypeHHG,
				Status:                models.MTOShipmentStatusApproved,
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedDeliveryDate,
				SITDaysAllowance:      &sitDaysAllowance,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &originalPostalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
	}
}
