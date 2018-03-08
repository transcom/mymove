package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

// Move is an object representing a move
type Move struct {
	ID               uuid.UUID `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	SelectedMoveType *string   `json:"selected_move_type" db:"selected_move_type"`
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

// GetOrCreateMove gets or creates a move
func GetOrCreateMove(db *pop.Connection, userID uuid.UUID) (Move, error) {
	// Check if move already exists
	query := db.Where("user_id = $1", userID)
	var moves []Move
	err := query.All(&moves)
	if err != nil {
		err = errors.Wrap(err, "DB Query Error")
		return (Move{}), err
	}

	// If move is not in DB, create it
	if len(moves) == 0 {
		newMove := Move{
			UserID: userID,
		}
		verrs, err := db.ValidateAndCreate(&newMove)
		if verrs.HasAny() {
			return (Move{}), verrs
		} else if err != nil {
			err = errors.Wrap(err, "Unable to create move")
			return (Move{}), err
		}
		return newMove, nil
	}
	// one move was found, return it
	return moves[0], nil
}
