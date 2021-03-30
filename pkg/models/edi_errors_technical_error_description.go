package models

import (
	"time"

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
	Source                     EDIType                                  `json:"source" db:"source"`
}

// EdiErrorsTechnicalErrorDescriptions is a list of EDI Technical Error Descriptions (TEDs)
type EdiErrorsTechnicalErrorDescriptions []EdiErrorsTechnicalErrorDescription
