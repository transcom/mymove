package adminapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	adminuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/admin_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForAdminUser(o models.AdminUser) *adminmessages.AdminUser {
	return &adminmessages.AdminUser{
		ID:             handlers.FmtUUID(o.ID),
		FirstName:      handlers.FmtString(o.FirstName),
		LastName:       handlers.FmtString(o.LastName),
		Email:          handlers.FmtString(o.Email),
		UserID:         handlers.FmtUUIDPtr(o.UserID),
		OrganizationID: handlers.FmtUUIDPtr(o.OrganizationID),
		Deactivated:    handlers.FmtBool(o.Deactivated),
		CreatedAt:      handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:      handlers.FmtDateTime(o.UpdatedAt),
	}
}

// IndexAdminUsersHandler returns a list of admin users via GET /admin_users
type IndexAdminUsersHandler struct {
	handlers.HandlerContext
	services.AdminUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of admin users
func (h IndexAdminUsersHandler) Handle(params adminuserop.IndexAdminUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	pagination := h.NewPagination(params.Page, params.PerPage)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	adminUsers, err := h.AdminUserListFetcher.FetchAdminUserList(queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalAdminUsersCount, err := h.DB().Count(&models.AdminUser{})
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedAdminUsersCount := len(adminUsers)

	payload := make(adminmessages.AdminUsers, queriedAdminUsersCount)

	for i, s := range adminUsers {
		payload[i] = payloadForAdminUser(s)
	}

	return adminuserop.NewIndexAdminUsersOK().WithContentRange(fmt.Sprintf("admin users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedAdminUsersCount, totalAdminUsersCount)).WithPayload(payload)
}

type GetAdminUserHandler struct {
	handlers.HandlerContext
	services.AdminUserFetcher
	services.NewQueryFilter
}

func (h GetAdminUserHandler) Handle(params adminuserop.GetAdminUserParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	adminUserID := params.AdminUserID

	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", adminUserID)}

	adminUser, err := h.AdminUserFetcher.FetchAdminUser(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := payloadForAdminUser(adminUser)

	return adminuserop.NewGetAdminUserOK().WithPayload(payload)
}
