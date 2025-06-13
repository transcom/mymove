package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models/roles"
)

// PrivilegeAssociator is the service object interface for fetching privileges for a user id
//
//go:generate mockery --name PrivilegeAssociator
type PrivilegeAssociator interface {
	FetchPrivilegesForUser(appCtx appcontext.AppContext, userID uuid.UUID) (roles.Privileges, error)
	FetchPrivilegeTypes(appCtx appcontext.AppContext) ([]roles.PrivilegeType, error)
}
