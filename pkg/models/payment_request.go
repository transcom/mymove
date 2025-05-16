package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
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
	// PaymentRequestStatusReviewedAllRejected is reviewed
	PaymentRequestStatusReviewedAllRejected PaymentRequestStatus = "REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED"
	// PaymentRequestStatusSentToGex is sent to gex
	PaymentRequestStatusSentToGex PaymentRequestStatus = "SENT_TO_GEX"
	// PaymentRequestStatusTppsReceived is received by gex
	PaymentRequestStatusTppsReceived PaymentRequestStatus = "TPPS_RECEIVED"
	// PaymentRequestStatusPaid is paid
	PaymentRequestStatusPaid PaymentRequestStatus = "PAID"
	// PaymentRequestStatusEDIError an error has occurred
	PaymentRequestStatusEDIError PaymentRequestStatus = "EDI_ERROR"
	// PaymentRequestStatusDeprecated is deprecated
	PaymentRequestStatusDeprecated PaymentRequestStatus = "DEPRECATED"
)

const (
	// QueuePaymentRequestPaymentRequested status payment requested
	QueuePaymentRequestPaymentRequested = "Payment requested"
	// QueuePaymentRequestReviewed status Payment request reviewed
	QueuePaymentRequestReviewed = "Reviewed"
	// QueuePaymentRequestRejected status Payment request rejected
	QueuePaymentRequestRejected = "Rejected"
	// QueuePaymentRequestPaid status PaymentRequest paid
	QueuePaymentRequestPaid = "Paid"
	// QueuePaymentRequestDeprecated status PaymentRequest deprecated
	QueuePaymentRequestDeprecated = "Deprecated"
	// QueuePaymentRequestError status PaymentRequest error
	QueuePaymentRequestError = "Error"
)

var validPaymentRequestStatus = []string{
	string(PaymentRequestStatusPending),
	string(PaymentRequestStatusReviewed),
	string(PaymentRequestStatusReviewedAllRejected),
	string(PaymentRequestStatusSentToGex),
	string(PaymentRequestStatusTppsReceived),
	string(PaymentRequestStatusPaid),
	string(PaymentRequestStatusEDIError),
	string(PaymentRequestStatusDeprecated),
}

// PaymentRequest is an object representing a payment request on a move task order
type PaymentRequest struct {
	ID                              uuid.UUID            `json:"id" db:"id"`
	MoveTaskOrderID                 uuid.UUID            `db:"move_id"`
	IsFinal                         bool                 `json:"is_final" db:"is_final"`
	Status                          PaymentRequestStatus `json:"status" db:"status"`
	RejectionReason                 *string              `json:"rejection_reason" db:"rejection_reason"`
	PaymentRequestNumber            string               `json:"payment_request_number" db:"payment_request_number"`
	SequenceNumber                  int                  `json:"sequence_number" db:"sequence_number"`
	RequestedAt                     time.Time            `json:"requested_at" db:"requested_at"`
	ReviewedAt                      *time.Time           `json:"reviewed_at" db:"reviewed_at"`
	SentToGexAt                     *time.Time           `json:"sent_to_gex_at" db:"sent_to_gex_at"`
	ReceivedByGexAt                 *time.Time           `json:"received_by_gex_at" db:"received_by_gex_at"`
	PaidAt                          *time.Time           `json:"paid_at" db:"paid_at"`
	CreatedAt                       time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time            `json:"updated_at" db:"updated_at"`
	RecalculationOfPaymentRequestID *uuid.UUID           `json:"recalculation_of_payment_request_id" db:"recalculation_of_payment_request_id"`

	// Associations
	MoveTaskOrder                 Move                        `belongs_to:"moves" fk_id:"move_id"`
	PaymentServiceItems           PaymentServiceItems         `has_many:"payment_service_items" fk_id:"payment_request_id"`
	ProofOfServiceDocs            ProofOfServiceDocs          `has_many:"proof_of_service_docs" fk_id:"payment_request_id"`
	EdiErrors                     EdiErrors                   `has_many:"edi_errors" fk_id:"payment_request_id"`
	RecalculationOfPaymentRequest *PaymentRequest             `belongs_to:"payment_requests" fk_id:"recalculation_of_payment_request_id"`
	TPPSPaidInvoiceReports        TPPSPaidInvoiceReportEntrys `has_many:"tpps_paid_invoice_reports" fk_id:"payment_request_number"`
}

type PaymentRequestWithEarliestRequestedDate struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	EarliestRequestedDate time.Time `db:"requested_at"`
}

// TableName overrides the table name used by Pop.
func (p PaymentRequest) TableName() string {
	return "payment_requests"
}

// PaymentRequests is a slice of PaymentRequest
type PaymentRequests []PaymentRequest

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentRequest) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.MoveTaskOrderID, Name: "MoveTaskOrderID"},
		&validators.StringInclusion{Field: p.Status.String(), Name: "Status", List: validPaymentRequestStatus},
		&validators.StringIsPresent{Field: p.PaymentRequestNumber, Name: "PaymentRequestNumber"},
		&validators.IntIsGreaterThan{Field: p.SequenceNumber, Name: "SequenceNumber", Compared: 0},
		&OptionalUUIDIsPresent{Field: p.RecalculationOfPaymentRequestID, Name: "RecalculationOfPaymentRequestID"},
	), nil
}
