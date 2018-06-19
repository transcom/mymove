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
	"github.com/transcom/mymove/pkg/auth"
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
	// MoveStatusCANCELED captures enum value "CANCELED"
	MoveStatusCANCELED MoveStatus = "CANCELED"
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
	CancelReason            *string                            `json:"cancel_reason" db:"cancel_reason"`
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

// State Machine
// Avoid calling Move.Status = ... ever. Use these methods to change the state.

// Submit submits the Move
func (m *Move) Submit() error {
	if m.Status != MoveStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Submit")
	}

	m.Status = MoveStatusSUBMITTED

	//TODO: update PPM status too
	// for i, _ := range m.PersonallyProcuredMoves {
	// 	err := m.PersonallyProcuredMoves[i].Submit()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	for _, ppm := range m.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Cancel cancels the Move and its associated PPMs
func (m *Move) Cancel(reason string) error {
	if m.Status != MoveStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	m.Status = MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if reason != "" {
		m.CancelReason = &reason
	}

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range m.PersonallyProcuredMoves {
		err := m.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

// FetchMove fetches and validates a Move for this User
func FetchMove(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Move, error) {
	var move Move
	err := db.Q().Eager("PersonallyProcuredMoves.Advance", "SignedCertifications").Find(&move, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	// Ensure that the logged-in user is authorized to access this move
	_, authErr := FetchOrder(db, session, move.OrdersID)
	if authErr != nil {
		return nil, authErr
	}

	return &move, nil
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	size *internalmessages.TShirtSize,
	weightEstimate *int64,
	plannedMoveDate *time.Time,
	pickupPostalCode *string,
	hasAdditionalPostalCode *bool,
	additionalPickupPostalCode *string,
	destinationPostalCode *string,
	hasSit *bool,
	daysInStorage *int64,
	estimatedStorageReimbursement *string,
	hasRequestedAdvance bool,
	advance *Reimbursement) (*PersonallyProcuredMove, *validate.Errors, error) {

	newPPM := PersonallyProcuredMove{
		MoveID:                     m.ID,
		Move:                       m,
		Size:                       size,
		WeightEstimate:             weightEstimate,
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
		EstimatedStorageReimbursement: estimatedStorageReimbursement,
	}

	verrs, err := SavePersonallyProcuredMove(db, &newPPM)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newPPM, verrs, nil
}

// CreateSignedCertification creates a new SignedCertification associated with this move
func (m Move) CreateSignedCertification(db *pop.Connection,
	submittingUserID uuid.UUID,
	certificationText string,
	signature string,
	date time.Time) (*SignedCertification, *validate.Errors, error) {

	newSignedCertification := SignedCertification{
		MoveID:            m.ID,
		SubmittingUserID:  submittingUserID,
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

// SaveMoveStatuses safely saves a Move status and its associated PPMs' Advances' statuses.
// TODO: Add functionality to save more than just status on these objects
func SaveMoveStatuses(db *pop.Connection, move *Move) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		for _, ppm := range move.PersonallyProcuredMoves {
			if ppm.Advance != nil {
				if verrs, err := db.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Advance")
					return transactionError
				}
			}
			// TODO: Add back in once we are updating PPM Status
			// if verrs, err := db.ValidateAndSave(ppm); verrs.HasAny() || err != nil {
			// 	responseVErrors.Append(verrs)
			// 	responseError = errors.Wrap(err, "Error Saving PPM")
			// 	return transactionError
			// }
		}

		if verrs, err := db.ValidateAndSave(move); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move")
			return transactionError
		}

		return nil

	})

	return responseVErrors, responseError
}

// FetchMoveForAdvancePaperwork returns a Move with all of the associations required
// to generate the Advance paperwork.
func FetchMoveForAdvancePaperwork(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	if err := db.Q().Eager("Orders.NewDutyStation", "Orders.ServiceMember.BackupContacts", "PersonallyProcuredMoves.Advance").Find(&move, moveID); err != nil {
		return move, errors.Wrap(err, "could not load move")
	}
	return move, nil
}
