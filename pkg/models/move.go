package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// Move is an object representing a move
type Move struct {
	ID                      uuid.UUID                          `json:"id" db:"id"`
	CreatedAt               time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                          `json:"updated_at" db:"updated_at"`
	UserID                  uuid.UUID                          `json:"user_id" db:"user_id"`
	User                    User                               `belongs_to:"user"`
	SelectedMoveType        *internalmessages.SelectedMoveType `json:"selected_move_type" db:"selected_move_type"`
	PersonallyProcuredMoves PersonallyProcuredMoves            `has_many:"personally_procured_moves"`
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

// MoveResult is returned by GetMoveForUser and encapsulates whether the call succeeded and why it failed.
type MoveResult struct {
	valid     bool
	errorCode FetchError
	move      Move
}

// IsValid indicates whether the MoveResult is valid.
func (m MoveResult) IsValid() bool {
	return m.valid
}

// Move returns the move if and only if the move was correctly fetched
func (m MoveResult) Move() Move {
	if !m.valid {
		zap.L().Fatal("Check if this isValid before accessing the Move()!")
	}
	return m.move
}

// ErrorCode returns the error if and only if the move was not correctly fetched
func (m MoveResult) ErrorCode() FetchError {
	if m.valid {
		zap.L().Fatal("Check that this !isValid before accessing the ErrorCode()!")
	}
	return m.errorCode
}

// NewInvalidMoveResult creates an invalid MoveResult
func NewInvalidMoveResult(errorCode FetchError) MoveResult {
	return MoveResult{
		errorCode: errorCode,
	}
}

// NewValidMoveResult creates a valid MoveResult
func NewValidMoveResult(move Move) MoveResult {
	return MoveResult{
		valid: true,
		move:  move,
	}
}

// FetchMove fetches and validates a Move for this User
func FetchMove(db *pop.Connection, authUser User, id uuid.UUID) (*Move, error) {
	var move Move
	err := db.Q().Eager().Find(&move, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// TODO: Handle case where more than one user is authorized to modify move
	if move.UserID != authUser.ID {
		return nil, ErrFetchForbidden
	}

	return &move, nil
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	size *internalmessages.TShirtSize,
	weightEstimate *int64,
	estimatedIncentive *string,
	plannedMoveDate *time.Time,
	pickupZip *string,
	additionalPickupZip *string,
	destinationZip *string,
	daysInStorage *int64) (*PersonallyProcuredMove, *validate.Errors, error) {

	newPPM := PersonallyProcuredMove{
		MoveID:              m.ID,
		Move:                m,
		Size:                size,
		WeightEstimate:      weightEstimate,
		EstimatedIncentive:  estimatedIncentive,
		PlannedMoveDate:     plannedMoveDate,
		PickupZip:           pickupZip,
		AdditionalPickupZip: additionalPickupZip,
		DestinationZip:      destinationZip,
		DaysInStorage:       daysInStorage,
	}

	verrs, err := db.ValidateAndCreate(&newPPM)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newPPM, verrs, nil
}

// GetMoveForUser returns a move only if it is allowed for the given user to access that move.
// If the user is not authorized to access that move, it behaves as if no such move exists.
func GetMoveForUser(db *pop.Connection, userID uuid.UUID, id uuid.UUID) (MoveResult, error) {
	var result MoveResult
	var move Move
	err := db.Find(&move, id)
	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			result = NewInvalidMoveResult(FetchErrorNotFound)
			err = nil
		}
		// Otherwise, it's an unexpected err so we return that.
	} else {
		// TODO: Handle case where more than one user is authorized to modify move
		if move.UserID != userID {
			result = NewInvalidMoveResult(FetchErrorForbidden)
		} else {
			result = NewValidMoveResult(move)
		}
	}

	return result, err
}

// GetMovesForUserID gets all move models for a given user ID
func GetMovesForUserID(db *pop.Connection, userID uuid.UUID) (Moves, error) {
	var moves Moves
	query := db.Where("user_id = $1", userID)
	err := query.All(&moves)
	return moves, err
}
