package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/gofrs/uuid"
)

// EDIResponseMessageTypes represents types of EDI Responses
type EDIResponseMessageTypes string

const (
	// EDI997 captures enum value "997"
	EDI997 EDIResponseMessageTypes = "997"
	// EDI824 captures enum value "824"
	EDI824 EDIResponseMessageTypes = "824"
	// EDI810 captures enum value "810"
	EDI810 EDIResponseMessageTypes = "810"
)

// EDIResponseProcessing represents an email sent to a service member
type EDIResponseProcessing struct {
	ID               uuid.UUID               `db:"id"`
	MessageType      EDIResponseMessageTypes `db:"edi_response_message_type"`
	ProcessStartedAt time.Time               `db:"process_started_at"`
	ProcessEndedAt   time.Time               `db:"process_ended_at"`
}

// EDIResponseProcessings is a slice of notification structs
type EDIResponseProcessings []EDIResponseProcessing

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EDIResponseProcessing) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringInclusion{Field: string(e.MessageType), Name: "MessageType", List: []string{
			string(EDI997),
			string(EDI824),
			string(EDI810),
		}},
		&validators.TimeIsPresent{Field: e.ProcessStartedAt, Name: "ProcessStartedAt"},
		&validators.TimeIsPresent{Field: e.ProcessEndedAt, Name: "ProcessEndedAt"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *EDIResponseProcessing) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *EDIResponseProcessing) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
