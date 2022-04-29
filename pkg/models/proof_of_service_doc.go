package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ProofOfServiceDoc represents a document for proof of service
type ProofOfServiceDoc struct {
	ID               uuid.UUID `json:"id" db:"id"`
	PaymentRequestID uuid.UUID `json:"payment_request_id" db:"payment_request_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	//Associations
	PaymentRequest PaymentRequest `belongs_to:"payment_request" fk_id:"payment_request_id"`
	PrimeUploads   PrimeUploads   `has_many:"prime_uploads" fk_id:"proof_of_service_docs_id" order_by:"created_at asc"`
}

// ProofOfServiceDocs is not required by pop and may be deleted
type ProofOfServiceDocs []ProofOfServiceDoc

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ProofOfServiceDoc) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "PaymentRequestID"},
	), nil
}
