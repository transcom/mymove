package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// PrimeRequester is an object representing the Prime API Requester
type PrimeRequester struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	ClientCertID uuid.UUID `db:"client_cert_id"`
	AllowAccess  bool      `db:"allow_access"`
	LastSeenAt   time.Time `db:"last_seen_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	// Associations
	ClientCert ClientCert `belongs_to:"client_certs"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *PrimeRequester) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs,
		&validators.UUIDIsPresent{Field: m.ID, Name: "ID"},
		&validators.UUIDIsPresent{Field: m.ClientCertID, Name: "ClientCertID"},
		&validators.StringIsPresent{Field: m.Name, Name: "Name"})
	return validate.Validate(vs...), nil
}

// PrimeRequesters is a list of move task orders
type PrimeRequesters []PrimeRequester
