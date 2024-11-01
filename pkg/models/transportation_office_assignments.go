package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// TransportationAssignment is the transportation office the OfficeUser is assigned to
type TransportationAssignment struct {
	ID                     uuid.UUID            `json:"id" db:"id"`
	TransportationOfficeID uuid.UUID            `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice   TransportationOffice `belongs_to:"transportation_office" fk_id:"transportation_office_id"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
	PrimaryOffice          bool                 `json:"primary_office" db:"primary_office"`
}

// TableName overrides the table name used by Pop.
func (o TransportationAssignment) TableName() string {
	return "transportation_office_assignments"
}

type TransportationAssignments []TransportationAssignment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *TransportationAssignment) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
	), nil
}
