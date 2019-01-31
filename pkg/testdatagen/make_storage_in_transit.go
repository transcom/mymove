package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeStorageInTransit creates a single StorageInTransit with associations
func MakeStorageInTransit(db *pop.Connection, assertions Assertions) models.StorageInTransit {
	shipment := assertions.StorageInTransit.Shipment
	if isZeroUUID(shipment.ID) {
		shipment = MakeShipment(db, assertions)
	}

	address := assertions.StorageInTransit.WarehouseAddress
	if isZeroUUID(address.ID) {
		address = MakeAddress(db, assertions)
	}

	// Filled in dummy data.
	storageInTransit := models.StorageInTransit{
		ShipmentID:         shipment.ID,
		Status:             models.StorageInTransitStatusREQUESTED,
		Location:           models.StorageInTransitLocationDESTINATION,
		EstimatedStartDate: time.Now(),
		Notes:              swag.String("Shipper phoned to let us know he is delayed until next week."),
		WarehouseID:        "000383",
		WarehouseName:      "Hercules Hauling",
		WarehouseAddressID: address.ID,
		WarehousePhone:     swag.String("(713) 868-3497"),
		WarehouseEmail:     swag.String("joe@herculeshauling.com"),
	}

	// Overwrite values with those from assertions
	mergeModels(&storageInTransit, assertions.StorageInTransit)

	mustCreate(db, &storageInTransit)

	return storageInTransit
}

// MakeDefaultStorageInTransit makes a single StorageInTransit with default values
func MakeDefaultStorageInTransit(db *pop.Connection) models.StorageInTransit {
	return MakeStorageInTransit(db, Assertions{})
}
