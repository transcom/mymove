package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// RoleType defines a role type for a user
type RoleType string

const (
	RoleTypeTio                RoleType = "tio"
	RoleTypeToo                RoleType = "too"
	RoleTypeContractingOfficer RoleType = "contractingOfficer"
	RoleTypeOffice             RoleType = "office"
	RoleTypeCustomer           RoleType = "customer"
)

// Role is an object representing the types of users who can authenticate in the admin app
type Role struct {
	ID        int       `json:"id" db:"id"`
	RoleType  RoleType  `json:"role_type" db:"role_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Roles []Role

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Role) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: int(r.ID), Name: "ID"},
		&validators.StringIsPresent{Field: string(r.RoleType), Name: "RoleType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Role) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Role) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
