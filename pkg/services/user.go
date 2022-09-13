package services

import (
	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// UserFetcher is the service object interface for FetchUser
//
//go:generate mockery --name UserFetcher --disable-version-string
type UserFetcher interface {
	FetchUser(appCtx appcontext.AppContext, filters []QueryFilter) (models.User, error)
}

// UserUpdater is the service object interface for UpdateUser
//
//go:generate mockery --name UserUpdater --disable-version-string
type UserUpdater interface {
	UpdateUser(appCtx appcontext.AppContext, id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error)
}

// UserSessionRevocation is the exported interface for revoking a user session
//
//go:generate mockery --name UserSessionRevocation --disable-version-string
type UserSessionRevocation interface {
	RevokeUserSession(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.UserUpdatePayload, sessionManagers [3]*scs.SessionManager) (*models.User, *validate.Errors, error)
}
