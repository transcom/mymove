package adminapi

import (
	"github.com/go-openapi/runtime/middleware"

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
