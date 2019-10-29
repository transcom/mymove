package adminapi

import (
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	adminuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/admin_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForAdminUserModel(o models.AdminUser) *adminmessages.AdminUser {
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
		payload[i] = payloadForAdminUserModel(s)
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

	payload := payloadForAdminUserModel(adminUser)

	return adminuserop.NewGetAdminUserOK().WithPayload(payload)
}

type CreateAdminUserHandler struct {
	handlers.HandlerContext
	services.AdminUserCreator
	services.NewQueryFilter
}

func (h CreateAdminUserHandler) Handle(params adminuserop.CreateAdminUserParams) middleware.Responder {
	payload := params.AdminUser
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	organizationID, err := uuid.FromString(payload.OrganizationID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", payload.OrganizationID.String()), zap.Error(err))
	}

	adminUser := models.AdminUser{
		LastName:       payload.LastName,
		FirstName:      payload.FirstName,
		Email:          payload.Email,
		Role:           models.SystemAdminRole,
		OrganizationID: &organizationID,
	}

	organizationIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", organizationID),
	}

	createdAdminUser, verrs, err := h.AdminUserCreator.CreateAdminUser(&adminUser, organizationIDFilter)
	if err != nil || verrs != nil {
		logger.Error("Error saving user", zap.Error(verrs))
		return adminuserop.NewCreateAdminUserInternalServerError()
	}

	logger.Info("Create Admin User", zap.String("admin_user_id", createdAdminUser.ID.String()), zap.String("responsible_user_id", session.AdminUserID.String()), zap.String("event_type", "create_admin_user"))
	returnPayload := payloadForAdminUserModel(*createdAdminUser)
	return adminuserop.NewCreateAdminUserCreated().WithPayload(returnPayload)
}

type UpdateAdminUserHandler struct {
	handlers.HandlerContext
	services.AdminUserUpdater
	services.NewQueryFilter
}

func (h UpdateAdminUserHandler) Handle(params adminuserop.UpdateAdminUserParams) middleware.Responder {
	payload := params.AdminUser
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	adminUserID, err := uuid.FromString(params.AdminUserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.AdminUserID.String()), zap.Error(err))
	}

	// Don't allow Admin Users to deactivate themselves
	if adminUserID == session.AdminUserID && payload.Deactivated {
		return adminuserop.NewUpdateAdminUserForbidden()
	}

	adminUser := models.AdminUser{
		ID:          adminUserID,
		LastName:    payload.LastName,
		FirstName:   payload.FirstName,
		Deactivated: payload.Deactivated,
	}

	updatedAdminUser, verrs, err := h.AdminUserUpdater.UpdateAdminUser(&adminUser)

	if err != nil || verrs != nil {
		fmt.Printf("%#v", verrs)
		logger.Error("Error saving user", zap.Error(err))
		return adminuserop.NewUpdateAdminUserInternalServerError()
	}

	logger.Info("Update admin User", zap.String("admin_user_id", updatedAdminUser.ID.String()), zap.String("responsible_user_id", session.AdminUserID.String()), zap.String("event_type", "update_admin_user"))
	returnPayload := payloadForAdminUserModel(*updatedAdminUser)

	return adminuserop.NewUpdateAdminUserOK().WithPayload(returnPayload)
}
