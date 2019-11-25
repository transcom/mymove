package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
)

// Contractor is an object representing an access code for a service member
type Contractor struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Name           string     `json:"code" db:"name"`
	Type           string     `json:"type" db:"type"`
	ContractNumber string     `json:"contract_number" db:"contract_number"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ClaimedAt      *time.Time `json:"claimed_at" db:"claimed_at"`
}

type Contractors []Contractor

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Contractor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(c.Name), Name: "Name"},
		&validators.StringIsPresent{Field: string(c.Type), Name: "Type"},
		&validators.StringIsPresent{Field: string(c.ContractNumber), Name: "ContractNumber"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Contractor) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Contractor) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
