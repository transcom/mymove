package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// UserFetcher is the service object interface for FetchUser
//
//go:generate mockery --name UserFetcher
type UserFetcher interface {
	FetchUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.User, error)
}

// UserUpdater is the service object interface for UpdateUser
//
//go:generate mockery --name UserUpdater
type UserUpdater interface {
	UpdateUser(appCtx appcontext.AppContext, id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error)
}

// UserSessionRevocation is the exported interface for revoking a user session
//
//go:generate mockery --name UserSessionRevocation
type UserSessionRevocation interface {
	RevokeUserSession(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.UserUpdate, sessionManagers auth.AppSessionManagers) (*models.User, *validate.Errors, error)
}

// UserDeleter is the exported interface for hard deleting a user and its associations (roles, privileges)
//
//go:generate mockery --name UserDeleter
type UserDeleter interface {
	DeleteUser(appCtx appcontext.AppContext, id uuid.UUID) error
}
