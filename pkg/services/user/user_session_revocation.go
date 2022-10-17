package user

import (
	"github.com/alexedwards/scs/v2"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userSessionRevocation struct {
	builder userQueryBuilder
}

// RevokeUserSession revokes the user's session
func (o *userSessionRevocation) RevokeUserSession(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.UserUpdatePayload, sessionManagers auth.AppSessionManagers) (*models.User, *validate.Errors, error) {
	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	redisErr := deleteSessionsFromRedis(appCtx, foundUser, payload, sessionManagers)
	if redisErr != nil {
		return nil, nil, redisErr
	}

	return deleteSessionIDFromDB(appCtx, o, foundUser, payload)
}

// NewUserSessionRevocation returns a new admin user creator builder
func NewUserSessionRevocation(builder userQueryBuilder) services.UserSessionRevocation {
	return &userSessionRevocation{builder}
}

// deleteSessionIDFromRedis deletes a sessionID from a particular sessionStore
func deleteSessionIDFromRedis(appCtx appcontext.AppContext, app auth.Application, userID uuid.UUID, sessionID string, sessionStore scs.Store) error {
	_, exists, err := sessionStore.Find(sessionID)
	if err != nil {
		appCtx.Logger().Error("Error looking sessionID in redis for app", zap.Any("app", app), zap.String("UserID", userID.String()), zap.Error(err))
		return err
	}

	if !exists {
		appCtx.Logger().Info("Not found in Redis; nothing to revoke", zap.Any("app", app))
	} else {
		appCtx.Logger().Info("Found for user ID; deleting it from Redis", zap.Any("app", app), zap.String("UserID", userID.String()))
		err := sessionStore.Delete(sessionID)
		if err != nil {
			appCtx.Logger().Error("Error deleting session from Redis for user ID", zap.Any("app", app), zap.String("UserID", userID.String()), zap.Error(err))
			return err
		}
	}

	return nil
}

// deleteSessionsFromRedis deletes all sessions as configured in the
// payload form the sessionManagers
func deleteSessionsFromRedis(appCtx appcontext.AppContext, user models.User, payload *adminmessages.UserUpdatePayload, sessionManagers auth.AppSessionManagers) error {
	userID := user.ID

	var adminErr, officeErr, milErr error
	if payload.RevokeAdminSession != nil && *payload.RevokeAdminSession {
		adminErr = deleteSessionIDFromRedis(appCtx, auth.AdminApp, userID, user.CurrentAdminSessionID, sessionManagers.Admin.Store())
	}

	if payload.RevokeOfficeSession != nil && *payload.RevokeOfficeSession {
		officeErr = deleteSessionIDFromRedis(appCtx, auth.OfficeApp, userID, user.CurrentOfficeSessionID, sessionManagers.Office.Store())
	}

	if payload.RevokeMilSession != nil && *payload.RevokeMilSession {
		milErr = deleteSessionIDFromRedis(appCtx, auth.MilApp, userID, user.CurrentMilSessionID, sessionManagers.Mil.Store())
	}

	// wait to check errors at the end so we try to delete all
	// sessions from redis even if one operation fails
	if adminErr != nil {
		return adminErr
	}
	if officeErr != nil {
		return officeErr
	}
	if milErr != nil {
		return milErr
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
