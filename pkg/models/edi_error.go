package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// EdiError table is used to collect all errors for a payment requests
type EdiError struct {
	ID                                  uuid.UUID                           `json:"id" db:"id"`
	CreatedAt                           time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time                           `json:"updated_at" db:"updated_at"`
	PaymentRequestID                    uuid.UUID                           `json:"payment_request_id" db:"payment_request_id"`
	PaymentRequest                      PaymentRequest                      `belongs_to:"payment_requests"`
	EdiErrorsTechnicalErrorDescriptions EdiErrorsTechnicalErrorDescriptions `has_many:"edi_errors_technical_error_descriptions" fk_id:"edi_error_id"`
	EdiErrorsAcknowledgementCodeErrors  EdiErrorsAcknowledgementCodeErrors  `has_many:"edi_errors_acknowledgement_code_errors" fk_id:"edi_error_id"`
	EdiErrorsSendToSyncadaErrors        EdiErrorsSendToSyncadaErrors        `has_many:"edi_errors_send_to_syncada_errors" fk_id:"edi_error_id"`
}
