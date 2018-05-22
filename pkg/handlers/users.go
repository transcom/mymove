package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForUserModel(serviceMember *models.ServiceMember) *internalmessages.LoggedInUserPayload {
	var smPayload *internalmessages.ServiceMemberPayload

	if serviceMember != nil {
		smPayload = payloadForServiceMemberModel(*serviceMember)
	}

	userPayload := internalmessages.LoggedInUserPayload{
		ID:            fmtUUID(serviceMember.UserID),
		CreatedAt:     fmtDateTime(serviceMember.User.CreatedAt),
		ServiceMember: smPayload,
		UpdatedAt:     fmtDateTime(serviceMember.User.UpdatedAt),
	}
	return &userPayload
}

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	serviceMember, err := models.GetFullServiceMemberProfile(h.db, session)
	if err != nil {
		h.logger.Error("Error retrieving service_member", zap.Error(err))
		response := userop.NewShowLoggedInUserUnauthorized()
		return response
	}

	userPayload := payloadForUserModel(serviceMember)
	response := userop.NewShowLoggedInUserOK().WithPayload(userPayload)
	return response
}
