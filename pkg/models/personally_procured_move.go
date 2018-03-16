package models

import (
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID             uuid.UUID                    `json:"id" db:"id"`
	MoveID         uuid.UUID                    `json:"move_id" db:"move_id"`
	CreatedAt      time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at" db:"updated_at"`
	Size           *internalmessages.TShirtSize `json:"size" db:"size"`
	WeightEstimate *int64                       `json:"weight_estimate" db:"weight_estimate"`
}

// PersonallyProcuredMoves is a list of PPMs
type PersonallyProcuredMoves []PersonallyProcuredMove

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
