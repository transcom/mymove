package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/gofrs/uuid"
)

// EDIType represents types of EDI Responses
type EDIType string

const (
	// EDIType810 captures enum value "810"
	EDIType810 EDIType = "810"
	// EDIType824 captures enum value "824"
	EDIType824 EDIType = "824"
	// EDIType858 captures enum value "858"
	EDIType858 EDIType = "858"
	// EDIType997 captures enum value "997"
	EDIType997 EDIType = "997"
)

// EDIProcessing represents an email sent to a service member
type EDIProcessing struct {
	ID               uuid.UUID `db:"id"`
	EDIType          EDIType   `db:"edi_type"`
	ProcessStartedAt time.Time `db:"process_started_at"`
	ProcessEndedAt   time.Time `db:"process_ended_at"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	NumEDIsProcessed int       `db:"num_edis_processed"`
}

// EDIProcessings is a slice of notification structs
type EDIProcessings []EDIProcessing

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EDIProcessing) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: e.ProcessStartedAt, Name: "ProcessStartedAt"},
		&validators.TimeIsPresent{Field: e.ProcessEndedAt, Name: "ProcessEndedAt"},
		&validators.IntIsPresent{Field: e.NumEDIsProcessed, Name: "NumEDIsProcessed"},
		&validators.StringInclusion{Field: string(e.EDIType), Name: "EDIType", List: []string{
			string(EDIType810),
			string(EDIType824),
			string(EDIType858),
			string(EDIType997),
		}},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *EDIProcessing) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *EDIProcessing) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
