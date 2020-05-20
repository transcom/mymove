package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// SessionRevocation is the exported interface for revoking a user session
//go:generate mockery -name SessionRevocation
type SessionRevocation interface {
	RevokeUserSession(userID uuid.UUID, payload *adminmessages.RevokedSessionPayload, redisPool *redis.Pool) (*models.User, *validate.Errors, error)
}
