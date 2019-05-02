package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

// AccessCode is an object representing an access code for a service member
type AccessCode struct {
	ID              uuid.UUID         `json:"id" db:"id"`
	ServiceMemberID *uuid.UUID        `json:"service_member_id" db:"service_member_id"`
	Code            string            `json:"code" db:"code"`
	MoveType        *SelectedMoveType `json:"move_type" db:"move_type"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	ClaimedAt       *time.Time        `json:"claimed_at" db:"claimed_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (ac *AccessCode) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(ac.Code), Name: "Code"},
		&validators.StringIsPresent{Field: string(*ac.MoveType), Name: "MoveType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (ac *AccessCode) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (ac *AccessCode) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
