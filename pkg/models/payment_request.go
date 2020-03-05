package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// PaymentRequestStatus is a type of Payment Request Status
type PaymentRequestStatus string

// String is a string representation of a Payment Request Status
func (p PaymentRequestStatus) String() string {
	return string(p)
}

const (
	// PaymentRequestStatusPending is pending
	PaymentRequestStatusPending PaymentRequestStatus = "PENDING"
	// PaymentRequestStatusReviewed is reviewed
	PaymentRequestStatusReviewed PaymentRequestStatus = "REVIEWED"
	// PaymentRequestStatusSentToGex is sent to gex
	PaymentRequestStatusSentToGex PaymentRequestStatus = "SENT_TO_GEX"
	// PaymentRequestStatusReceivedByGex is received by gex
	PaymentRequestStatusReceivedByGex PaymentRequestStatus = "RECEIVED_BY_GEX"
	// PaymentRequestStatusPaid is paid
	PaymentRequestStatusPaid PaymentRequestStatus = "PAID"
)

var validPaymentRequestStatus = []string{
	string(PaymentRequestStatusPending),
	string(PaymentRequestStatusReviewed),
	string(PaymentRequestStatusSentToGex),
	string(PaymentRequestStatusReceivedByGex),
	string(PaymentRequestStatusPaid),
}

// PaymentRequest is an object representing a payment request on a move task order
type PaymentRequest struct {
	ID                   uuid.UUID            `json:"id" db:"id"`
	MoveTaskOrderID      uuid.UUID            `db:"move_task_order_id"`
	IsFinal              bool                 `json:"is_final" db:"is_final"`
	Status               PaymentRequestStatus `json:"status" db:"status"`
	RejectionReason      *string              `json:"rejection_reason" db:"rejection_reason"`
	PaymentRequestNumber string               `json:"payment_request_number" db:"payment_request_number"`
	SequenceNumber       int                  `json:"sequence_number" db:"sequence_number"`
	RequestedAt          time.Time            `json:"requested_at" db:"requested_at"`
	ReviewedAt           *time.Time           `json:"reviewed_at" db:"reviewed_at"`
	SentToGexAt          *time.Time           `json:"sent_to_gex_at" db:"sent_to_gex_at"`
	ReceivedByGexAt      *time.Time           `json:"received_by_gex_at" db:"received_by_gex_at"`
	PaidAt               *time.Time           `json:"paid_at" db:"paid_at"`
	CreatedAt            time.Time            `db:"created_at"`
	UpdatedAt            time.Time            `db:"updated_at"`

	// Associations
	MoveTaskOrder       MoveTaskOrder       `belongs_to:"move_task_orders"`
	PaymentServiceItems PaymentServiceItems `has_many:"payment_service_items"`
	ProofOfServiceDocs  ProofOfServiceDocs  `has_many:"proof_of_service_docs"`
}

// PaymentRequests is a slice of PaymentRequest
type PaymentRequests []PaymentRequest

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentRequest) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.MoveTaskOrderID, Name: "MoveTaskOrderID"},
		&validators.StringInclusion{Field: p.Status.String(), Name: "Status", List: validPaymentRequestStatus},
		&validators.StringIsPresent{Field: p.PaymentRequestNumber, Name: "PaymentRequestNumber"},
		&validators.IntIsGreaterThan{Field: p.SequenceNumber, Name: "SequenceNumber", Compared: 0},
	), nil
}
