package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// PaymentServiceItemStatus is a type of Payment Service Item Status
type PaymentServiceItemStatus string

// PaymentServiceItemStatus is a string representation of a Payment Service Item Status
func (p PaymentServiceItemStatus) String() string {
	return string(p)
}

const (
	// PaymentServiceItemStatusRequested is the requested status
	PaymentServiceItemStatusRequested PaymentServiceItemStatus = "REQUESTED"
	// PaymentServiceItemStatusApproved is the approved status
	PaymentServiceItemStatusApproved PaymentServiceItemStatus = "APPROVED"
	// PaymentServiceItemStatusDenied is the denied status
	PaymentServiceItemStatusDenied PaymentServiceItemStatus = "DENIED"
	// PaymentServiceItemStatusSentToGex is the sent-to-gex status
	PaymentServiceItemStatusSentToGex PaymentServiceItemStatus = "SENT_TO_GEX"
	// PaymentServiceItemStatusPaid is the paid status
	PaymentServiceItemStatusPaid PaymentServiceItemStatus = "PAID"
)

var validPaymentServiceItemStatus = []string{
	string(PaymentServiceItemStatusRequested),
	string(PaymentServiceItemStatusApproved),
	string(PaymentServiceItemStatusDenied),
	string(PaymentServiceItemStatusSentToGex),
	string(PaymentServiceItemStatusPaid),
}

// PaymentServiceItem represents a payment service item
type PaymentServiceItem struct {
	ID               uuid.UUID                `json:"id" db:"id"`
	PaymentRequestID uuid.UUID                `json:"payment_request_id" db:"payment_request_id"`
	MTOServiceItemID uuid.UUID                `json:"mto_service_item_id" db:"mto_service_item_id"`
	Status           PaymentServiceItemStatus `json:"status" db:"status"`
	PriceCents       *unit.Cents              `json:"price_cents" db:"price_cents"`
	RejectionReason  *string                  `json:"rejection_reason" db:"rejection_reason"`
	RequestedAt      time.Time                `json:"requested_at" db:"requested_at"`
	ApprovedAt       *time.Time               `json:"approved_at" db:"approved_at"`
	DeniedAt         *time.Time               `json:"denied_at" db:"denied_at"`
	SentToGexAt      *time.Time               `json:"sent_to_gex_at" db:"sent_to_gex_at"`
	PaidAt           *time.Time               `json:"paid_at" db:"paid_at"`
	CreatedAt        time.Time                `db:"created_at"`
	UpdatedAt        time.Time                `db:"updated_at"`

	//Associations
	PaymentRequest           PaymentRequest           `belongs_to:"payment_request"`
	MTOServiceItem           MTOServiceItem           `belongs_to:"mto_service_item"`
	PaymentServiceItemParams PaymentServiceItemParams `has_many:"payment_service_item_params"`
}

// PaymentServiceItems is not required by pop and may be deleted
type PaymentServiceItems []PaymentServiceItem

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentServiceItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "PaymentRequestID"},
		&validators.UUIDIsPresent{Field: p.MTOServiceItemID, Name: "MTOServiceItemID"},
		&validators.StringInclusion{Field: p.Status.String(), Name: "Status", List: validPaymentServiceItemStatus},
		&validators.TimeIsPresent{Field: p.RequestedAt, Name: "RequestedAt"},
	), nil
}
