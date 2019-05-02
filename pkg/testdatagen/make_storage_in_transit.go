package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

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
		EstimatedStartDate: NextValidMoveDate,
		Notes:              swag.String("Shipper phoned to let us know he is delayed until next week."),
		WarehouseID:        "000383",
		WarehouseName:      "Hercules Hauling",
		WarehouseAddressID: address.ID,
		WarehousePhone:     swag.String("(713) 868-3497"),
		WarehouseEmail:     swag.String("joe@herculeshauling.com"),
		Shipment:           shipment,
		WarehouseAddress:   address,
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

// ResetStorageInTransitSequenceNumber resets the storage in transit sequence number for a given year/dayOfYear.
func ResetStorageInTransitSequenceNumber(db *pop.Connection, year int, dayOfYear int) error {
	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	if dayOfYear <= 0 {
		return errors.Errorf("Day of year (%d) must be non-negative", dayOfYear)
	}

	sql := `DELETE FROM storage_in_transit_number_trackers WHERE year = $1 and day_of_year = $2`
	return db.RawQuery(sql, year, dayOfYear).Exec()
}

// SetStorageInTransitSequenceNumber sets the storage in transit sequence number for a given year/dayOfYear.
func SetStorageInTransitSequenceNumber(db *pop.Connection, year int, dayOfYear int, sequenceNumber int) error {
	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	if dayOfYear <= 0 {
		return errors.Errorf("Day of year (%d) must be non-negative", dayOfYear)
	}

	sql := `INSERT INTO storage_in_transit_number_trackers as trackers (year, day_of_year, sequence_number)
			VALUES ($1, $2, $3)
		ON CONFLICT (year, day_of_year)
		DO
			UPDATE
				SET sequence_number = $3
				WHERE trackers.year = $1 AND trackers.day_of_year = $2
	`

	return db.RawQuery(sql, year, dayOfYear, sequenceNumber).Exec()
}
