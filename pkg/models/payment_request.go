package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type PaymentRequest struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	IsFinal         bool          `json:"is_final" db:"is_final"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID uuid.UUID     `json:"move_task_order_id" db:"move_task_order_id"`
	RejectionReason string        `json:"rejection_reason" db:"rejection_reason"`
	//TODO DocumentPackage
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentRequest) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.MoveTaskOrderID, Name: "MoveTaskOrderID"},
	), nil
}
