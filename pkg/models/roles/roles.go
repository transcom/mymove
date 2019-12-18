package roles

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Role is an object representing the types of users who can authenticate in the admin app
type RoleType string

const (
	TOO      RoleType = "transportation_ordering_officer"
	Customer RoleType = "customer"
)

type Role struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoleType  RoleType  `json:"role_type" db:"role_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Roles []Role

func (rs Roles) HasRole(roleType RoleType) bool {
	for _, r := range rs {
		if r.RoleType == roleType {
			return true
		}
	}
	return false
}

func (rs Roles) GetRole(roleType RoleType) (Role, bool) {
	for _, r := range rs {
		if r.RoleType == roleType {
			return r, true
		}
	}
	return Role{}, false
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Role) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ID, Name: "ID"},
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
