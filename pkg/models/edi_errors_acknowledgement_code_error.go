package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// EdiErrorsAcknowledgementCodeError records function acknowledgment codes received from an EDI 997 for a payment request
type EdiErrorsAcknowledgementCodeError struct {
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

// EdiErrorsAcknowledgementCodeErrors list of EdiErrorsAcknowledgementCodeError
type EdiErrorsAcknowledgementCodeErrors []EdiErrorsAcknowledgementCodeError
