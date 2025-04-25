package roles

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PrivilegeType represents the types of privileges available in the admin app
type PrivilegeType string

// PrivilegeName represents the names of privileges
type PrivilegeName string

const (
	// PrivilegeTypeSupervisor is the Task Ordering Officer Role
	PrivilegeTypeSupervisor PrivilegeType = "supervisor"
	PrivilegeTypeSafety     PrivilegeType = "safety"
)

const (
	// PrivilegeTypeSupervisor is the Task Ordering Officer Role
	PrivilegeSearchTypeSupervisor PrivilegeType = "SUPERVISOR"
	PrivilegeSearchTypeSafety     PrivilegeType = "SAFETY"
)

// Privilege represents a Privilege for users
type Privilege struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	PrivilegeType PrivilegeType `json:"privilege_type" db:"privilege_type"`
	PrivilegeName PrivilegeName `json:"privilege_name" db:"privilege_name"`
	Sort          int32         `json:"sort" db:"sort"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (r Privilege) TableName() string {
	return "privileges"
}

// Privileges is a slice of Privilege objects
type Privileges []Privilege

// HasPrivilege validates if Privilege has a privilege of a particular type
func (rs Privileges) HasPrivilege(privilegeType PrivilegeType) bool {
	for _, r := range rs {
		if r.PrivilegeType == privilegeType {
			return true
		}
	}
	return false
}

// GetPrivilege returns the privilege a Privilege Type
func (rs Privileges) GetPrivilege(privilegeType PrivilegeType) (Privilege, bool) {
	for _, r := range rs {
		if r.PrivilegeType == privilegeType {
			return r, true
		}
	}
	return Privilege{}, false
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *Privilege) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ID, Name: "ID"},
		&validators.StringIsPresent{Field: string(r.PrivilegeType), Name: "PrivilegeType"},
	), nil
}

// FetchPrivilegesForUser gets the active PrivilegeTypes for the user
func FetchPrivilegesForUser(db *pop.Connection, userID uuid.UUID) (Privileges, error) {
	var privileges Privileges

	err := db.Q().Join("users_privileges", "users_privileges.privilege_id = privileges.id").
		Where("users_privileges.deleted_at IS NULL AND users_privileges.user_id = ?", (userID)).
		All(&privileges)
	return privileges, err
}
