package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// EdiError stores errors found while sending an 858 and being reported from EDI response files (824, 997)
type EdiError struct {
	ID                         uuid.UUID                                `json:"id" db:"id"`
	CreatedAt                  time.Time                                `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                                `json:"updated_at" db:"updated_at"`
	PaymentRequestID           uuid.UUID                                `json:"payment_request_id" db:"payment_request_id"`
	PaymentRequest             PaymentRequest                           `belongs_to:"payment_requests"`
	InterchangeControlNumberID *uuid.UUID                               `json:"interchange_control_number_id" db:"interchange_control_number_id"`
	InterchangeControlNumber   PaymentRequestToInterchangeControlNumber `belongs_to:"payment_request_to_interchange_control_numbers"`
	Code                       *string                                  `json:"code" db:"code"`
	Description                *string                                  `json:"description" db:"description"`
	EDIType                    EDIType                                  `json:"edi_type" db:"edi_type"`
}

// EdiErrors is a list of EDI Error
type EdiErrors []EdiError

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EdiError) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(e.EDIType), Name: "EDIType", List: allowedEDITypes})
	vs = append(vs, &validators.UUIDIsPresent{Field: e.PaymentRequestID, Name: "PaymentRequestID"})
	if e.InterchangeControlNumberID != nil {
		vs = append(vs, &validators.UUIDIsPresent{Field: *e.InterchangeControlNumberID, Name: "InterchangeControlNumberID"})
	}
	vs = append(vs, &AtLeastOneNotNil{FieldName1: "Code", FieldValue1: e.Code, FieldName2: "Description", FieldValue2: e.Description})
	if e.Code != nil {
		vs = append(vs, &validators.StringIsPresent{Field: *e.Code, Name: "Code", Message: "Code string if present should not be empty"})
	}
	if e.Description != nil {
		vs = append(vs, &validators.StringIsPresent{Field: *e.Description, Name: "Description", Message: "Description string if present should not be empty"})
	}

	return validate.Validate(vs...), nil
}
