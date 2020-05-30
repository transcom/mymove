package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"

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
	RevokeUserSession(id uuid.UUID, payload *adminmessages.UserRevokeSessionPayload, redisPool *redis.Pool) (*models.User, *validate.Errors, error)
}
