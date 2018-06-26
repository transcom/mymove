package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/unit"
)

// MethodOfReceipt is how the SM will be paid
type MethodOfReceipt string

const (
	// MethodOfReceiptMILPAY captures enum value MIL_PAY
	MethodOfReceiptMILPAY MethodOfReceipt = "MIL_PAY"
	// MethodOfReceiptOTHERDD captures enum value OTHER_DD
	MethodOfReceiptOTHERDD MethodOfReceipt = "OTHER_DD"
	// MethodOfReceiptGTCC captures enum value GTCC
	MethodOfReceiptGTCC MethodOfReceipt = "GTCC"
)

// ReimbursementStatus is the status of the Reimbursement
type ReimbursementStatus string

const (
	// ReimbursementStatusDRAFT captures enum value "DRAFT"
	ReimbursementStatusDRAFT ReimbursementStatus = "DRAFT"
	// ReimbursementStatusREQUESTED captures enum value "REQUESTED"
	ReimbursementStatusREQUESTED ReimbursementStatus = "REQUESTED"
	// ReimbursementStatusAPPROVED captures enum value "APPROVED"
	ReimbursementStatusAPPROVED ReimbursementStatus = "APPROVED"
	// ReimbursementStatusREJECTED captures enum value "REJECTED"
	ReimbursementStatusREJECTED ReimbursementStatus = "REJECTED"
	// ReimbursementStatusPAID captures enum value "PAID"
	ReimbursementStatusPAID ReimbursementStatus = "PAID"
)

// Reimbursement is money that is intended to be paid to the servicemember
type Reimbursement struct {
	ID              uuid.UUID           `json:"id" db:"id"`
	CreatedAt       time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" db:"updated_at"`
	RequestedAmount unit.Cents          `json:"requested_amount" db:"requested_amount"`
	MethodOfReceipt MethodOfReceipt     `json:"method_of_receipt" db:"method_of_receipt"`
	Status          ReimbursementStatus `json:"status" db:"status"`
	RequestedDate   *time.Time          `json:"requested_date" db:"requested_date"`
}

// State Machine
// Avoid calling Reimbursement.Status = ... ever. Use these methods to change the state.

// ErrInvalidTransition is an error representing an invalid transition.
var ErrInvalidTransition = errors.New("INVALID_TRANSITION")

// Request officially requests the reimbursement.
func (r *Reimbursement) Request() error {
	if r.Status != ReimbursementStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Request")
	}

	r.Status = ReimbursementStatusREQUESTED
	today := time.Now()
	r.RequestedDate = &today
	return nil
}

// Approve approves the Reimbursement
func (r *Reimbursement) Approve() error {
	if r.Status != ReimbursementStatusREQUESTED {
		return errors.Wrap(ErrInvalidTransition, "Approve")
	}

	r.Status = ReimbursementStatusAPPROVED
	return nil
}

// Reject rejects the Reimbursement
func (r *Reimbursement) Reject() error {
	if r.Status != ReimbursementStatusREQUESTED {
		return errors.Wrap(ErrInvalidTransition, "Reject")
	}

	r.Status = ReimbursementStatusREJECTED
	return nil
}

// Pay pays the Reimbursement
func (r *Reimbursement) Pay() error {
	if r.Status != ReimbursementStatusAPPROVED {
		return errors.Wrap(ErrInvalidTransition, "Pay")
	}

	r.Status = ReimbursementStatusPAID
	return nil
}

// END State Machine

// BuildDraftReimbursement makes a Reimbursement in the DRAFT state, but does not save it
func BuildDraftReimbursement(requestedAmount unit.Cents, methodOfReceipt MethodOfReceipt) Reimbursement {
	return Reimbursement{
		Status:          ReimbursementStatusDRAFT,
		RequestedAmount: requestedAmount,
		MethodOfReceipt: methodOfReceipt,
	}
}

// BuildRequestedReimbursement makes a Reimbursement in the REQUEST state, but does not save it
// This will be useful for reimbursements that are filed after the initial move is created
func BuildRequestedReimbursement(requestedAmount unit.Cents, methodOfReceipt MethodOfReceipt) Reimbursement {
	today := time.Now()
	return Reimbursement{
		Status:          ReimbursementStatusREQUESTED,
		RequestedAmount: requestedAmount,
		MethodOfReceipt: methodOfReceipt,
		RequestedDate:   &today,
	}
}

// Reimbursements is not required by pop and may be deleted
type Reimbursements []Reimbursement

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Reimbursement) Validate(tx *pop.Connection) (*validate.Errors, error) {
	if r == nil {
		return validate.NewErrors(), nil
	}

	validStatuses := []string{
		string(ReimbursementStatusDRAFT),
		string(ReimbursementStatusREQUESTED),
		string(ReimbursementStatusAPPROVED),
		string(ReimbursementStatusREJECTED),
		string(ReimbursementStatusPAID),
	}

	validMethodsOfReceipt := []string{
		string(MethodOfReceiptMILPAY),
		string(MethodOfReceiptOTHERDD),
		string(MethodOfReceiptGTCC),
	}

	return validate.Validate(
		&validators.IntIsGreaterThan{Field: int(r.RequestedAmount), Name: "RequestedAmount", Compared: 0},
		&validators.StringInclusion{Field: string(r.Status), Name: "Status", List: validStatuses},
		&validators.StringInclusion{Field: string(r.MethodOfReceipt), Name: "Status", List: validMethodsOfReceipt},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Reimbursement) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Reimbursement) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchReimbursement Fetches and Validates a Reimbursement model
func FetchReimbursement(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Reimbursement, error) {
	var reimbursement Reimbursement
	err := db.Q().Find(&reimbursement, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	return &reimbursement, nil
}
