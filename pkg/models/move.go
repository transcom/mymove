package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"

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

// ModelFetchErrors describe the expected errors returned by model fetch methods
const (
	ModelFetchErrorNotFound      string = "NOT_FOUND"
	ModelFetchErrorNotAuthorized string = "NOT_AUTHORIZED"
)

// GetMoveForUser returns a move only if it is allowed for the given user to access that move.
// If the user is not authorized to access that move, it behaves as if no such move exists.
func GetMoveForUser(db *pop.Connection, userID uuid.UUID, id uuid.UUID) (Move, error) {
	move := Move{}
	err := db.Find(&move, id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			err = errors.Wrap(err, ModelFetchErrorNotFound)
		}
	} else {
		// TODO: Handle case where more than one user is authorized to modify move
		if move.UserID != userID {
			move = Move{} // make sure and return a blank move in this case.
			err = errors.Wrap(err, ModelFetchErrorNotAuthorized)
		}
	}

	return move, err
}

// GetMoveByID fetches a Move model by their database ID
func GetMoveByID(db *pop.Connection, id uuid.UUID) (Move, error) {
	move := Move{}
	err := db.Find(&move, id)
	return move, err
}
