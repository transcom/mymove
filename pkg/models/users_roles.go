package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// UsersRoles represents a user and a role
type UsersRoles struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	RoleID    uuid.UUID  `db:"role_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdateAt  time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (u UsersRoles) TableName() string {
	return "users_roles"
}
