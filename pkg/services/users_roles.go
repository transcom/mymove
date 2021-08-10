package services

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gofrs/uuid"
)

// UserRoleAssociator is the service object interface for UpdateUserRoles
//go:generate mockery --name UserRoleAssociator --disable-version-string
type UserRoleAssociator interface {
	UpdateUserRoles(appCfg appconfig.AppConfig, userID uuid.UUID, roles []roles.RoleType) ([]models.UsersRoles, error)
}
