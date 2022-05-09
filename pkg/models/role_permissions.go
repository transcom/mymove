package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// RolePermissions represents a role and a permission. Roles grant permissions to users.
type RolePermissions struct {
	ID           uuid.UUID `db:"id"`
	PermissionID uuid.UUID `db:"permission_id"`
	RoleID       uuid.UUID `db:"role_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdateAt     time.Time `db:"updated_at"`
}
