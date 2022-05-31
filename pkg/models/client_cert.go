package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ClientCert represents a known x509 Certificate in the database.
type ClientCert struct {
	ID                          uuid.UUID `json:"id" db:"id"`
	Sha256Digest                string    `db:"sha256_digest"`
	Subject                     string    `db:"subject"`
	AllowDpsAuthAPI             bool      `db:"allow_dps_auth_api"`
	AllowOrdersAPI              bool      `db:"allow_orders_api"`
	CreatedAt                   time.Time `db:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at"`
	AllowAirForceOrdersRead     bool      `db:"allow_air_force_orders_read"`
	AllowAirForceOrdersWrite    bool      `db:"allow_air_force_orders_write"`
	AllowArmyOrdersRead         bool      `db:"allow_army_orders_read"`
	AllowArmyOrdersWrite        bool      `db:"allow_army_orders_write"`
	AllowCoastGuardOrdersRead   bool      `db:"allow_coast_guard_orders_read"`
	AllowCoastGuardOrdersWrite  bool      `db:"allow_coast_guard_orders_write"`
	AllowMarineCorpsOrdersRead  bool      `db:"allow_marine_corps_orders_read"`
	AllowMarineCorpsOrdersWrite bool      `db:"allow_marine_corps_orders_write"`
	AllowNavyOrdersRead         bool      `db:"allow_navy_orders_read"`
	AllowNavyOrdersWrite        bool      `db:"allow_navy_orders_write"`
	AllowPrime                  bool      `db:"allow_prime"`
	UserID                      uuid.UUID `db:"user_id"`
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
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	return &cert, nil
}

// GetAllowedOrdersIssuersRead returns a slice with the issuers of Orders that this ClientCert is allowed to read
func (c *ClientCert) GetAllowedOrdersIssuersRead() []string {
	var issuers []string
	if c.AllowAirForceOrdersRead {
		issuers = append(issuers, string(IssuerAirForce))
	}
	if c.AllowArmyOrdersRead {
		issuers = append(issuers, string(IssuerArmy))
	}
	if c.AllowCoastGuardOrdersRead {
		issuers = append(issuers, string(IssuerCoastGuard))
	}
	if c.AllowMarineCorpsOrdersRead {
		issuers = append(issuers, string(IssuerMarineCorps))
	}
	if c.AllowNavyOrdersRead {
		issuers = append(issuers, string(IssuerNavy))
	}
	return issuers
}
