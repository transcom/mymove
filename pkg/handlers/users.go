package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	// "go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForUserModel(storer storage.FileStorer, user *models.User, serviceMember *models.ServiceMember) *internalmessages.LoggedInUserPayload {
	var smPayload *internalmessages.ServiceMemberPayload

	if serviceMember != nil {
		smPayload = payloadForServiceMemberModel(storer, *serviceMember)
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	var user *models.User
	serviceMember, err := models.GetFullServiceMemberProfile(h.db, session)
	if err == nil {
		if serviceMember == nil {
			user, err = models.GetUser(h.db, session.UserID)
		} else {
			user = &serviceMember.User
		}
	}

	var response middleware.Responder
	if err != nil {
		fmt.Println("ERROR IN HANDLER", err)
		// h.logger.Error("Error retrieving service_member", zap.Error(err))
		response = userop.NewShowLoggedInUserUnauthorized()
	} else {
		userPayload := payloadForUserModel(h.storage, user, serviceMember)
		response = userop.NewShowLoggedInUserOK().WithPayload(userPayload)
	}
	return response
}
