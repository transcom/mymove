package roles

import (
	"time"

	"github.com/gofrs/uuid"
)

// UsersRoles represents a user and a role
type RolePrivilege struct {
	ID          uuid.UUID `db:"id"`
	RoleID      uuid.UUID `db:"role_id"`
	PrivilegeID uuid.UUID `db:"privilege_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (u RolePrivilege) TableName() string {
	return "roles_privileges"
}
