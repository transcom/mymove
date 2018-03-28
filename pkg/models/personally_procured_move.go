package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID             uuid.UUID                    `json:"id" db:"id"`
	MoveID         uuid.UUID                    `json:"move_id" db:"move_id"`
	Move           Move                         `belongs_to:"move"`
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

// GetPersonallyProcuredMovesForMoveID gets all PPMs models for a given move ID
func GetPersonallyProcuredMovesForMoveID(db *pop.Connection, moveID uuid.UUID) (PersonallyProcuredMoves, error) {
	var ppms PersonallyProcuredMoves
	query := db.Where("move_id = $1", moveID)
	err := query.All(&ppms)
	return ppms, err
}

// GetPersonallyProcuredMoveForID returns a PersonallyProcuredMove model for a given ID
func GetPersonallyProcuredMoveForID(db *pop.Connection, id uuid.UUID) (PersonallyProcuredMove, error) {
	var ppm PersonallyProcuredMove
	err := db.Find(&ppm, id)
	return ppm, err
}
