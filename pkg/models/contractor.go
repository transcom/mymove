package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Contractor is an object representing an access code for a service member
type Contractor struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"code" db:"name"`
	Type           string    `json:"type" db:"type"`
	ContractNumber string    `json:"contract_number" db:"contract_number"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Contractors is a slice of Contractor objects
type Contractors []Contractor

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (c *Contractor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: string(c.Name), Name: "Name"},
		&validators.StringIsPresent{Field: string(c.Type), Name: "Type"},
		&validators.StringIsPresent{Field: string(c.ContractNumber), Name: "ContractNumber"},
	), nil
}

// FetchGHCPrimeTestContractor returns a test contractor for dev
func FetchGHCPrimeTestContractor(db *pop.Connection) (*Contractor, error) {
	var contractor Contractor
	err := db.Q().Where("contract_number='HTC111-11-1-1111'").First(&contractor)
	if err != nil {
		err = db.Q().Where(`contract_number = $1`, "TEST").First(&contractor)
		if err != nil {
			if errors.Cause(err).Error() == RecordNotFoundErrorString {
				return nil, errors.Wrap(ErrFetchNotFound, "error fetching contractor")
			}
		}
		// Otherwise, it's an unexpected err so we return that.
		return &contractor, err
	}
	return &contractor, nil
}
