package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	customeruserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/customer_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// IndexCustomerUsersHandler returns an user via GET /users/{userID}
type IndexCustomerUsersHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

func payloadForCustomerUserModel(o models.User) *adminmessages.CustomerUser {
	payload := &adminmessages.CustomerUser{
		ID:        handlers.FmtUUID(o.ID),
		Email:     handlers.FmtString(o.LoginGovEmail),
		Active:    handlers.FmtBool(o.Active),
		CreatedAt: handlers.FmtDateTime(o.CreatedAt),
	}

	return payload
}

// Handle lists all users (customers/service members)
func (h IndexCustomerUsersHandler) Handle(params customeruserop.IndexCustomerUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	ordering := query.NewQueryOrder(params.Sort, params.Order)
	pagination := h.NewPagination(params.Page, params.PerPage)

	var users models.Users
	err := h.ListFetcher.FetchRecordList(&users, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalCustomerUsersCount, err := h.ListFetcher.FetchRecordCount(&users, queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedCustomerUsersCount := len(users)

	payload := make(adminmessages.CustomerUsers, queriedCustomerUsersCount)

	for i, s := range users {
		payload[i] = payloadForCustomerUserModel(s)
	}

	return customeruserop.NewIndexCustomerUsersOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedCustomerUsersCount, totalCustomerUsersCount)).WithPayload(payload)
}
