package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"crypto/sha256"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// MoveStatus represents the status of an order record's lifecycle
type MoveStatus string

const (
	// MoveStatusDRAFT captures enum value "DRAFT"
	MoveStatusDRAFT MoveStatus = "DRAFT"
	// MoveStatusSUBMITTED captures enum value "SUBMITTED"
	MoveStatusSUBMITTED MoveStatus = "SUBMITTED"
	// MoveStatusAPPROVED captures enum value "APPROVED"
	MoveStatusAPPROVED MoveStatus = "APPROVED"
	// MoveStatusCOMPLETED captures enum value "COMPLETED"
	MoveStatusCOMPLETED MoveStatus = "COMPLETED"
)

const maxLocatorAttempts = 3
const locatorLength = 6

var locatorLetters = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZ")

// Move is an object representing a move
type Move struct {
	ID                      uuid.UUID                          `json:"id" db:"id"`
	Locator                 string                             `json:"locator" db:"locator"`
	CreatedAt               time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                          `json:"updated_at" db:"updated_at"`
	OrdersID                uuid.UUID                          `json:"orders_id" db:"orders_id"`
	Orders                  Order                              `belongs_to:"orders"`
	SelectedMoveType        *internalmessages.SelectedMoveType `json:"selected_move_type" db:"selected_move_type"`
	PersonallyProcuredMoves PersonallyProcuredMoves            `has_many:"personally_procured_moves" order_by:"created_at desc"`
	Status                  MoveStatus                         `json:"status" db:"status"`
	SignedCertifications    SignedCertifications               `has_many:"signed_certifications" order_by:"created_at desc"`
}

// Moves is not required by pop and may be deleted
type Moves []Move

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Move) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.OrdersID, Name: "OrdersID"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
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
func FetchMove(db *pop.Connection, authUser User, reqApp string, id uuid.UUID) (*Move, error) {
	var move Move
	err := db.Q().Eager().Find(&move, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	// Ensure that the logged-in user is authorized to access this move
	_, authErr := FetchOrder(db, authUser, reqApp, move.OrdersID)
	if authErr != nil {
		return nil, authErr
	}

	return &move, nil
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	size *internalmessages.TShirtSize,
	weightEstimate *int64,
	estimatedIncentive *string,
	plannedMoveDate *time.Time,
	pickupPostalCode *string,
	hasAdditionalPostalCode *bool,
	additionalPickupPostalCode *string,
	destinationPostalCode *string,
	hasSit *bool,
	daysInStorage *int64,
	hasRequestedAdvance bool,
	advance *Reimbursement) (*PersonallyProcuredMove, *validate.Errors, error) {

	newPPM := PersonallyProcuredMove{
		MoveID:                     m.ID,
		Move:                       m,
		Size:                       size,
		WeightEstimate:             weightEstimate,
		EstimatedIncentive:         estimatedIncentive,
		PlannedMoveDate:            plannedMoveDate,
		PickupPostalCode:           pickupPostalCode,
		HasAdditionalPostalCode:    hasAdditionalPostalCode,
		AdditionalPickupPostalCode: additionalPickupPostalCode,
		DestinationPostalCode:      destinationPostalCode,
		HasSit:                     hasSit,
		DaysInStorage:              daysInStorage,
		Status:                     PPMStatusDRAFT,
		HasRequestedAdvance:        hasRequestedAdvance,
		Advance:                    advance,
	}

	verrs, err := SavePersonallyProcuredMove(db, &newPPM)
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

// generateLocator constructs a record locator - a unique 6 character alphanumeric string
func generateLocator() string {
	// Get a UUID as a source of (almost certainly) unique bytes
	seed, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	// Scramble them via SHA256 in case UUID has structure
	scrambledBytes := sha256.Sum256(seed.Bytes())
	// Now convert bytes to letters
	locatorRunes := make([]rune, locatorLength)
	for idx := 0; idx < locatorLength; idx++ {
		j := int(scrambledBytes[idx]) % len(locatorLetters)
		locatorRunes[idx] = locatorLetters[j]
	}
	return string(locatorRunes)
}

// createNewMove adds a new Move record into the DB. In the (unlikely) event that we have a clash on Locators we
// retry with a new record locator.
func createNewMove(db *pop.Connection,
	ordersID uuid.UUID,
	selectedType *internalmessages.SelectedMoveType) (*Move, *validate.Errors, error) {

	for i := 0; i < maxLocatorAttempts; i++ {
		move := Move{
			OrdersID:         ordersID,
			Locator:          generateLocator(),
			SelectedMoveType: selectedType,
			Status:           MoveStatusDRAFT,
		}
		verrs, err := db.ValidateAndCreate(&move)
		if verrs.HasAny() {
			return nil, verrs, nil
		}
		if err != nil {
			if strings.HasPrefix(errors.Cause(err).Error(), uniqueConstraintViolationErrorPrefix) {
				// If we have a collision, try again for maxLocatorAttempts
				continue
			}
			return nil, verrs, err
		}

		return &move, verrs, nil
	}
	// the only way we get here is if we got a unique constraint error maxLocatorAttempts times.
	verrs := validate.NewErrors()
	return nil, verrs, ErrLocatorGeneration
}
