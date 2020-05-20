package user

import (
	"fmt"

	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type sessionRevocation struct {
	builder userQueryBuilder
}

// RevokeUserSession revoke's the user's session
func (o *sessionRevocation) RevokeUserSession(userID uuid.UUID, payload *adminmessages.RevokedSessionPayload, redisPool *redis.Pool) (*models.User, *validate.Errors, error) {
	user, err := foundUser(userID, o)
	if err != nil {
		return nil, nil, err
	}

	var currentSessionID string

	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession == true {
		currentSessionID = user.CurrentAdminSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession == true {
		currentSessionID = user.CurrentOfficeSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession == true {
		currentSessionID = user.CurrentMilSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	return deleteSessionIDFromDB(o, userID, payload)
}

// NewSessionRevocation returns a new admin user creator builder
func NewSessionRevocation(builder userQueryBuilder) services.SessionRevocation {
	return &sessionRevocation{builder}
}

func deleteSessionIDFromRedis(currentSessionID string, redisPool *redis.Pool) error {
	conn := redisPool.Get()
	defer conn.Close()

	redisKey := fmt.Sprintf("scs:session:%s", currentSessionID)
	_, redisErr := redis.Bytes(conn.Do("GET", redisKey))
	if redisErr == redis.ErrNil {
		fmt.Println("session token not found, nothing to do")
		return nil
	} else if redisErr != nil {
		return redisErr
	}

	_, err := conn.Do("DEL", redisKey)
	if err != nil {
		return err
	}

	return nil
}

func deleteSessionIDFromDB(o *sessionRevocation, userID uuid.UUID, payload *adminmessages.RevokedSessionPayload) (*models.User, *validate.Errors, error) {
	user, err := foundUser(userID, o)
	if err != nil {
		return nil, nil, err
	}

	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession == true {
		user.CurrentAdminSessionID = ""
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession == true {
		user.CurrentOfficeSessionID = ""
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession == true {
		user.CurrentMilSessionID = ""
	}

	verrs, err := o.builder.UpdateOne(user, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return user, nil, nil
}

func foundUser(userID uuid.UUID, o *sessionRevocation) (*models.User, error) {
	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", userID.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}
