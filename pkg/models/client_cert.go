package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ClientCert represents a known x509 Certificate in the database. It stores the SSN securely by hashing it.
type ClientCert struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Sha256Digest    string    `db:"sha256_digest"`
	Subject         string    `db:"subject"`
	AllowDpsAuthAPI bool      `db:"allow_dps_auth_api"`
	AllowOrdersAPI  bool      `db:"allow_orders_api"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *ClientCert) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Sha256Digest, Name: "Sha256Digest"},
		&validators.StringIsPresent{Field: c.Subject, Name: "Subject"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *ClientCert) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *ClientCert) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchClientCert fetches and validates a client certificate by digest
func FetchClientCert(db *pop.Connection, sha256Digest string) (*ClientCert, error) {
	var cert ClientCert
	err := db.Eager().Where("sha256_digest = $1", sha256Digest).First(&cert)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	return &cert, nil
}
