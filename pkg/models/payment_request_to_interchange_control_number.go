package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PaymentRequestToInterchangeControlNumber is an object that links payment requests to an Interchange Control Number used in the EDI 858 invoice
type PaymentRequestToInterchangeControlNumber struct {
	ID                       uuid.UUID `db:"id"`
	PaymentRequestID         uuid.UUID `db:"payment_request_id"`
	InterchangeControlNumber int       `db:"interchange_control_number"`
	EDIType                  EDIType   `db:"edi_type"`

	// Associations
	PaymentRequest PaymentRequest `belongs_to:"payment_requests" fk_id:"payment_request_id"`
}

// TableName overrides the table name used by Pop.
func (p PaymentRequestToInterchangeControlNumber) TableName() string {
	return "payment_request_to_interchange_control_numbers"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *PaymentRequestToInterchangeControlNumber) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringInclusion{Field: p.EDIType.String(), Name: "EDIType", List: allowedEDITypes},
		&validators.UUIDIsPresent{Field: p.PaymentRequestID, Name: "PaymentRequestID"},
		// minimum interchange control number must be greater than 0
		&validators.IntIsGreaterThan{Field: p.InterchangeControlNumber, Name: "InterchangeControlNumber", Compared: 0},
		// max interchange control number must be less than 1000000000
		&validators.IntIsLessThan{Field: p.InterchangeControlNumber, Name: "InterchangeControlNumber", Compared: 1000000000},
	), nil
}
