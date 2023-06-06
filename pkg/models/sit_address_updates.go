package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// SITAddressUpdateStatus represents the possible statuses for a mto shipment
type SITAddressUpdateStatus string

const (
	// SITAddressUpdateStatusRequested is the requested status type for SIT Address Update Items
	SITAddressUpdateStatusRequested SITAddressUpdateStatus = "REQUESTED"
	// SITAddressUpdateStatusRejected is the rejected status type for SIT Address Update Items
	SITAddressUpdateStatusRejected SITAddressUpdateStatus = "REJECTED"
	// SITAddressUpdateStatusApproved is the approved status type for SIT Address Update Items
	SITAddressUpdateStatusApproved SITAddressUpdateStatus = "APPROVED"
)

var AllowedSITAddressStatuses = []string{
	string(SITAddressUpdateStatusRequested),
	string(SITAddressUpdateStatusRejected),
	string(SITAddressUpdateStatusApproved),
}

type SITAddressUpdate struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	ContractorRemarks *string                `json:"contractor_remarks" db:"contractor_remarks"`
	Distance          int                    `json:"distance" db:"distance"`
	OfficeRemarks     *string                `json:"office_remarks" db:"office_remarks"`
	Status            SITAddressUpdateStatus `json:"status" db:"status"`
	CreatedAt         time.Time              `db:"created_at"`
	UpdatedAt         time.Time              `db:"updated_at"`

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
		&StringIsNilOrNotBlank{Name: "ContractorRemarks", Field: s.ContractorRemarks},
		&StringIsNilOrNotBlank{Name: "OfficeRemarks", Field: s.OfficeRemarks},
	), nil
}

// TableName overrides the table name used by Pop.
func (s SITAddressUpdate) TableName() string {
	return "sit_address_updates"
}

// SITAddressUpdates is a slice containing of SITAddressUpdates
type SITAddressUpdates []SITAddressUpdate

func FetchSITAddressUpdate(db *pop.Connection, sitAddressUpdateID uuid.UUID) (SITAddressUpdate, error) {
	var sitAddressUpdate SITAddressUpdate
	err := db.Eager("NewAddress").Find(&sitAddressUpdate, sitAddressUpdateID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return SITAddressUpdate{}, ErrFetchNotFound
		}
		return SITAddressUpdate{}, err
	}

	return sitAddressUpdate, nil
}
