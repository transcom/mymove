package factory

import (
	"time"

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

	shipment := BuildMTOShipment(db, customs, traits)

	// Find ShipmentAddressUpdate assertion and convert to models ShipmentAddressUpdate
	var newShipmentAddress models.ShipmentAddressUpdate
	if result := findValidCustomization(customs, ShipmentAddressUpdate); result != nil {
		newShipmentAddress = result.Model.(models.ShipmentAddressUpdate)
		if result.LinkOnly {
			return newShipmentAddress
		}
	}

	// Use shipment dest address as original address unless customizations are provided
	originalAddress := *shipment.DestinationAddress
	tempOrigAddressCustoms := customs
	validOrigCustoms := findValidCustomization(customs, Addresses.OriginalAddress)
	if validOrigCustoms != nil {
		tempOrigAddressCustoms = convertCustomizationInList(tempOrigAddressCustoms, Addresses.OriginalAddress, Address)
		// Create Original Address
		originalAddress = BuildAddress(db, tempOrigAddressCustoms, traits)
	}

	// Find New Address Customizations
	tempNewAddressCustoms := customs
	validNewCustoms := findValidCustomization(customs, Addresses.NewAddress)
	if validNewCustoms != nil {
		tempNewAddressCustoms = convertCustomizationInList(tempNewAddressCustoms, Addresses.NewAddress, Address)
	}
	// Create New Address
	newAddress := BuildAddress(db, tempNewAddressCustoms, traits)

	shipmentAddressUpdate := models.ShipmentAddressUpdate{
		ID:                uuid.Must(uuid.NewV4()),
		ContractorRemarks: "Customer reached out to me this week & let me know they want to move closer to a sick dependent who needs care.",
		OfficeRemarks:     nil,
		Status:            models.ShipmentAddressUpdateStatusRequested,
		NewAddress:        newAddress,
		NewAddressID:      newAddress.ID,
		OriginalAddress:   originalAddress,
		OriginalAddressID: originalAddress.ID,
		ShipmentID:        shipment.ID,
		Shipment:          shipment,
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

func GetTraitShipmentAddressUpdateRequested() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRequested,
			},
		},
		{
			Model: models.Move{
				Locator:            "CRQST1",
				Status:             models.MoveStatusAPPROVALSREQUESTED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				ApprovedAt:         models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApprovalsRequested,
			},
		},
	}
}

func GetTraitShipmentAddressUpdateApproved() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusApproved,
			},
		},
		{
			Model: models.Move{
				Locator:            "CRQST2",
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				ApprovedAt:         models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}
}

func GetTraitShipmentAddressUpdateRejected() []Customization {
	return []Customization{
		{
			Model: models.ShipmentAddressUpdate{
				Status: models.ShipmentAddressUpdateStatusRejected,
			},
		},
		{
			Model: models.Move{
				Locator:            "CRQST3",
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				ApprovedAt:         models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}
}
