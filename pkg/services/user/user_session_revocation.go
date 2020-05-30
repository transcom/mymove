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

type userSessionRevocation struct {
	builder userQueryBuilder
}

// RevokeUserSession revokes the user's session
func (o *userSessionRevocation) RevokeUserSession(id uuid.UUID, payload *adminmessages.UserRevokeSessionPayload, redisPool *redis.Pool) (*models.User, *validate.Errors, error) {
	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	var currentSessionID string

	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession == true {
		currentSessionID = foundUser.CurrentAdminSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession == true {
		currentSessionID = foundUser.CurrentOfficeSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession == true {
		currentSessionID = foundUser.CurrentMilSessionID
		deleteSessionIDFromRedis(currentSessionID, redisPool)
	}

	return deleteSessionIDFromDB(o, foundUser, payload)
}

// NewUserSessionRevocation returns a new admin user creator builder
func NewUserSessionRevocation(builder userQueryBuilder) services.UserSessionRevocation {
	return &userSessionRevocation{builder}
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

func deleteSessionIDFromDB(o *userSessionRevocation, user models.User, payload *adminmessages.UserRevokeSessionPayload) (*models.User, *validate.Errors, error) {
	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession == true {
		user.CurrentAdminSessionID = ""
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession == true {
		user.CurrentOfficeSessionID = ""
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession == true {
		user.CurrentMilSessionID = ""
	}

	verrs, err := o.builder.UpdateOne(&user, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &user, nil, nil
}

// func foundUser(userID uuid.UUID, o *sessionRevocation) (*models.User, error) {
// 	var foundUser models.User
// 	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", userID.String())}
// 	err := o.builder.FetchOne(&foundUser, filters)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &foundUser, nil
// }
