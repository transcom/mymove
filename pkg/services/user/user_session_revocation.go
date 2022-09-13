package user

import (
	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userSessionRevocation struct {
	builder userQueryBuilder
}

// RevokeUserSession revokes the user's session
func (o *userSessionRevocation) RevokeUserSession(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.UserUpdatePayload, sessionManagers [3]*scs.SessionManager) (*models.User, *validate.Errors, error) {
	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	redisErr := deleteSessionIDFromRedis(appCtx, foundUser, payload, sessionManagers)
	if redisErr != nil {
		return nil, nil, redisErr
	}

	return deleteSessionIDFromDB(appCtx, o, foundUser, payload)
}

// NewUserSessionRevocation returns a new admin user creator builder
func NewUserSessionRevocation(builder userQueryBuilder) services.UserSessionRevocation {
	return &userSessionRevocation{builder}
}

func revokeSession(appCtx appcontext.AppContext, sessionManager *scs.SessionManager, sessionID string, userID uuid.UUID) error {
	_, exists, err := sessionManager.Store.Find(sessionID)

	if err != nil {
		appCtx.Logger().Error("Error looking up field in Redis for user ID", zap.String("cookie_name", sessionManager.Cookie.Name), zap.String("UserID", userID.String()), zap.Error(err))
		return err
	}

	if !exists {
		appCtx.Logger().Info("Not found in Redis; nothing to revoke", zap.String("cooke_name", sessionManager.Cookie.Name))
	} else {
		appCtx.Logger().Info("Found for user ID; deleting it from Redis", zap.String("cooke_name", sessionManager.Cookie.Name), zap.String("UserID", userID.String()))
		err := sessionManager.Store.Delete(sessionID)
		if err != nil {
			appCtx.Logger().Error("Error deleting field from Redis for user ID", zap.String("cooke_name", sessionManager.Cookie.Name), zap.String("UserID", userID.String()), zap.Error(err))
			return err
		}
	}
	return nil
}

func deleteSessionIDFromRedis(appCtx appcontext.AppContext, user models.User, payload *adminmessages.UserUpdatePayload, sessionManagers [3]*scs.SessionManager) error {

	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession {
		err := revokeSession(appCtx, sessionManagers[1], user.CurrentAdminSessionID, user.ID)
		if err != nil {
			return err
		}
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession {
		err := revokeSession(appCtx, sessionManagers[2], user.CurrentOfficeSessionID, user.ID)
		if err != nil {
			return err
		}
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession {
		err := revokeSession(appCtx, sessionManagers[0], user.CurrentMilSessionID, user.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteSessionIDFromDB(appCtx appcontext.AppContext, o *userSessionRevocation, user models.User, payload *adminmessages.UserUpdatePayload) (*models.User, *validate.Errors, error) {
	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession {
		user.CurrentAdminSessionID = ""
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession {
		user.CurrentOfficeSessionID = ""
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession {
		user.CurrentMilSessionID = ""
	}

	verrs, err := o.builder.UpdateOne(appCtx, &user, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &user, nil, nil
}
