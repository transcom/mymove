package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// MTOServiceItem is an object representing service items for a move task order.
type MTOServiceItem struct {
	ID               uuid.UUID     `db:"id"`
	MoveTaskOrder    MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID  uuid.UUID     `db:"move_task_order_id"`
	MTOShipment      MTOShipment   `belongs_to:"mto_shipments"`
	MTOShipmentID    *uuid.UUID    `db:"mto_shipment_id"`
	ReService        ReService     `belongs_to:"re_services"`
	ReServiceID      uuid.UUID     `db:"re_service_id"`
	Reason           *string       `db:"reason"`
	PickupPostalCode *string       `db:"pickup_postal_code"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
}

// MTOServiceItems is a slice containing MTOServiceItems
type MTOServiceItems []MTOServiceItem

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOServiceItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	vs = append(vs, &OptionalUUIDIsPresent{Field: m.MTOShipmentID, Name: "MTOShipmentID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.ReServiceID, Name: "ReServiceID"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.Reason, Name: "Reason"})
	vs = append(vs, &StringIsNilOrNotBlank{Field: m.PickupPostalCode, Name: "PickupPostalCode"})

	return validate.Validate(vs...), nil
}

// TableName overrides the table name used by Pop.
func (m MTOServiceItem) TableName() string {
	return "mto_service_items"
}
