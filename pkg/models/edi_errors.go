package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"go.uber.org/zap/zapcore"
)

// EdiError stores errors found while sending an 858 and being reported from EDI response files (824, 997)
type EdiError struct {
	ID                         uuid.UUID                                `json:"id" db:"id"`
	CreatedAt                  time.Time                                `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                                `json:"updated_at" db:"updated_at"`
	PaymentRequestID           uuid.UUID                                `json:"payment_request_id" db:"payment_request_id"`
	PaymentRequest             PaymentRequest                           `belongs_to:"payment_requests" fk_id:"payment_request_id"`
	InterchangeControlNumberID *uuid.UUID                               `json:"interchange_control_number_id" db:"interchange_control_number_id"`
	InterchangeControlNumber   PaymentRequestToInterchangeControlNumber `belongs_to:"payment_request_to_interchange_control_numbers" fk_id:"interchange_control_number_id"`
	Code                       *string                                  `json:"code" db:"code"`
	Description                *string                                  `json:"description" db:"description"`
	EDIType                    EDIType                                  `json:"edi_type" db:"edi_type"`
}

// TableName overrides the table name used by Pop.
func (e EdiError) TableName() string {
	return "edi_errors"
}

// EdiErrors is a list of EDI Error
type EdiErrors []EdiError

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EdiError) Validate(_ *pop.Connection) (*validate.Errors, error) {
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

// MarshalLogObject is required to be able to zap.Object log this model.
func (e *EdiError) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("EDIType", e.EDIType.String())
	encoder.AddString("PaymentRequestID", e.PaymentRequestID.String())
	encoder.AddString("Code", *e.Code)
	encoder.AddString("Description ", *e.Description)

	return nil
}

func FetchEdiErrorByPaymentRequestID(db *pop.Connection, paymentRequestID uuid.UUID) (EdiError, error) {
	var ediError EdiError
	err := db.Where("payment_request_id = $1", paymentRequestID).First(&ediError)
	if err != nil {
		return EdiError{}, err
	}
	return ediError, nil
}
