package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type PaymentServiceItemStatus string

func (p PaymentServiceItemStatus) String() string {
	return string(p)
}

const (
	PaymentServiceItemStatusRequested PaymentServiceItemStatus = "REQUESTED"
	PaymentServiceItemStatusApproved  PaymentServiceItemStatus = "APPROVED"
	PaymentServiceItemStatusDenied    PaymentServiceItemStatus = "DENIED"
	PaymentServiceItemStatusSentToGex PaymentServiceItemStatus = "SENT_TO_GEX"
	PaymentServiceItemStatusPaid      PaymentServiceItemStatus = "PAID"
)

var validPaymentServiceItemStatus = []string{
	string(PaymentServiceItemStatusRequested),
	string(PaymentServiceItemStatusApproved),
	string(PaymentServiceItemStatusDenied),
	string(PaymentServiceItemStatusSentToGex),
	string(PaymentServiceItemStatusPaid),
}

type PaymentServiceItem struct {
	ID               uuid.UUID                `json:"id" db:"id"`
	PaymentRequestID uuid.UUID                `json:"payment_request_id" db:"payment_request_id"`
	ServiceItemID    uuid.UUID                `json:"service_item_id" db:"service_item_id"`
	Status           PaymentServiceItemStatus `json:"status" db:"status"`
	PriceCents       unit.Cents               `json:"price_cents" db:"price_cents"`
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
	ServiceItem              MTOServiceItem           `belongs_to:"mto_service_item"`
	PaymentServiceItemParams PaymentServiceItemParams `has_many:"payment_service_item_params"`
}

// PaymentServiceItems is not required by pop and may be deleted
type PaymentServiceItems []PaymentServiceItem

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentServiceItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "PaymentRequestID"},
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "ServiceItemID"},
		&validators.StringInclusion{Field: p.Status.String(), Name: "Status", List: validPaymentServiceItemStatus},
		&validators.TimeIsPresent{Field: p.RequestedAt, Name: "RequestedAt"},
		// TODO: Removing this until we have pricing to populate
		// &validators.IntIsPresent{Field: p.PriceCents.Int(), Name: "PriceCents"},
		// &validators.IntIsGreaterThan{Field: p.PriceCents.Int(), Name: "PriceCents", Compared: 0},
	), nil
}
