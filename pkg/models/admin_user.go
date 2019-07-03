package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type AdminUser struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Role           string    `json:"role" db:"role"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
}

// String is not required by pop and may be deleted
func (a AdminUser) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AdminUsers is not required by pop and may be deleted
type AdminUsers []AdminUser

// String is not required by pop and may be deleted
func (a AdminUsers) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AdminUser) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validAdminRoles := []string{"SYSTEM_ADMIN", "PROGRAM_ADMIN"}
	return validate.Validate(
		&validators.StringIsPresent{Field: a.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: a.LastName, Name: "LastName"},
		&validators.UUIDIsPresent{Field: a.UserID, Name: "UserID"},
		&validators.StringInclusion{Field: a.Role, Name: "Role", List: validAdminRoles},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AdminUser) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AdminUser) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
