package handlers

import (
	"github.com/go-openapi/runtime/middleware"

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
		Type:          user.Type,
		UpdatedAt:     fmtDateTime(user.UpdatedAt),
	}
	return userPayload
}

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response := userop.NewShowLoggedInUserUnauthorized()
		return response
	}
	serviceMember, err := user.GetServiceMemberProfile(h.db)
	if err != nil {
		response := userop.NewShowLoggedInUserUnauthorized()
		return response
	}

	userPayload := payloadForUserModel(user, serviceMember)
	response := userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
	return response
}
