package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PaymentRequestEDI is an object that links payment requests to an Interchange Control Number used in the EDI 858 invoice
type PaymentRequestEDI struct {
	ID                       uuid.UUID `db:"id"`
	PaymentRequestID         uuid.UUID `db:"payment_request_id"`
	InterchangeControlNumber int       `db:"interchange_control_number"`
	EDIType                  EDIType   `db:"edi_type"`
	EDIText                  string    `db:"edi_text"`

	// Associations
	PaymentRequest PaymentRequest `belongs_to:"payment_requests" fk_id:"payment_request_id"`

	// POP managed fields
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (m PaymentRequestEDI) TableName() string {
	return "payment_request_edis"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *PaymentRequestEDI) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringInclusion{Field: m.EDIType.String(), Name: "EDIType", List: allowedEDITypes},
		&validators.UUIDIsPresent{Field: m.PaymentRequestID, Name: "PaymentRequestID"},
		// minimum interchange control number must be greater than 0
		&validators.IntIsGreaterThan{Field: m.InterchangeControlNumber, Name: "InterchangeControlNumber", Compared: 0},
		// max interchange control number must be less than 1000000000
		&validators.IntIsLessThan{Field: m.InterchangeControlNumber, Name: "InterchangeControlNumber", Compared: 1000000000},
		// edi text is present
		&validators.StringIsPresent{Field: m.EDIText, Name: "EDIText"},
	), nil
}
