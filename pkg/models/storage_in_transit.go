package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// StorageInTransitStatus represents the status of a SIT request
type StorageInTransitStatus string

// StorageInTransitLocation represents the location of the SIT request
type StorageInTransitLocation string

const (
	// StorageInTransitStatusREQUESTED represents an initial SIT request
	StorageInTransitStatusREQUESTED StorageInTransitStatus = "REQUESTED"
	// StorageInTransitStatusAPPROVED represents an approved SIT request
	StorageInTransitStatusAPPROVED StorageInTransitStatus = "APPROVED"
	// StorageInTransitStatusDENIED represents a denied SIT request
	StorageInTransitStatusDENIED StorageInTransitStatus = "DENIED"
	// StorageInTransitStatusINSIT represents a shipment that is current in SIT
	StorageInTransitStatusINSIT StorageInTransitStatus = "IN_SIT"
	// StorageInTransitStatusRELEASED represents a shipment that has been released from SIT
	StorageInTransitStatusRELEASED StorageInTransitStatus = "RELEASED"
	// StorageInTransitStatusDELIVERED represents a shipment that has been delivered
	StorageInTransitStatusDELIVERED StorageInTransitStatus = "DELIVERED"

	// StorageInTransitLocationORIGIN represents SIT at the origin
	StorageInTransitLocationORIGIN StorageInTransitLocation = "ORIGIN"
	// StorageInTransitLocationDESTINATION represents SIT at the destination
	StorageInTransitLocationDESTINATION StorageInTransitLocation = "DESTINATION"
)

var storageInTransitStatuses = []string{
	string(StorageInTransitStatusREQUESTED),
	string(StorageInTransitStatusAPPROVED),
	string(StorageInTransitStatusDENIED),
	string(StorageInTransitStatusINSIT),
	string(StorageInTransitStatusRELEASED),
	string(StorageInTransitStatusDELIVERED),
}

var storageInTransitLocations = []string{
	string(StorageInTransitLocationORIGIN),
	string(StorageInTransitLocationDESTINATION),
}

// StorageInTransit represents a single SIT request for a shipment
type StorageInTransit struct {
	ID                  uuid.UUID                `json:"id" db:"id"`
	CreatedAt           time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at" db:"updated_at"`
	ShipmentID          uuid.UUID                `json:"shipment_id" db:"shipment_id"`
	SITNumber           *string                  `json:"sit_number" db:"sit_number"`
	Status              StorageInTransitStatus   `json:"status" db:"status"`
	Location            StorageInTransitLocation `json:"location" db:"location"`
	EstimatedStartDate  time.Time                `json:"estimated_start_date" db:"estimated_start_date"`
	AuthorizedStartDate *time.Time               `json:"authorized_start_date" db:"authorized_start_date"`
	ActualStartDate     *time.Time               `json:"actual_start_date" db:"actual_start_date"`
	OutDate             *time.Time               `json:"out_date" db:"out_date"`
	Notes               *string                  `json:"notes" db:"notes"`
	AuthorizationNotes  *string                  `json:"authorization_notes" db:"authorization_notes"`
	WarehouseID         string                   `json:"warehouse_id" db:"warehouse_id"`
	WarehouseName       string                   `json:"warehouse_name" db:"warehouse_name"`
	WarehouseAddressID  uuid.UUID                `json:"warehouse_address_id" db:"warehouse_address_id"`
	WarehousePhone      *string                  `json:"warehouse_phone" db:"warehouse_phone"`
	WarehouseEmail      *string                  `json:"warehouse_email" db:"warehouse_email"`

	// Associations
	Shipment         Shipment `belongs_to:"shipment"`
	WarehouseAddress Address  `belongs_to:"address"`
}

// StorageInTransits is not required by pop and may be deleted
type StorageInTransits []StorageInTransit

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *StorageInTransit) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.ShipmentID, Name: "ShipmentID"},
		&StringIsNilOrNotBlank{Field: s.SITNumber, Name: "SITNumber"},
		&validators.StringInclusion{Field: string(s.Status), Name: "Status", List: storageInTransitStatuses},
		&validators.StringInclusion{Field: string(s.Location), Name: "Location", List: storageInTransitLocations},
		&validators.TimeIsPresent{Field: s.EstimatedStartDate, Name: "EstimatedStartDate"},
		&OptionalTimeIsPresent{Field: s.AuthorizedStartDate, Name: "AuthorizedStartDate"},
		&OptionalTimeIsPresent{Field: s.ActualStartDate, Name: "ActualStartDate"},
		&OptionalTimeIsPresent{Field: s.OutDate, Name: "OutDate"},
		&StringIsNilOrNotBlank{Field: s.Notes, Name: "Notes"},
		&StringIsNilOrNotBlank{Field: s.AuthorizationNotes, Name: "AuthorizationNotes"},
		&validators.StringIsPresent{Field: s.WarehouseID, Name: "WarehouseID"},
		&validators.StringIsPresent{Field: s.WarehouseName, Name: "WarehouseName"},
		&validators.UUIDIsPresent{Field: s.WarehouseAddressID, Name: "WarehouseAddressID"},
		&StringIsNilOrNotBlank{Field: s.WarehousePhone, Name: "WarehousePhone"},
		&StringIsNilOrNotBlank{Field: s.WarehouseEmail, Name: "WarehouseEmail"},
	), nil
}

// FetchStorageInTransitsOnShipment retrieves Storage In Transit objects using the shipment ID
func FetchStorageInTransitsOnShipment(tx *pop.Connection, shipmentID uuid.UUID) (StorageInTransits, error) {
	storageInTransits := StorageInTransits{}

	err := tx.Eager().Where("shipment_id = $1", shipmentID).
		LeftJoin("addresses", "storage_in_transits.warehouse_address_id=addresses.id").
		All(&storageInTransits)

	if err != nil {
		return nil, err
	}

	return storageInTransits, nil

}

// FetchStorageInTransitByID retrieves a single Storage In Transit object based on its own ID
func FetchStorageInTransitByID(tx *pop.Connection, storageInTransitID uuid.UUID) (*StorageInTransit, error) {
	var storageInTransit StorageInTransit
	err := tx.Eager().Where("storage_in_transits.id = $1", storageInTransitID).
		LeftJoin("addresses", "storage_in_transits.warehouse_address_id=addresses.id").First(&storageInTransit)

	if err != nil {
		// If we fail to get rows let's pass up a ErrFetchNotFound so that handlers wind up passing a 404
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}

	return &storageInTransit, nil

}

// DeleteStorageInTransit deletes a Storage In Transit object based on the provided ID
func DeleteStorageInTransit(tx *pop.Connection, storageInTransitID uuid.UUID) (err error) {
	var storageInTransit StorageInTransit

	// Identify the record we're going to delete by its ID
	// If we can't find it we return an ErrFetchNotFound.
	err = tx.Find(&storageInTransit, storageInTransitID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ErrFetchNotFound
		}
		return err
	}

	// Execute the deletion
	err = tx.Destroy(&storageInTransit)
	if err != nil {
		return err
	}

	return nil

}
