package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// UsersPrivileges represents a user and a privilege
type UsersPrivileges struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	PrivilegeID uuid.UUID  `db:"privilege_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (u UsersPrivileges) TableName() string {
	return "users_privileges"
}
