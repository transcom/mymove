package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type PaymentRequest struct {
	ID              uuid.UUID     `db:"id"`
	IsFinal         bool          `db:"is_final"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID uuid.UUID     `db:"move_task_order_id"`
	ServiceItemIDs  []uuid.UUID   `db:"service_item_id_s"`
	RejectionReason string        `db:"rejection_reason"`
	//TODO DocumentPackage
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentRequest) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.MoveTaskOrderID, Name: "MoveTaskOrderID"},
		&UUIDArrayIsPresent{Field: p.ServiceItemIDs, Name: "ServiceItemIDs"},
		// TODO: make sure serviceItemIDs are unique
		&validators.StringIsPresent{Field: p.RejectionReason, Name: "RejectionReason"},
	), nil
}
