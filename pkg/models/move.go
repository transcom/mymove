package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// Move is an object representing a move
type Move struct {
	ID               uuid.UUID                         `json:"id" db:"id"`
	CreatedAt        time.Time                         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                         `json:"updated_at" db:"updated_at"`
	UserID           uuid.UUID                         `json:"user_id" db:"user_id"`
	SelectedMoveType internalmessages.SelectedMoveType `json:"selected_move_type" db:"selected_move_type"`
}

// String is not required by pop and may be deleted
func (m Move) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Moves is not required by pop and may be deleted
type Moves []Move

// String is not required by pop and may be deleted
func (m Moves) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Move) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.UserID, Name: "UserID"},
		&validators.StringIsPresent{Field: string(m.SelectedMoveType), Name: "SelectedMoveType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Move) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Move) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetMoveByID fetches a Move model by their database ID
func GetMoveByID(db *pop.Connection, id uuid.UUID) (Move, error) {
	move := Move{}
	err := db.Find(&move, id)
	return move, err
}

// GetMovesForUserID gets all move models for a given user ID
func GetMovesForUserID(db *pop.Connection, userID uuid.UUID) (Moves, error) {
	var moves Moves
	query := db.Where("user_id = $1", userID)
	err := query.All(&moves)
	return moves, err
}
