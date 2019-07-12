package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type AdminRole string

const (
	SystemAdminRole  AdminRole = "SYSTEM_ADMIN"
	ProgramAdminRole AdminRole = "PROGRAM_ADMIN"
)

func (ar *AdminRole) ValidRoles() []AdminRole {
	return []AdminRole{
		SystemAdminRole,
		ProgramAdminRole,
	}
}

func (ar *AdminRole) String() string {
	return string(*ar)
}

type AdminUser struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
	UserID         *uuid.UUID   `json:"user_id" db:"user_id"`
	User           User         `belongs_to:"user"`
	Role           AdminRole    `json:"role" db:"role"`
	Email          string       `json:"email" db:"email"`
	FirstName      string       `json:"first_name" db:"first_name"`
	LastName       string       `json:"last_name" db:"last_name"`
	OrganizationID *uuid.UUID   `json:"organization_id" db:"organization_id"`
	Organization   Organization `belongs_to:"organization"`
	Disabled       bool         `json:"disabled" db:"disabled"`
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

type RoleInclusion struct {
	Name    string
	Field   AdminRole
	List    []AdminRole
	Message string
}

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
