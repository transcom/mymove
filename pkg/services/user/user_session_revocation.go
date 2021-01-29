package user

import (
	"fmt"

	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userSessionRevocation struct {
	builder userQueryBuilder
}

// RevokeUserSession revokes the user's session
func (o *userSessionRevocation) RevokeUserSession(id uuid.UUID, payload *adminmessages.UserUpdatePayload, sessionStore scs.Store) (*models.User, *validate.Errors, error) {
	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	redisErr := deleteSessionIDFromRedis(foundUser, payload, sessionStore)
	if redisErr != nil {
		return nil, nil, redisErr
	}

	return deleteSessionIDFromDB(o, foundUser, payload)
}

// NewUserSessionRevocation returns a new admin user creator builder
func NewUserSessionRevocation(builder userQueryBuilder) services.UserSessionRevocation {
	return &userSessionRevocation{builder}
}

func deleteSessionIDFromRedis(user models.User, payload *adminmessages.UserUpdatePayload, sessionStore scs.Store) error {
	var currentAdminSessionID, currentOfficeSessionID, currentMilSessionID string
	userID := user.ID

	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession == true {
		currentAdminSessionID = user.CurrentAdminSessionID
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession == true {
		currentOfficeSessionID = user.CurrentOfficeSessionID
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession == true {
		currentMilSessionID = user.CurrentMilSessionID
	}

	var sessionIDMap = map[string]string{
		"adminSessionID":  currentAdminSessionID,
		"officeSessionID": currentOfficeSessionID,
		"milSessionID":    currentMilSessionID,
	}

	for field, sessionID := range sessionIDMap {
		_, exists, err := sessionStore.Find(sessionID)
		if err != nil {
			fmt.Printf("Error looking up %s in Redis for user ID %s", field, userID)
			return err
		}

		if !exists {
			fmt.Printf("%s not found in Redis; nothing to revoke. \n", field)
		} else {
			fmt.Printf("%s found for user ID %s; deleting it from Redis. \n", field, userID)
			err := sessionStore.Delete(sessionID)
			if err != nil {
				fmt.Printf("Error deleting %s from Redis for user ID %s.", field, userID)
				return err
			}
		}
	}

	return nil
}

func deleteSessionIDFromDB(o *userSessionRevocation, user models.User, payload *adminmessages.UserUpdatePayload) (*models.User, *validate.Errors, error) {
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
