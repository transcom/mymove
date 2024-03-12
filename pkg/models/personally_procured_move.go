package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/unit"
)

// ProGearStatus represents the status of a pro-gear question
type ProGearStatus string

const (
	// ProGearStatusYes captures enum value "YES"
	ProGearStatusYes ProGearStatus = "YES"
	// ProGearStatusNo captures enum value "NO"
	ProGearStatusNo ProGearStatus = "NO"
	// ProGearStatusNotSure captures enum value "YES"
	ProGearStatusNotSure ProGearStatus = "NOT SURE"
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
	// PPMStatusPAYMENTREQUESTED captures enum value "PAYMENT_REQUESTED"
	PPMStatusPAYMENTREQUESTED PPMStatus = "PAYMENT_REQUESTED"
	// PPMStatusCOMPLETED captures enum value "COMPLETED"
	PPMStatusCOMPLETED PPMStatus = "COMPLETED"
	// PPMStatusCANCELED captures enum value "CANCELED"
	PPMStatusCANCELED PPMStatus = "CANCELED"
)

// PersonallyProcuredMove is the portion of a move that a service member performs themselves
type PersonallyProcuredMove struct {
	ID                            uuid.UUID      `json:"id" db:"id"`
	MoveID                        uuid.UUID      `json:"move_id" db:"move_id"`
	Move                          Move           `belongs_to:"move" fk_id:"move_id"`
	CreatedAt                     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time      `json:"updated_at" db:"updated_at"`
	WeightEstimate                *unit.Pound    `json:"weight_estimate" db:"weight_estimate"`
	OriginalMoveDate              *time.Time     `json:"original_move_date" db:"original_move_date"`
	ActualMoveDate                *time.Time     `json:"actual_move_date" db:"actual_move_date"`
	SubmitDate                    *time.Time     `json:"submit_date" db:"submit_date"`
	ApproveDate                   *time.Time     `json:"approve_date" db:"approve_date"`
	ReviewedDate                  *time.Time     `json:"reviewed_date" db:"reviewed_date"`
	NetWeight                     *unit.Pound    `json:"net_weight" db:"net_weight"`
	PickupPostalCode              *string        `json:"pickup_postal_code" db:"pickup_postal_code"`
	HasAdditionalPostalCode       *bool          `json:"has_additional_postal_code" db:"has_additional_postal_code"`
	AdditionalPickupPostalCode    *string        `json:"additional_pickup_postal_code" db:"additional_pickup_postal_code"`
	DestinationPostalCode         *string        `json:"destination_postal_code" db:"destination_postal_code"`
	HasSit                        *bool          `json:"has_sit" db:"has_sit"`
	DaysInStorage                 *int64         `json:"days_in_storage" db:"days_in_storage"`
	EstimatedStorageReimbursement *string        `json:"estimated_storage_reimbursement" db:"estimated_storage_reimbursement"`
	Mileage                       *int64         `json:"mileage" db:"mileage"`
	PlannedSITMax                 *unit.Cents    `json:"planned_sit_max" db:"planned_sit_max"`
	SITMax                        *unit.Cents    `json:"sit_max" db:"sit_max"`
	IncentiveEstimateMin          *unit.Cents    `json:"incentive_estimate_min" db:"incentive_estimate_min"`
	IncentiveEstimateMax          *unit.Cents    `json:"incentive_estimate_max" db:"incentive_estimate_max"`
	Status                        PPMStatus      `json:"status" db:"status"`
	HasRequestedAdvance           bool           `json:"has_requested_advance" db:"has_requested_advance"`
	AdvanceID                     *uuid.UUID     `json:"advance_id" db:"advance_id"`
	Advance                       *Reimbursement `belongs_to:"reimbursements" fk_id:"advance_id"`
	AdvanceWorksheet              Document       `belongs_to:"documents" fk_id:"advance_worksheet_id"`
	AdvanceWorksheetID            *uuid.UUID     `json:"advance_worksheet_id" db:"advance_worksheet_id"`
	TotalSITCost                  *unit.Cents    `json:"total_sit_cost" db:"total_sit_cost"`
	HasProGear                    *ProGearStatus `json:"has_pro_gear" db:"has_pro_gear"`
	HasProGearOverThousand        *ProGearStatus `json:"has_pro_gear_over_thousand" db:"has_pro_gear_over_thousand"`
}

// TableName overrides the table name used by Pop.
func (p PersonallyProcuredMove) TableName() string {
	return "personally_procured_moves"
}

// PersonallyProcuredMoves is a list of PPMs
type PersonallyProcuredMoves []PersonallyProcuredMove

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PersonallyProcuredMove) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(p.Status), Name: "Status"},
	), nil
}

// State Machinery
// Avoid calling PersonallyProcuredMove.Status = ... ever. Use these methods to change the state.

// Submit marks the PPM request for review
func (p *PersonallyProcuredMove) Submit(submitDate time.Time) error {
	p.Status = PPMStatusSUBMITTED
	p.SubmitDate = &submitDate
	return nil
}

// Approve approves the PPM to go forward.
func (p *PersonallyProcuredMove) Approve(approveDate time.Time) error {
	p.Status = PPMStatusAPPROVED
	p.ApproveDate = &approveDate
	return nil
}

// RequestPayment requests payment for the PPM
func (p *PersonallyProcuredMove) RequestPayment() error {
	p.Status = PPMStatusPAYMENTREQUESTED
	return nil
}

// Complete marks the PPM as completed
func (p *PersonallyProcuredMove) Complete(reviewedDate time.Time) error {
	p.Status = PPMStatusCOMPLETED
	p.ReviewedDate = &reviewedDate
	return nil
}

// Cancel marks the PPM as Canceled
func (p *PersonallyProcuredMove) Cancel() error {
	p.Status = PPMStatusCANCELED
	return nil
}

// FetchPersonallyProcuredMove Fetches and Validates a PPM model
func FetchPersonallyProcuredMove(db *pop.Connection, _ *auth.Session, id uuid.UUID) (*PersonallyProcuredMove, error) {
	var ppm PersonallyProcuredMove
	err := db.Q().Eager("Move.Orders.OriginDutyLocation.Address", "Move.Orders.NewDutyLocation.Address", "Advance").Find(&ppm, id)
	if err != nil {
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	return &ppm, nil
}

// FetchPersonallyProcuredMoveByOrderID Fetches and Validates a PPM model
func FetchPersonallyProcuredMoveByOrderID(db *pop.Connection, orderID uuid.UUID) (*PersonallyProcuredMove, error) {
	var ppm PersonallyProcuredMove
	err := db.Q().
		LeftJoin("moves as m", "m.id = personally_procured_moves.move_id").
		Where("m.orders_id = ?", orderID).
		First(&ppm)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return &PersonallyProcuredMove{}, ErrFetchNotFound
		}
		return &PersonallyProcuredMove{}, err
	}

	return &ppm, nil
}

// SavePersonallyProcuredMove Safely saves a PPM and it's associated Advance.
func SavePersonallyProcuredMove(db *pop.Connection, ppm *PersonallyProcuredMove) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionErr := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if ppm.HasRequestedAdvance {
			if ppm.Advance != nil {
				if verrs, err := db.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Advance")
					return transactionError
				}
				ppm.AdvanceID = &ppm.Advance.ID
			}
		}

		if verrs, err := db.ValidateAndSave(ppm); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving PPM")
			return transactionError
		}

		return nil

	})

	if transactionErr != nil {
		return responseVErrors, responseError
	}

	return responseVErrors, responseError
}
