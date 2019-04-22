package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// AccessCode is an object representing an access code for a service member
type AccessCode struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	MoveID    uuid.UUID        `json:"move_id" db:"move_id"`
	Code      string           `json:"code" db:"code"`
	MoveType  SelectedMoveType `json:"move_type" db:"move_type"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}

//ValidateAccessCode checks whether an access code is valid and unused
func ValidateAccessCode(db *pop.Connection, code string, moveType SelectedMoveType) (*AccessCode, bool) {
	ac := AccessCode{}
	err := db.
		Where("code = ?", code).
		Where("move_id IS NULL").
		Where("move_type = ?", moveType).
		First(&ac)
	if err != nil {
		return &ac, false
	}
	return &ac, true
}
