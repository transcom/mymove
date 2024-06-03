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
	// RoleTypeServicesCounselor is the Services Counselor Role
	RoleTypeServicesCounselor RoleType = "services_counselor"
	// RoleTypePrimeSimulator is the PrimeSimulator Role
	RoleTypePrimeSimulator RoleType = "prime_simulator"
	// RoleTypeQae is the Quality Assurance Evaluator Role
	RoleTypeQae RoleType = "qae"
	// RoleTypeCustomerServiceRepresentative is the Customer Support Representative Role
	RoleTypeCustomerServiceRepresentative RoleType = "customer_service_representative"
	// RoleTypePrime is the Role associated with actions performed by the Prime
	RoleTypePrime RoleType = "prime"
	// RoleTypeHQ is the Headquarters Role
	RoleTypeHQ RoleType = "headquarters"
)

// Role represents a Role for users
type Role struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoleType  RoleType  `json:"role_type" db:"role_type"`
	RoleName  RoleName  `json:"role_name" db:"role_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (r Role) TableName() string {
	return "roles"
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
func (r *Role) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ID, Name: "ID"},
		&validators.StringIsPresent{Field: string(r.RoleType), Name: "RoleType"},
	), nil
}

// FetchRolesForUser gets the active RoleTypes for the user
func FetchRolesForUser(db *pop.Connection, userID uuid.UUID) (Roles, error) {
	var roles Roles

	err := db.Q().Join("users_roles", "users_roles.role_id = roles.id").
		Where("users_roles.deleted_at IS NULL AND users_roles.user_id = ?", (userID)).
		All(&roles)
	return roles, err
}
