package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

// PerformancePeriod defines periods of time across the year - per James,
// there are five periods across the year of unequal lengths. We have gotten
// some conflicting information about this and should double-check before calling
// it done.
type PerformancePeriod struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}

// String is not required by pop and may be deleted
func (a PerformancePeriod) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// PerformancePeriods is not required by pop and may be deleted
type PerformancePeriods []PerformancePeriod

// String is not required by pop and may be deleted
func (a PerformancePeriods) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method. Todo: more fields to validate?
// This method is not required and may be deleted.
func (a *PerformancePeriod) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&v.UUIDIsPresent{Field: a.ID, Name: "ID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *PerformancePeriod) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *PerformancePeriod) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
