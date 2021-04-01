package models

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
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
	// PaymentServiceItemStatusEDIError is an error status
	PaymentServiceItemStatusEDIError PaymentServiceItemStatus = "EDI_ERROR"

	// PaymentServiceItemMaxReferenceIDLength is the maximum overall length allowed for a reference ID
	// (given the EDI field's max length)
	PaymentServiceItemMaxReferenceIDLength = 30
	// PaymentServiceItemMinReferenceIDSuffixLength is the minimum suffix length for the PSI's reference ID
	PaymentServiceItemMinReferenceIDSuffixLength = 8
)

var validPaymentServiceItemStatus = []string{
	string(PaymentServiceItemStatusRequested),
	string(PaymentServiceItemStatusApproved),
	string(PaymentServiceItemStatusDenied),
	string(PaymentServiceItemStatusSentToGex),
	string(PaymentServiceItemStatusPaid),
	string(PaymentServiceItemStatusEDIError),
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
	ReferenceID      string                   `json:"reference_id" db:"reference_id"`
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
		// Note: ReferenceID is required, but can't be checked here because it's set in the BeforeCreate which
		// is called by Pop after Validate.
	), nil
}

// BeforeCreate is a Pop callback that is called before a PaymentServiceItem is created.
func (p *PaymentServiceItem) BeforeCreate(db *pop.Connection) error {
	// If we don't already have a UUID, create a new one so we can use it when assembling the reference ID.
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewV4())
	}

	// If we don't already have a reference ID, create a new one.
	if p.ReferenceID == "" {
		newReferenceID, err := p.GeneratePSIReferenceID(db)
		if err != nil {
			return err
		}

		p.ReferenceID = newReferenceID
	}

	return nil
}

// GeneratePSIReferenceID returns a reference ID for the PaymentServiceItem it is being called on.
// The format should be <MTO reference ID>-<part of PSI ID to make it unique>
func (p *PaymentServiceItem) GeneratePSIReferenceID(db *pop.Connection) (string, error) {
	var psiReferenceID string

	// Get the MTO's reference ID as that will serve as the prefix of the PSI's reference ID.
	var mto Move
	err := db.Q().
		InnerJoin("payment_requests pr", "moves.id = pr.move_id").
		Where("pr.id = $1", p.PaymentRequestID).
		First(&mto)
	if err != nil {
		return "", fmt.Errorf("could not find MTO for payment request ID %s: %w", p.PaymentRequestID, err)
	}

	if mto.ReferenceID == nil {
		return "", fmt.Errorf("nil reference ID for MTO %s: %w", mto.ID, err)
	}

	// Get a version of the PSI's ID without dashes.
	psiID := fmt.Sprintf("%x", p.ID)
	psiIDLength := len(psiID)

	// We want at least 8 hex digits of the ID, but no more than the max allowed by the EDI (if needed for uniqueness).
	currentLength := PaymentServiceItemMinReferenceIDSuffixLength
	maxSuffixLength := PaymentServiceItemMaxReferenceIDLength - len(*mto.ReferenceID) - 1 // including the hyphen
	if maxSuffixLength > psiIDLength {
		maxSuffixLength = psiIDLength
	}

	for currentLength <= maxSuffixLength {
		psiReferenceID = *mto.ReferenceID + "-" + psiID[:currentLength]

		// Check to see if it already exists
		count, err := db.Where("reference_id = $1", psiReferenceID).Count(&PaymentServiceItem{})
		if err != nil {
			return "", fmt.Errorf("could not count payment service item records for reference ID %s: %w", psiReferenceID, err)
		}

		// If we can't find another PSI with this reference ID, we're done.
		if count == 0 {
			return psiReferenceID, nil
		}

		// Add one to the length and try again
		currentLength++
	}

	// If we get here, we've exhausted all the possible hex digits of the UUID.
	return "", fmt.Errorf("cannot find unique PSI reference ID")
}
