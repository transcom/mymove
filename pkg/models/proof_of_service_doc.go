package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type ProofOfServiceDoc struct {
	ID               uuid.UUID `json:"id" db:"id"`
	PaymentRequestID uuid.UUID `json:"payment_request_id" db:"payment_request_id"`
	UploadID         uuid.UUID `json:"upload_id" db:"upload_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	//Associations
	PaymentRequest PaymentRequest `belongs_to:"payment_request"`
	Upload         PaymentRequest `belongs_to:"upload"`
}

// ProofOfServiceDocs is not required by pop and may be deleted
type ProofOfServiceDocs []ProofOfServiceDoc

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ProofOfServiceDoc) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "PaymentRequestID"},
		&validators.UUIDIsPresent{Field: p.UploadID, Name: "UploadID"},
	), nil
}
