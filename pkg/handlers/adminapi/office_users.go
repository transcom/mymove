package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// TODO: fill this in
func payloadForOfficeUserModel(o models.OfficeUser) *adminmessages.OfficeUser {
	return &adminmessages.OfficeUser{ID: *handlers.FmtUUID(o.ID)}
}

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerContext
	services.NewQueryFilter
	services.OfficeUserListFetcher
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {

	logger := h.LoggerFromRequest(params.HTTPRequest)

	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{
		h.NewQueryFilter("id", "=", "d874d002-5582-4a91-97d3-786e8f66c763"),
	}

	officeUsers, err := h.OfficeUserListFetcher.FetchOfficeUserList(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := make(adminmessages.OfficeUsers, len(officeUsers))
	for i, s := range officeUsers {
		payload[i] = payloadForOfficeUserModel(s)
	}

	return officeuserop.NewIndexOfficeUsersOK().WithPayload(payload)
}

type CreateOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserCreator
	services.NewQueryFilter
}

func (h CreateOfficeUserHandler) Handle(params officeuserop.CreateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	transporationOfficeID, err := uuid.FromString(payload.TransportationOfficeID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", payload.TransportationOfficeID.String()), zap.Error(err))
	}

	officeUser := models.OfficeUser{
		LastName:               payload.LastName,
		FirstName:              payload.FirstName,
		Telephone:              payload.Telephone,
		Email:                  payload.Email,
		TransportationOfficeID: transporationOfficeID,
	}

	transporationIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", transporationOfficeID),
	}

	createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(&officeUser, transporationIDFilter)
	if err != nil || verrs.HasAny() {
		logger.Error("Error saving user", zap.Error(err), zap.Error(verrs))
		return officeuserop.NewCreateOfficeUserInternalServerError()
	}

	returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
	return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload)
}
