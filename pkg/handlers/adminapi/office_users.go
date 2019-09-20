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
		ID:             handlers.FmtUUID(o.ID),
		FirstName:      handlers.FmtString(o.FirstName),
		MiddleInitials: handlers.FmtStringPtr(o.MiddleInitials),
		LastName:       handlers.FmtString(o.LastName),
		Telephone:      handlers.FmtString(o.Telephone),
		Email:          handlers.FmtString(o.Email),
		Disabled:       handlers.FmtBool(o.Disabled),
		CreatedAt:      handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:      handlers.FmtDateTime(o.UpdatedAt),
	}
}

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerContext
	services.OfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	pagination := h.NewPagination(params.Page, params.PerPage)

	officeUsers, err := h.OfficeUserListFetcher.FetchOfficeUserList(queryFilters, pagination)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalOfficeUsersCount, err := h.DB().Count(&models.OfficeUser{})
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedOfficeUsersCount := len(officeUsers)

	payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

	for i, s := range officeUsers {
		payload[i] = payloadForOfficeUserModel(s)
	}

	return officeuserop.NewIndexOfficeUsersOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalOfficeUsersCount)).WithPayload(payload)
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
	if verrs != nil {
		payload := &adminmessages.ValidationError{
			InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
		}

		payload.Title = handlers.FmtString(handlers.ValidationErrMessage)
		payload.Detail = handlers.FmtString("The information you provided is invalid.")

		return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(payload)
	}

	if err != nil {
		return officeuserop.NewCreateOfficeUserInternalServerError()
	}

	returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
	return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload)
}

type UpdateOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserUpdater
	services.NewQueryFilter
}

func (h UpdateOfficeUserHandler) Handle(params officeuserop.UpdateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	officeUserID, err := uuid.FromString(params.OfficeUserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.OfficeUserID.String()), zap.Error(err))
	}

	officeUser := models.OfficeUser{
		ID:             officeUserID,
		MiddleInitials: handlers.FmtStringPtr(payload.MiddleInitials),
		LastName:       payload.LastName,
		FirstName:      payload.FirstName,
		Telephone:      payload.Telephone,
	}

	updatedOfficeUser, verrs, err := h.OfficeUserUpdater.UpdateOfficeUser(&officeUser)
	if err != nil || verrs != nil {
		fmt.Printf("%#v", verrs)
		logger.Error("Error saving user", zap.Error(err))
		return officeuserop.NewUpdateOfficeUserInternalServerError()
	}

	returnPayload := payloadForOfficeUserModel(*updatedOfficeUser)

	return officeuserop.NewUpdateOfficeUserOK().WithPayload(returnPayload)
}
