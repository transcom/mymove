package services

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gofrs/uuid"
)

// UserRoleAssociator is the service object interface for AssociateUserRoles
//go:generate mockery -name UserRoleAssociator
type UserRoleAssociator interface {
	AssociateUserRoles(userID uuid.UUID, roles roles.Roles) ([]models.UsersRoles, error)
}
