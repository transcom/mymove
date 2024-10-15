package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type State struct {
	ID        uuid.UUID `json:"id" db:"id"`
	State     string    `json:"state" db:"state"`
	StateName string    `json:"state_name" db:"state_name"`
	IsOconus  bool      `json:"is_oconus" db:"is_oconus"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (s State) TableName() string {
	return "re_states"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s State) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(s.State), Name: "State"},
		&validators.StringIsPresent{Field: string(s.StateName), Name: "StateName"},
		&BooleanMustBeSet{Field: &s.IsOconus, Name: "IsOconus"},
	), nil
}

// fetches countries by the two digit code
func FetchStateByCode(db *pop.Connection, code string) (State, error) {
	var state State
	err := db.Where("state = ?", code).First(&state)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return State{}, errors.Wrap(ErrFetchNotFound, "the state code provided in the request was not found")
		}
		return State{}, err
	}

	return state, nil
}

// fetches countries by the two digit code
func FetchStateByID(db *pop.Connection, id uuid.UUID) (State, error) {
	var state State
	err := db.Q().Find(&state, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return State{}, ErrFetchNotFound
		}
		return State{}, err
	}

	return state, nil
}
