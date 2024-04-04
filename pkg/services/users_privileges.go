package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// UserPrivilegeAssociator is the service object interface for UpdateUserPrivileges
//
//go:generate mockery --name UserPrivilegeAssociator
type UserPrivilegeAssociator interface {
	UpdateUserPrivileges(appCtx appcontext.AppContext, userID uuid.UUID, privileges []models.PrivilegeType) ([]models.UsersPrivileges, error)
}
