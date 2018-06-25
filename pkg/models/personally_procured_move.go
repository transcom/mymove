package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// PPMStatus represents the status of an order record's lifecycle
type PPMStatus string

const (
	// PPMStatusDRAFT captures enum value "DRAFT"
	PPMStatusDRAFT PPMStatus = "DRAFT"
	// PPMStatusSUBMITTED captures enum value "SUBMITTED"
	PPMStatusSUBMITTED PPMStatus = "SUBMITTED"
	// PPMStatusAPPROVED captures enum value "APPROVED"
	PPMStatusAPPROVED PPMStatus = "APPROVED"
	// PPMStatusINPROGRESS captures enum value "IN_PROGRESS"
	PPMStatusINPROGRESS PPMStatus = "IN_PROGRESS"
	// PPMStatusCANCELED captures enum value "CANCELED"
	PPMStatusCANCELED PPMStatus = "CANCELED"
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID                            uuid.UUID                    `json:"id" db:"id"`
	MoveID                        uuid.UUID                    `json:"move_id" db:"move_id"`
	Move                          Move                         `belongs_to:"move"`
	CreatedAt                     time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time                    `json:"updated_at" db:"updated_at"`
	Size                          *internalmessages.TShirtSize `json:"size" db:"size"`
	WeightEstimate                *int64                       `json:"weight_estimate" db:"weight_estimate"`
	PlannedMoveDate               *time.Time                   `json:"planned_move_date" db:"planned_move_date"`
	PickupPostalCode              *string                      `json:"pickup_postal_code" db:"pickup_postal_code"`
	HasAdditionalPostalCode       *bool                        `json:"has_additional_postal_code" db:"has_additional_postal_code"`
	AdditionalPickupPostalCode    *string                      `json:"additional_pickup_postal_code" db:"additional_pickup_postal_code"`
	DestinationPostalCode         *string                      `json:"destination_postal_code" db:"destination_postal_code"`
	HasSit                        *bool                        `json:"has_sit" db:"has_sit"`
	DaysInStorage                 *int64                       `json:"days_in_storage" db:"days_in_storage"`
	EstimatedStorageReimbursement *string                      `json:"estimated_storage_reimbursement" db:"estimated_storage_reimbursement"`
	Mileage                       *int64                       `json:"mileage" db:"mileage"`
	PlannedSITMax                 *unit.Cents                  `json:"planned_sit_max" db:"planned_sit_max"`
	SITMax                        *unit.Cents                  `json:"sit_max" db:"sit_max"`
	IncentiveEstimateMin          *unit.Cents                  `json:"incentive_estimate_min" db:"incentive_estimate_min"`
	IncentiveEstimateMax          *unit.Cents                  `json:"incentive_estimate_max" db:"incentive_estimate_max"`
	Status                        PPMStatus                    `json:"status" db:"status"`
	HasRequestedAdvance           bool                         `json:"has_requested_advance" db:"has_requested_advance"`
	AdvanceID                     *uuid.UUID                   `json:"advance_id" db:"advance_id"`
	Advance                       Reimbursement                `belongs_to:"reimbursements"`
	AdvanceWorksheet              Document                     `belongs_to:"documents"`
	AdvanceWorksheetID            *uuid.UUID                   `json:"advance_worksheet_id" db:"advance_worksheet_id"`
}

// PersonallyProcuredMoves is a list of PPMs
type PersonallyProcuredMoves []PersonallyProcuredMove

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(p.Status), Name: "Status"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PersonallyProcuredMove) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// State Machine
// Avoid calling PersonallyProcuredMove.Status = ... ever. Use these methods to change the state.

// Cancel cancels the PPM
func (p *PersonallyProcuredMove) Cancel() error {
	if p.Status != PPMStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	p.Status = PPMStatusCANCELED
	return nil
}

// FetchPersonallyProcuredMove Fetches and Validates a PPM model
func FetchPersonallyProcuredMove(db *pop.Connection, session *auth.Session, id uuid.UUID) (*PersonallyProcuredMove, error) {
	var ppm PersonallyProcuredMove
	err := db.Q().Eager("Move.Orders.ServiceMember", "Advance").Find(&ppm, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// TODO: Handle case where more than one user is authorized to modify ppm
	if session.IsMyApp() && ppm.Move.Orders.ServiceMember.ID != session.ServiceMemberID {
		return nil, ErrFetchForbidden
	}

	return &ppm, nil
}

// SavePersonallyProcuredMove Safely saves a PPM and it's associated Advance.
func SavePersonallyProcuredMove(db *pop.Connection, ppm *PersonallyProcuredMove) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if ppm.HasRequestedAdvance {
			if ppm.AdvanceID != nil {
				// GTCC isn't a valid method of receipt for PPM Advances, so reject if that's the case.
				if ppm.Advance.MethodOfReceipt == MethodOfReceiptGTCC {
					responseVErrors.Add("MethodOfReceipt", "GTCC is not a valid receipt method for PPM Advances.")
					return transactionError
				}

				if verrs, err := db.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Advance")
					return transactionError
				}
				ppm.AdvanceID = &ppm.Advance.ID
			} else if ppm.AdvanceID == nil {
				// if Has Requested Advance is set, but there is nothing saved or to save, that's an error.
				responseError = ErrInvalidPatchGate
				return transactionError
			}
		} else {
			if ppm.AdvanceID != nil {
				// If HasRequstedAdvance is false, we need to delete the record
				reimbursement := Reimbursement{}
				err := db.Find(&reimbursement, *ppm.AdvanceID)
				if err != nil {
					responseError = errors.Wrap(err, "Error finding Advance for Advance ID")
					return transactionError
				}
				ppm.Advance = reimbursement

				err = db.Destroy(ppm.Advance)
				if err != nil {
					responseError = errors.Wrap(err, "Error Deleting Advance record")
					return transactionError
				}
				ppm.AdvanceID = nil
				ppm.Advance = Reimbursement{}
			}
		}

		if verrs, err := db.ValidateAndSave(ppm); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving PPM")
			return transactionError
		}

		return nil

	})

	return responseVErrors, responseError
}

// createNewPPM adds a new Personally Procured Move record into the DB.
func createNewPPM(db *pop.Connection, moveID uuid.UUID) (*PersonallyProcuredMove, *validate.Errors, error) {
	ppm := PersonallyProcuredMove{
		MoveID: moveID,
		Status: PPMStatusDRAFT,
	}
	verrs, err := db.ValidateAndCreate(&ppm)
	if verrs.HasAny() {
		return nil, verrs, nil
	}
	if err != nil {
		return nil, verrs, err
	}

	return &ppm, verrs, nil

}
