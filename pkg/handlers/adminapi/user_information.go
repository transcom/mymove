package adminapi

import (
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/user"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForUser(u services.UserInformation) *adminmessages.UserInformation {
	return &adminmessages.UserInformation{
		ID: strfmt.UUID(u.UserID.String()),
		User: &adminmessages.User{
			LoginGovEmail:          u.LoginGovEmail,
			CurrentAdminSessionID:  u.CurrentAdminSessionID,
			CurrentOfficeSessionID: u.CurrentOfficeSessionID,
			CurrentMilSessionID:    u.CurrentMilSessionID,
		},
	}
}

// GetUserHandler returns an user via GET /users/{userID}
type GetUserHandler struct {
	handlers.HandlerContext
	services.UserInformationFetcher
}

// Handle retrieves a specific user
func (h GetUserHandler) Handle(params userop.GetUserParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	userID := uuid.FromStringOrNil(params.UserID.String())
	userInformation, err := h.FetchUserInformation(userID)
	if err != nil {
		switch err.(type) {
		case services.NotFoundError:
			logger.Error("adminapi.GetUserHandler not found error:", zap.Error(err))
			return userop.NewGetUserNotFound()
		default:
			logger.Error("adminapi.GetUserHandler error:", zap.Error(err))
			return handlers.ResponseForError(logger, err)
		}
	}
	payload := payloadForUser(userInformation)
	return userop.NewGetUserOK().WithPayload(payload)
}
