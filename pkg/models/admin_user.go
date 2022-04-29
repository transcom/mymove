package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// AdminRole represents administrative roles
type AdminRole string

const (
	// SystemAdminRole represents a role for managing the system
	SystemAdminRole AdminRole = "SYSTEM_ADMIN"
	// ProgramAdminRole represents a role for managing the program (Note: This is deprecated and should be removed)
	ProgramAdminRole AdminRole = "PROGRAM_ADMIN"
)

// ValidRoles returns a slice of valid roles for an admin
func (ar *AdminRole) ValidRoles() []AdminRole {
	return []AdminRole{
		SystemAdminRole,
		ProgramAdminRole,
	}
}

// String returns a string representation of the admin role
func (ar *AdminRole) String() string {
	return string(*ar)
}

// AdminUser is someone who operates the Milmove systems
type AdminUser struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	UserID         *uuid.UUID   `json:"user_id" db:"user_id"`
	User           User         `belongs_to:"user" fk_id:"user_id"`
	Role           AdminRole    `json:"role" db:"role"`
	Email          string       `json:"email" db:"email"`
	FirstName      string       `json:"first_name" db:"first_name"`
	LastName       string       `json:"last_name" db:"last_name"`
	OrganizationID *uuid.UUID   `json:"organization_id" db:"organization_id"`
	Organization   Organization `belongs_to:"organization" fk_id:"organization_id"`
	Active         bool         `json:"active" db:"active"`
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

// RoleInclusion is used to validate if a role is valid for inclusion
type RoleInclusion struct {
	Name    string
	Field   AdminRole
	List    []AdminRole
	Message string
}

// IsValid validates if RoleInclusion is valid
func (v *RoleInclusion) IsValid(errors *validate.Errors) {
	found := false
	for _, l := range v.List {
		if l == v.Field {
			found = true
			break
		}
	}
	if !found {
		if len(v.Message) > 0 {
			errors.Add(validators.GenerateKey(v.Name), v.Message)
			return
		}

		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s is not in the list %+v.", v.Name, v.List))
	}
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AdminUser) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: a.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: a.Email, Name: "Email"},
		&RoleInclusion{Field: a.Role, Name: "Role", List: new(AdminRole).ValidRoles()},
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
