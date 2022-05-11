package roles

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// RoleType represents the types of users who can authenticate in the admin app
type RoleType string

// RoleName represents the names of roles
type RoleName string

const (
	// RoleTypeTOO is the Transportation Ordering Officer Role
	RoleTypeTOO RoleType = "transportation_ordering_officer"
	// RoleTypeCustomer is the Customer Role
	RoleTypeCustomer RoleType = "customer"
	// RoleTypeTIO is the Transportation Invoicing Officer Role
	RoleTypeTIO RoleType = "transportation_invoicing_officer"
	// RoleTypeContractingOfficer is the Contracting Officer Role
	RoleTypeContractingOfficer RoleType = "contracting_officer"
	// RoleTypePPMOfficeUsers is the PPM Office User Role
	RoleTypePPMOfficeUsers RoleType = "ppm_office_users"
	// RoleTypeServicesCounselor is the Services Counselor Role
	RoleTypeServicesCounselor RoleType = "services_counselor"
	// RoleTypePrimeSimulator is the PrimeSimulator Role
	RoleTypePrimeSimulator RoleType = "prime_simulator"
	// RoleTypeQaeCsr is the Quality Assurance and Customer Support Role
	RoleTypeQaeCsr RoleType = "qae_csr"
)

// Role represents a Role for users
type Role struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoleType  RoleType  `json:"role_type" db:"role_type"`
	RoleName  RoleName  `json:"role_name" db:"role_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Roles is a slice of Role objects
type Roles []Role

// HasRole validates if Role has a role of a particular type
func (rs Roles) HasRole(roleType RoleType) bool {
	for _, r := range rs {
		if r.RoleType == roleType {
			return true
		}
	}
	return false
}

// GetRole returns the role a Role Type
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
