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

	// Find OfficePhoneLine assertion and convert to models officephoneline
	var newShipmentAddress models.ShipmentAddressUpdate
	if result := findValidCustomization(customs, ShipmentAddressUpdate); result != nil {
		newShipmentAddress = result.Model.(models.ShipmentAddressUpdate)
		if result.LinkOnly {
			return newShipmentAddress
		}
	}
	hardcodeUUID, _ := uuid.FromString("01b9671e-b268-4906-967b-ba661a1d3933")
	// create newShipmentAddress
	shipmentAddressUpdate := models.ShipmentAddressUpdate{
		ID: uuid.Must(uuid.NewV4()),
		// ContractorRemarks:
		// OfficeRemarks:
		// Status:
		// CreatedAt:
		// UpdatedAt:
		// we need to know the shipment id before this point
		ShipmentID: hardcodeUUID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&shipmentAddressUpdate, newShipmentAddress)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &shipmentAddressUpdate)
	}

	return shipmentAddressUpdate
}

// ID                uuid.UUID                   `json:"id" db:"id"`
// ContractorRemarks string                      `json:"contractor_remarks" db:"contractor_remarks"`
// OfficeRemarks     *string                     `json:"office_remarks" db:"office_remarks"`
// Status            ShipmentAddressUpdateStatus `json:"status" db:"status"`
// CreatedAt         time.Time                   `db:"created_at"`
// UpdatedAt         time.Time                   `db:"updated_at"`

// what do we do with default associations?
// Associations
// Shipment          MTOShipment `belongs_to:"mto_shipments" fk_id:"shipment_id"`
// ShipmentID        uuid.UUID   `db:"shipment_id"`
// OriginalAddress   Address     `belongs_to:"addresses" fk_id:"original_address_id"`
// OriginalAddressID uuid.UUID   `db:"original_address_id"`
// NewAddress        Address     `belongs_to:"addresses" fk_id:"new_address_id"`
// NewAddressID      uuid.UUID   `db:"new_address_id"`
