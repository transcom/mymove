package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// Move is an object representing a move
type Move struct {
	ID                      uuid.UUID                          `json:"id" db:"id"`
	CreatedAt               time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                          `json:"updated_at" db:"updated_at"`
	OrdersID                uuid.UUID                          `json:"orders_id" db:"orders_id"`
	Order                   Order                              `belongs_to:"order"`
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
		&validators.UUIDIsPresent{Field: m.OrdersID, Name: "OrdersID"},
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

// FetchMove fetches and validates a Move for this User
func FetchMove(db *pop.Connection, authUser User, id uuid.UUID) (*Move, error) {
	var move Move
	err := db.Q().Eager().Find(&move, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// Fetch orders from move.orders id. Get servicemember from orders. Make sure user id matches.
	// TODO: Handle case where more than one user is authorized to modify move
	// if move.UserID != authUser.ID {
	// 	return nil, ErrFetchForbidden
	// }

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

// CreateSignedCertification creates a new SignedCertification associated with this move
func (m Move) CreateSignedCertification(db *pop.Connection,
	submittingUser User,
	certificationText string,
	signature string,
	date time.Time) (*SignedCertification, *validate.Errors, error) {

	newSignedCertification := SignedCertification{
		MoveID:            m.ID,
		SubmittingUserID:  submittingUser.ID,
		CertificationText: certificationText,
		Signature:         signature,
		Date:              date,
	}

	verrs, err := db.ValidateAndCreate(&newSignedCertification)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newSignedCertification, verrs, nil
}

// GetMovesForUserID gets all move models for a given user ID
func GetMovesForUserID(db *pop.Connection, userID uuid.UUID) (Moves, error) {
	var moves Moves
	query := db.Where("user_id = $1", userID)
	err := query.All(&moves)
	return moves, err
}
