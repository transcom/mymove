package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"go.uber.org/zap/zapcore"

	"github.com/gofrs/uuid"
)

// EDIProcessing represents an email sent to a service member
type EDIProcessing struct {
	ID               uuid.UUID `db:"id"`
	EDIType          EDIType   `db:"edi_type"`
	NumEDIsProcessed int       `db:"num_edis_processed"`
	ProcessStartedAt time.Time `db:"process_started_at"`
	ProcessEndedAt   time.Time `db:"process_ended_at"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// EDIProcessings is a slice of notification structs
type EDIProcessings []EDIProcessing

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *EDIProcessing) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringInclusion{Field: string(e.EDIType), Name: "EDIType", List: allowedEDITypes},
		&validators.IntIsGreaterThan{Field: e.NumEDIsProcessed, Name: "NumEDIsProcessed", Compared: -1},
		&validators.TimeIsPresent{Field: e.ProcessStartedAt, Name: "ProcessStartedAt"},
		&validators.TimeIsPresent{Field: e.ProcessEndedAt, Name: "ProcessEndedAt"},
	), nil
}

// TableName overrides the table name used by Pop.
func (e *EDIProcessing) TableName() string {
	return "edi_processings"
}

// MarshalLogObject is required to be able to zap.Object log this model.
func (e *EDIProcessing) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("EDIType", e.EDIType.String())
	encoder.AddInt("NumEDIsProcessed", e.NumEDIsProcessed)
	encoder.AddTime("ProcessStartedAt", e.ProcessStartedAt)
	encoder.AddTime("ProcessEndedAt", e.ProcessEndedAt)
	return nil
}
