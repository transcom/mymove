package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUserModel(storage FileStorer, user models.User, serviceMember *models.ServiceMember) *internalmessages.LoggedInUserPayload {
	var smPayload *internalmessages.ServiceMemberPayload

	if serviceMember != nil {
		smPayload = payloadForServiceMemberModel(storage, user, *serviceMember)
	}

	userPayload := internalmessages.LoggedInUserPayload{
		ID:            fmtUUID(user.ID),
		CreatedAt:     fmtDateTime(user.CreatedAt),
		ServiceMember: smPayload,
		UpdatedAt:     fmtDateTime(user.UpdatedAt),
	}
	return &userPayload
}

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	user, ok := auth.GetUser(params.HTTPRequest.Context())
	if !ok {
		return userop.NewShowLoggedInUserInternalServerError()
	}
	serviceMember, err := user.GetServiceMemberProfile(h.db)
	if err != nil {
		h.logger.Error("Error retrieving service_member", zap.Error(err))
		response := userop.NewShowLoggedInUserUnauthorized()
		return response
	}

	userPayload := payloadForUserModel(h.storage, user, serviceMember)
	response := userop.NewShowLoggedInUserOK().WithPayload(userPayload)
	return response
}
