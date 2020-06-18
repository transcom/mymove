package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForUserModel(o models.User) *adminmessages.User {
	return &adminmessages.User{
		ID:                     *handlers.FmtUUID(o.ID),
		LoginGovEmail:          handlers.FmtString(o.LoginGovEmail),
		CurrentAdminSessionID:  handlers.FmtString(o.CurrentAdminSessionID),
		CurrentMilSessionID:    handlers.FmtString(o.CurrentMilSessionID),
		CurrentOfficeSessionID: handlers.FmtString(o.CurrentOfficeSessionID),
	}
}

// GetUserHandler returns an user via GET /users/{userID}
type GetUserHandler struct {
	handlers.HandlerContext
	services.UserFetcher
	services.NewQueryFilter
}

// Handle retrieves a specific user
func (h GetUserHandler) Handle(params userop.GetUserParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	userID := uuid.FromStringOrNil(params.UserID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", userID)}
	user, err := h.UserFetcher.FetchUser(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	payload := payloadForUserModel(user)
	return userop.NewGetUserOK().WithPayload(payload)
}

// RevokeUserSessionHandler is the handler for creating users.
type RevokeUserSessionHandler struct {
	handlers.HandlerContext
	services.UserSessionRevocation
	services.NewQueryFilter
}

// Handle revokes a user session
func (h RevokeUserSessionHandler) Handle(params userop.RevokeUserSessionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	payload := params.User

	userID, err := uuid.FromString(params.UserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.UserID.String()), zap.Error(err))
	}

	sessionStore := h.SessionManager(session).Store
	updatedUser, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(userID, payload, sessionStore)
	if revokeErr != nil || validationErrors != nil {
		fmt.Printf("%#v", validationErrors)
		logger.Error("Error revoking user session", zap.Error(revokeErr))
		return userop.NewRevokeUserSessionInternalServerError()
	}

	returnPayload := payloadForUserModel(*updatedUser)
	return userop.NewRevokeUserSessionOK().WithPayload(returnPayload)
}
