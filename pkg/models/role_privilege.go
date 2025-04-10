package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models/roles"
)

// RolePrivilege represents a role->privilege mapping
type RolePrivilege struct {
	ID            uuid.UUID      `db:"id" rw:"r"`
	RoleID        uuid.UUID      `db:"role_id" rw:"r"`
	Role          roles.Role     `belongs_to:"roles" fk_id:"role_id" rw:"r"`
	PrivilegeID   uuid.UUID      `db:"privilege_id" rw:"r"`
	Privilege     Privilege      `belongs_to:"privileges" fk_id:"privilege_id" rw:"r"`
	CreatedAt     time.Time      `db:"created_at" rw:"r"`
	UpdatedAt     time.Time      `db:"updated_at" rw:"r"`
	RoleType      roles.RoleType `db:"role_type" rw:"r"`
	PrivilegeType PrivilegeType  `db:"privilege_type" rw:"r"`
	RoleName      roles.RoleName `db:"role_name" rw:"r"`
	PrivilegeName PrivilegeName  `db:"privilege_name" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (u RolePrivilege) TableName() string {
	return "roles_privileges"
}
