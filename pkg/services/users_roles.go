package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// UserRoleAssociator is the service object interface for UpdateUserRoles
//
//go:generate mockery --name UserRoleAssociator
type UserRoleAssociator interface {
	UpdateUserRoles(appCtx appcontext.AppContext, userID uuid.UUID, roles []roles.RoleType) ([]models.UsersRoles, *validate.Errors, error)
}
