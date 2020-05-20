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

// RevokeUserSessionHandler is the handler for creating users.
type RevokeUserSessionHandler struct {
	handlers.HandlerContext
	services.SessionRevocation
	services.NewQueryFilter
}

// Handle revokes a user session
func (h RevokeUserSessionHandler) Handle(params userop.RevokeUserSessionParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.User

	userID, err := uuid.FromString(params.UserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.UserID.String()), zap.Error(err))
	}

	updatedUser, validationErrors, revokeErr := h.SessionRevocation.RevokeUserSession(userID, payload, h.RedisPool())
	if revokeErr != nil || validationErrors != nil {
		fmt.Printf("%#v", validationErrors)
		logger.Error("Error saving user", zap.Error(err))
		return userop.NewRevokeUserSessionInternalServerError()
	}

	returnPayload := payloadForUserModel(*updatedUser)
	return userop.NewRevokeUserSessionOK().WithPayload(returnPayload)
}
