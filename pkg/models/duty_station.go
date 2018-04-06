package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// DutyStation represents a military duty station for a specific branch
type DutyStation struct {
	ID        uuid.UUID                       `json:"id" db:"id"`
	CreatedAt time.Time                       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time                       `json:"updated_at" db:"updated_at"`
	Name      string                          `json:"name" db:"name"`
	Branch    internalmessages.MilitaryBranch `json:"branch" db:"branch"`
	AddressID uuid.UUID                       `json:"address_id" db:"address_id"`
	Address   Address                         `belongs_to:"address"`
}

// String is not required by pop and may be deleted
func (d DutyStation) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// DutyStations is not required by pop and may be deleted
type DutyStations []DutyStation

// String is not required by pop and may be deleted
func (d DutyStations) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DutyStation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Name, Name: "Name"},
		&validators.UUIDIsPresent{Field: d.AddressID, Name: "AddressID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *DutyStation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *DutyStation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
