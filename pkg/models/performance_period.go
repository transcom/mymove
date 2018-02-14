package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
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
func (p PerformancePeriod) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PerformancePeriods is not required by pop and may be deleted
type PerformancePeriods []PerformancePeriod

// String is not required by pop and may be deleted
func (p PerformancePeriods) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PerformancePeriod) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: p.StartDate, Name: "StartDate"},
		&validators.TimeIsPresent{Field: p.EndDate, Name: "EndDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PerformancePeriod) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PerformancePeriod) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
