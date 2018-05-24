package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID                         uuid.UUID                    `json:"id" db:"id"`
	MoveID                     uuid.UUID                    `json:"move_id" db:"move_id"`
	Move                       Move                         `belongs_to:"move"`
	CreatedAt                  time.Time                    `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                    `json:"updated_at" db:"updated_at"`
	Size                       *internalmessages.TShirtSize `json:"size" db:"size"`
	WeightEstimate             *int64                       `json:"weight_estimate" db:"weight_estimate"`
	EstimatedIncentive         *string                      `json:"estimated_incentive" db:"estimated_incentive"`
	PlannedMoveDate            *time.Time                   `json:"planned_move_date" db:"planned_move_date"`
	PickupPostalCode           *string                      `json:"pickup_postal_code" db:"pickup_postal_code"`
	HasAdditionalPostalCode    *bool                        `json:"has_additional_postal_code" db:"has_additional_postal_code"`
	AdditionalPickupPostalCode *string                      `json:"additional_pickup_postal_code" db:"additional_pickup_postal_code"`
	DestinationPostalCode      *string                      `json:"destination_postal_code" db:"destination_postal_code"`
	HasSit                     *bool                        `json:"has_sit" db:"has_sit"`
	DaysInStorage              *int64                       `json:"days_in_storage" db:"days_in_storage"`
	Status                     PPMStatus                    `json:"status" db:"status"`
	HasRequestedAdvance        bool                         `json:"has_requested_advance" db:"has_requested_advance"`
	AdvanceID                  *uuid.UUID                   `json:"advance_id" db:"advance_id"`
	Advance                    *Reimbursement               `belongs_to:"reimbursements"`
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

// FetchPersonallyProcuredMove Fetches and Validates a PPM model
func FetchPersonallyProcuredMove(db *pop.Connection, authUser User, reqApp string, id uuid.UUID) (*PersonallyProcuredMove, error) {
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
	if reqApp == app.MyApp && ppm.Move.Orders.ServiceMember.UserID != authUser.ID {
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

		if ppm.Advance != nil {
			if verrs, err := db.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = err
				return transactionError
			}
			ppm.AdvanceID = &ppm.Advance.ID
		}

		if verrs, err := db.ValidateAndSave(ppm); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
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
