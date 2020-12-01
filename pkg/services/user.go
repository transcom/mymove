package services

import (
	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// UserFetcher is the service object interface for FetchUser
//go:generate mockery -name UserFetcher
type UserFetcher interface {
	FetchUser(filters []QueryFilter) (models.User, error)
}

// UserSessionRevocation is the exported interface for revoking a user session
//go:generate mockery -name UserSessionRevocation
type UserSessionRevocation interface {
	RevokeUserSession(id uuid.UUID, payload *adminmessages.UserRevokeSessionPayload, sessionStore scs.Store) (*models.User, *validate.Errors, error)
}
