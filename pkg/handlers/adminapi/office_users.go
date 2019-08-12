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
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForOfficeUserModel(o models.OfficeUser) *adminmessages.OfficeUser {
	return &adminmessages.OfficeUser{
		ID:        *handlers.FmtUUID(o.ID),
		FirstName: o.FirstName,
		LastName:  o.LastName,
		Email:     o.Email,
	}
}

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerContext
	services.OfficeUserListFetcher
	services.NewQueryFilter
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	officeUsers, err := h.OfficeUserListFetcher.FetchOfficeUserList(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	officeUsersCount := len(officeUsers)

	payload := make(adminmessages.OfficeUsers, officeUsersCount)
	for i, s := range officeUsers {
		payload[i] = payloadForOfficeUserModel(s)
	}

	return officeuserop.NewIndexOfficeUsersOK().WithContentRange(fmt.Sprintf("office users 0-%d/%d", officeUsersCount, officeUsersCount)).WithPayload(payload)
}

type GetOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserFetcher
	services.NewQueryFilter
}

func (h GetOfficeUserHandler) Handle(params officeuserop.GetOfficeUserParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	officeUserID := params.OfficeUserID

	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", officeUserID)}

	officeUser, err := h.OfficeUserFetcher.FetchOfficeUser(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := payloadForOfficeUserModel(officeUser)

	return officeuserop.NewGetOfficeUserOK().WithPayload(payload)
}

type CreateOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserCreator
	services.NewQueryFilter
}

func (h CreateOfficeUserHandler) Handle(params officeuserop.CreateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	transportationOfficeID, err := uuid.FromString(payload.TransportationOfficeID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", payload.TransportationOfficeID.String()), zap.Error(err))
	}

	officeUser := models.OfficeUser{
		LastName:               payload.LastName,
		FirstName:              payload.FirstName,
		Telephone:              payload.Telephone,
		Email:                  payload.Email,
		TransportationOfficeID: transportationOfficeID,
	}

	transportationIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", transportationOfficeID),
	}

	createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(&officeUser, transportationIDFilter)
	if err != nil || verrs != nil {
		logger.Error("Error saving user", zap.Error(err))
		return officeuserop.NewCreateOfficeUserInternalServerError()
	}

	returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
	return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload)
}
