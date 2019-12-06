package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// MtoServiceItem is an object representing service items for a move task order
type MtoServiceItem struct {
	ID              uuid.UUID     `db:"id"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID uuid.UUID     `db:"move_task_order_id"`
	MtoShipment     MtoShipment   `belongs_to:"mto_shipments"`
	MtoShipmentID   uuid.UUID     `db:"mto_shipment_id"`
	ReService       ReService     `belongs_to:"re_services"`
	ReServiceID     uuid.UUID     `db:"re_service_id"`
	MetaID          uuid.UUID     `db:"meta_id"`
	MetaType        string        `db:"meta_type"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MtoServiceItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MtoShipmentID, Name: "MtoShipmentID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.ReServiceID, Name: "ReServiceID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MetaID, Name: "MetaID"})
	vs = append(vs, &validators.StringIsPresent{Field: m.MetaType, Name: "MetaType"})
	return validate.Validate(vs...), nil
}
