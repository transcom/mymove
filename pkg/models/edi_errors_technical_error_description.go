package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// EdiErrorsTechnicalErrorDescription stores the reports Technical Error Descriptions (TEDs) recorded from an EDI 824
type EdiErrorsTechnicalErrorDescription struct {
	ID                         uuid.UUID                                `json:"id" db:"id"`
	CreatedAt                  time.Time                                `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time                                `json:"updated_at" db:"updated_at"`
	EdiErrorID                 uuid.UUID                                `json:"edi_error_id" db:"edi_error_id"`
	EdiError                   EdiError                                 `belongs_to:"edi_errors"`
	PaymentRequestID           uuid.UUID                                `json:"payment_request_id" db:"payment_request_id"`
	PaymentRequest             PaymentRequest                           `belongs_to:"payment_requests"`
	InterchangeControlNumberID uuid.UUID                                `json:"interchange_control_number_id" db:"interchange_control_number_id"`
	InterchangeControlNumber   PaymentRequestToInterchangeControlNumber `belongs_to:"payment_request_to_interchange_control_numbers"`
	Code                       string                                   `json:"code" db:"code"`
	Description                string                                   `json:"description" db:"description"`
	EDIType                    EDIType                                  `json:"edi_type" db:"edi_type"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EdiErrorsTechnicalErrorDescription) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(e.EDIType), Name: "EDIType", List: []string{
		string(EDI824),
	}})
	vs = append(vs, &validators.UUIDIsPresent{Field: e.EdiErrorID, Name: "EdiErrorID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: e.PaymentRequestID, Name: "PaymentRequestID"})
	vs = append(vs, &validators.UUIDIsPresent{Field: e.InterchangeControlNumberID, Name: "InterchangeControlNumberID"})
	if strings.TrimSpace(e.Code) == "" && strings.TrimSpace(e.Description) == "" {
		vs = append(vs, &validators.StringIsPresent{Field: e.Code, Name: "Code", Message: "Code or Description must be present"})
		vs = append(vs, &validators.StringIsPresent{Field: e.Description, Name: "Description", Message: "Code or Description must be present"})
	}
	return validate.Validate(vs...), nil
}

// EdiErrorsTechnicalErrorDescriptions is a list of EDI Technical Error Descriptions (TEDs)
type EdiErrorsTechnicalErrorDescriptions []EdiErrorsTechnicalErrorDescription
