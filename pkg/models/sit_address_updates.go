package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// SITAddressStatus represents the possible statuses for a mto shipment
type SITAddressStatus string

const (
	// SITAddressStatusRequested is the requested status type for SIT Address Update Items
	SITAddressStatusRequested SITAddressStatus = "REQUESTED"
	// SITAddressStatusRejected is the rejected status type for SIT Address Update Items
	SITAddressStatusRejected SITAddressStatus = "REJECTED"
	// SITAddressStatusApproved is the approved status type for SIT Address Update Items
	SITAddressStatusApproved SITAddressStatus = "APPROVED"
)

var AllowedSITAddressStatuses = []string{
	string(SITAddressStatusRequested),
	string(SITAddressStatusRejected),
	string(SITAddressStatusApproved),
}

type SITAddressUpdate struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	ContractorRemarks string           `json:"contractor_remarks" db:"contractor_remarks"`
	Distance          int              `json:"distance" db:"distance"`
	OfficeRemarks     *string          `json:"office_remarks" db:"office_remarks"`
	Reason            string           `json:"reason" db:"reason"`
	Status            SITAddressStatus `json:"status" db:"status"`

	// Associations
	MTOServiceItem   MTOServiceItem `belongs_to:"mto_service_items" fk_id:"mto_service_item_id"`
	MTOServiceItemID uuid.UUID      `db:"mto_service_item_id"`
	OldAddress       Address        `belongs_to:"addresses" fk_id:"old_address_id"`
	OldAddressID     uuid.UUID      `db:"old_address_id"`
	NewAddress       Address        `belongs_to:"addresses" fk_id:"new_address_id"`
	NewAddressID     uuid.UUID      `db:"new_address_id"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (s *SITAddressUpdate) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "MTOServiceItemID", Field: s.MTOServiceItemID},
		&validators.UUIDIsPresent{Name: "OldAddressID", Field: s.OldAddressID},
		&validators.UUIDIsPresent{Name: "NewAddressID", Field: s.NewAddressID},
		&validators.StringInclusion{Name: "Status", Field: string(s.Status), List: AllowedSITAddressStatuses},
		&validators.IntIsPresent{Name: "Distance", Field: s.Distance},
		&validators.StringIsPresent{Name: "Reason", Field: s.Reason},
		&validators.StringIsPresent{Name: "ContractorRemarks", Field: s.ContractorRemarks},
		&StringIsNilOrNotBlank{Name: "OfficeRemarks", Field: s.OfficeRemarks},
	), nil
}

// TableName overrides the table name used by Pop.
func (s SITAddressUpdate) TableName() string {
	return "sit_address_updates"
}

// SITAddressUpdates is a slice containing of SITAddressUpdates
type SITAddressUpdates []SITAddressUpdate
