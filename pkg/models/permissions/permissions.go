package permissions

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// PermissionType represents the types of permissions a role can have
type PermissionType string

const (
	PermissionEditHistory PermissionType = "edit.shipment"
)

// Permission represents something a user can be allowed or restricted from doing/viewing
type Permission struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	PermissionType PermissionType `json:"permission_type" db:"permission_type"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

// Permissions is a slice of Permission objects
type Permissions []Permission

// HasPermission validates if Permission has a permission of a particular type
func (ps Permissions) HasPermission(permissionType PermissionType) bool {
	for _, p := range ps {
		if p.PermissionType == permissionType {
			return true
		}
	}
	return false
}

// GetPermission returns the permission for a permission type
func (ps Permissions) GetPermission(permissionType PermissionType) (Permission, bool) {
	for _, p := range ps {
		if p.PermissionType == permissionType {
			return p, true
		}
	}
	return Permission{}, false
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Permission) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ID, Name: "ID"},
		&validators.StringIsPresent{Field: string(r.PermissionType), Name: "PermissionType"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Permission) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Permission) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
