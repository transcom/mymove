package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildShipmentAddressUpdate creates a ShipmentAddressUpdate
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildShipmentAddressUpdate(db *pop.Connection, customs []Customization, traits []Trait) models.ShipmentAddressUpdate {
	customs = setupCustomizations(customs, traits)

	move := BuildMoveWithShipment(db, customs, traits)

	// Find ShipmentAddressUpdate assertion and convert to models ShipmentAddressUpdate
	var newShipmentAddress models.ShipmentAddressUpdate
	if result := findValidCustomization(customs, ShipmentAddressUpdate); result != nil {
		newShipmentAddress = result.Model.(models.ShipmentAddressUpdate)
		if result.LinkOnly {
			return newShipmentAddress
		}
	}

	// Create orig/new addresses
	originalAddress := BuildAddress(db, customs, traits)
	newAddress := BuildAddress(db, customs, traits)

	shipmentAddressUpdate := models.ShipmentAddressUpdate{
		ID:                uuid.Must(uuid.NewV4()),
		ContractorRemarks: "Test Contractor Remark",
		OfficeRemarks:     nil,
		Status:            models.ShipmentAddressUpdateStatusRequested,
		NewAddress:        newAddress,
		NewAddressID:      newAddress.ID,
		OriginalAddress:   originalAddress,
		OriginalAddressID: originalAddress.ID,
		ShipmentID:        move.MTOShipments[0].ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&shipmentAddressUpdate, newShipmentAddress)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &shipmentAddressUpdate)
	}

	return shipmentAddressUpdate
}

// ------------------------
//      TRAITS
// ------------------------

func GetTraitNonSITAddressUpdateRequested() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
	}
}

func GetTraitNonSITAddressUpdateApproved() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusApproved,
			},
		},
	}
}

func GetTraitNonSITAddressUpdateRejected() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRejected,
			},
		},
	}
}
