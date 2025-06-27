package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"
)

// RoleAssociater is the service object interface for fetching roles for a user id
//
//go:generate mockery --name RoleAssociater
type RoleAssociater interface {
	FetchRolesForUser(appCtx appcontext.AppContext, userID uuid.UUID) (roles.Roles, error)
	FetchRolesPrivileges(appCtx appcontext.AppContext) ([]roles.Role, error)
	FetchRoleTypes(appCtx appcontext.AppContext) ([]roles.RoleType, error)
}
