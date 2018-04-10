package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUserModel(user models.User, serviceMember *models.ServiceMember) internalmessages.UserPayload {
	var smPayload *internalmessages.ServiceMemberPayload

	if serviceMember != nil {
		smp := payloadForServiceMemberModel(user, *serviceMember)
		smPayload = &smp
	}

	userPayload := internalmessages.UserPayload{
		ID:            fmtUUID(user.ID),
		CreatedAt:     fmtDateTime(user.CreatedAt),
		ServiceMember: smPayload,
		UpdatedAt:     fmtDateTime(user.UpdatedAt),
	}
	return userPayload
}

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	user, ok := context.GetUser(params.HTTPRequest.Context())
	if !ok {
		return userop.NewShowLoggedInUserInternalServerError()
	}
	serviceMember, err := user.GetServiceMemberProfile(h.db)
	if err != nil {
		h.logger.Error("Error retrieving service_member", zap.Error(err))
		response := userop.NewShowLoggedInUserUnauthorized()
		return response
	}

	userPayload := payloadForUserModel(user, serviceMember)
	response := userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
	return response
}
