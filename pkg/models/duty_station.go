package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

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
		&BranchIsPresent{Field: d.Branch, Name: "Branch"},
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

// FindDutyStations returns all duty stations matching a search query and military branch
func FindDutyStations(tx *pop.Connection, search string, branch string) (DutyStations, error) {
	var stations DutyStations

	// ILIKE does case-insensitive pattern matching, "%" matches any string
	searchQuery := fmt.Sprintf("%%%s%%", search)
	query := tx.Eager().Where("branch = $1 AND name ILIKE $2", branch, searchQuery)

	if err := query.All(&stations); err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return stations, err
		}
	}

	return stations, nil
}

// FetchDutyStation returns all duty stations matching a search query and military branch
func FetchDutyStation(tx *pop.Connection, id uuid.UUID) (DutyStation, error) {
	var station DutyStation

	if err := tx.Eager().Find(&station, id); err != nil {
		return station, err
	}

	return station, nil
}
