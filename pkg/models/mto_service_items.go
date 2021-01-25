package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// MTOServiceItemStatus represents the possible statuses for a mto shipment
type MTOServiceItemStatus string

const (
	// MTOServiceItemStatusSubmitted is the submitted status type for MTO Service Items
	MTOServiceItemStatusSubmitted MTOServiceItemStatus = "SUBMITTED"
	// MTOServiceItemStatusApproved is the approved status type for MTO Service Items
	MTOServiceItemStatusApproved MTOServiceItemStatus = "APPROVED"
	// MTOServiceItemStatusRejected is the rejected status type for MTO Service Items
	MTOServiceItemStatusRejected MTOServiceItemStatus = "REJECTED"
)

// MTOServiceItem is an object representing service items for a move task order.
type MTOServiceItem struct {
	ID                            uuid.UUID                      `db:"id"`
	MoveTaskOrder                 Move                           `belongs_to:"moves"`
	MoveTaskOrderID               uuid.UUID                      `db:"move_id"`
	MTOShipment                   MTOShipment                    `belongs_to:"mto_shipments"`
	MTOShipmentID                 *uuid.UUID                     `db:"mto_shipment_id"`
	ReService                     ReService                      `belongs_to:"re_services"`
	ReServiceID                   uuid.UUID                      `db:"re_service_id"`
	Reason                        *string                        `db:"reason"`
	RejectionReason               *string                        `db:"rejection_reason"`
	Status                        MTOServiceItemStatus           `db:"status"`
	PickupPostalCode              *string                        `db:"pickup_postal_code"`
	SITPostalCode                 *string                        `db:"sit_postal_code"`
	SITEntryDate                  *time.Time                     `db:"sit_entry_date"`
	SITDepartureDate              *time.Time                     `db:"sit_departure_date"`
	SITOriginHHGOriginalAddress   *Address                       `belongs_to:"addresses"`
	SITOriginHHGOriginalAddressID *uuid.UUID                     `db:"sit_origin_hhg_original_address_id"`
	SITOriginHHGActualAddress     *Address                       `belongs_to:"addresses"`
	SITOriginHHGActualAddressID   *uuid.UUID                     `db:"sit_origin_hhg_actual_address_id"`
	SITDestinationFinalAddress    *Address                       `belongs_to:"addresses"`
	SITDestinationFinalAddressID  *uuid.UUID                     `db:"sit_destination_final_address_id"`
	Description                   *string                        `db:"description"`
	Dimensions                    MTOServiceItemDimensions       `has_many:"mto_service_item_dimensions" fk_id:"mto_service_item_id"`
	CustomerContacts              MTOServiceItemCustomerContacts `has_many:"mto_service_item_customer_contacts" fk_id:"mto_service_item_id"`
	CreatedAt                     time.Time                      `db:"created_at"`
	UpdatedAt                     time.Time                      `db:"updated_at"`
	ApprovedAt                    *time.Time                     `db:"approved_at"`
	RejectedAt                    *time.Time                     `db:"rejected_at"`
}

// MTOServiceItems is a slice containing MTOServiceItems
type MTOServiceItems []MTOServiceItem

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOServiceItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(MTOServiceItemStatusSubmitted),
		string(MTOServiceItemStatusApproved),
		string(MTOServiceItemStatusRejected),
	}})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	vs = append(vs, &OptionalUUIDIsPresent{Field: m.MTOShipmentID, Name: "MTOShipmentID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.ReServiceID, Name: "ReServiceID"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.Reason, Name: "Reason"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.PickupPostalCode, Name: "PickupPostalCode"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.Description, Name: "Description"})

	return validate.Validate(vs...), nil
}

// TableName overrides the table name used by Pop.
func (m MTOServiceItem) TableName() string {
	return "mto_service_items"
}
