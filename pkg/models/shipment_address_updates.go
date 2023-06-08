package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ShipmentAddressUpdateStatus represents the possible statuses for a mto shipment
type ShipmentAddressUpdateStatus string

const (
	// ShipmentAddressUpdateStatusRequested is the requested status type for Shipment Address Update Items
	ShipmentAddressUpdateStatusRequested ShipmentAddressUpdateStatus = "REQUESTED"
	// ShipmentAddressUpdateStatusRejected is the rejected status type for Shipment Address Update Items
	ShipmentAddressUpdateStatusRejected ShipmentAddressUpdateStatus = "REJECTED"
	// ShipmentAddressUpdateStatusApproved is the approved status type for Shipment Address Update Items
	ShipmentAddressUpdateStatusApproved ShipmentAddressUpdateStatus = "APPROVED"
)

var AllowedShipmentAddressStatuses = []string{
	string(ShipmentAddressUpdateStatusRequested),
	string(ShipmentAddressUpdateStatusRejected),
	string(ShipmentAddressUpdateStatusApproved),
}

type ShipmentAddressUpdate struct {
	ID                             uuid.UUID `json:"id" db:"id"`
	ContractorRemarks              string    `json:"contractor_remarks" db:"contractor_remarks"`
	ServiceAreaChanged             bool      `json:"service_area_changed" db:"service_area_changed"`
	MileageBracketChanged          bool      `json:"mileage_bracket_changed" db:"mileage_bracket_changed"`
	ChangedFromShortHaulToLongHaul bool      `json:"changed_from_short_haul_to_long_haul" db:"changed_from_short_haul_to_long_haul"`
	ChangedFromLongHaulToShortHaul bool      `json:"changed_from_long_haul_to_short_haul" db:"changed_from_long_haul_to_short_haul"`

	OfficeRemarks *string                     `json:"office_remarks" db:"office_remarks"`
	Status        ShipmentAddressUpdateStatus `json:"status" db:"status"`
	CreatedAt     time.Time                   `db:"created_at"`
	UpdatedAt     time.Time                   `db:"updated_at"`

	// Associations
	Shipment          MTOShipment `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	ShipmentID        uuid.UUID   `db:"shipment_id"`
	OriginalAddress   Address     `belongs_to:"addresses" fk_id:"original_address_id"`
	OriginalAddressID uuid.UUID   `db:"original_address_id"`
	NewAddress        Address     `belongs_to:"addresses" fk_id:"new_address_id"`
	NewAddressID      uuid.UUID   `db:"new_address_id"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (s *ShipmentAddressUpdate) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: s.ShipmentID},
		&validators.UUIDIsPresent{Name: "OriginalAddressID", Field: s.OriginalAddressID},
		&validators.UUIDIsPresent{Name: "NewAddressID", Field: s.NewAddressID},
		&validators.StringInclusion{Name: "Status", Field: string(s.Status), List: AllowedShipmentAddressStatuses},
		&validators.StringIsPresent{Name: "ContractorRemarks", Field: s.ContractorRemarks},
		&StringIsNilOrNotBlank{Name: "OfficeRemarks", Field: s.OfficeRemarks},
	), nil
}

// TableName overrides the table name used by Pop.
func (s ShipmentAddressUpdate) TableName() string {
	return "shipment_address_updates"
}

// ShipmentAddressUpdates is a slice containing of ShipmentAddressUpdates
type ShipmentAddressUpdates []ShipmentAddressUpdate
