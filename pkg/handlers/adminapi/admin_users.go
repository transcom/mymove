package adminapi

import (
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	adminuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/admin_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
)

func payloadForAdminUserModel(o models.AdminUser) *adminmessages.AdminUser {
	return &adminmessages.AdminUser{
		ID:             handlers.FmtUUID(o.ID),
		FirstName:      handlers.FmtString(o.FirstName),
		LastName:       handlers.FmtString(o.LastName),
		Email:          handlers.FmtString(o.Email),
		UserID:         handlers.FmtUUIDPtr(o.UserID),
		OrganizationID: handlers.FmtUUIDPtr(o.OrganizationID),
		Active:         handlers.FmtBool(o.Active),
		CreatedAt:      handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:      handlers.FmtDateTime(o.UpdatedAt),
	}
}

// IndexAdminUsersHandler returns a list of admin users via GET /admin_users
type IndexAdminUsersHandler struct {
	handlers.HandlerConfig
	services.AdminUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of admin users
func (h IndexAdminUsersHandler) Handle(params adminuserop.IndexAdminUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
			queryFilters := []services.QueryFilter{}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			adminUsers, err := h.AdminUserListFetcher.FetchAdminUserList(appCtx, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalAdminUsersCount, err := h.AdminUserListFetcher.FetchAdminUserCount(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedAdminUsersCount := len(adminUsers)

			payload := make(adminmessages.AdminUsers, queriedAdminUsersCount)

			for i, s := range adminUsers {
				payload[i] = payloadForAdminUserModel(s)
			}

			return adminuserop.NewIndexAdminUsersOK().WithContentRange(fmt.Sprintf("admin users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedAdminUsersCount, totalAdminUsersCount)).WithPayload(payload), nil
		})
}

// GetAdminUserHandler retrieves a handler for admin users
type GetAdminUserHandler struct {
	handlers.HandlerConfig
	services.AdminUserFetcher
	services.NewQueryFilter
}

// Handle retrieves a new admin user
func (h GetAdminUserHandler) Handle(params adminuserop.GetAdminUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			adminUserID := params.AdminUserID

			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", adminUserID)}

			adminUser, err := h.AdminUserFetcher.FetchAdminUser(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := payloadForAdminUserModel(adminUser)

			return adminuserop.NewGetAdminUserOK().WithPayload(payload), nil
		})
}

// CreateAdminUserHandler is the handler for creating users.
type CreateAdminUserHandler struct {
	handlers.HandlerConfig
	services.AdminUserCreator
	services.NewQueryFilter
}

// Handle creates an admin user
func (h CreateAdminUserHandler) Handle(params adminuserop.CreateAdminUserParams) middleware.Responder {
	payload := params.AdminUser
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			organizationID, err := uuid.FromString(payload.OrganizationID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", payload.OrganizationID.String()), zap.Error(err))
				return adminuserop.NewCreateAdminUserBadRequest(), err
			}

			adminUser := models.AdminUser{
				LastName:       payload.LastName,
				FirstName:      payload.FirstName,
				Email:          payload.Email,
				Role:           models.SystemAdminRole,
				OrganizationID: &organizationID,
				Active:         true,
			}

			organizationIDFilter := []services.QueryFilter{
				h.NewQueryFilter("id", "=", organizationID),
			}

			createdAdminUser, verrs, err := h.AdminUserCreator.CreateAdminUser(appCtx, &adminUser, organizationIDFilter)
			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(verrs))
				return adminuserop.NewCreateAdminUserInternalServerError(), err
			}

			_, err = audit.Capture(appCtx, createdAdminUser, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForAdminUserModel(*createdAdminUser)
			return adminuserop.NewCreateAdminUserCreated().WithPayload(returnPayload), nil
		})
}

// UpdateAdminUserHandler is the handler for updating users
type UpdateAdminUserHandler struct {
	handlers.HandlerConfig
	services.AdminUserUpdater
	services.NewQueryFilter
}

// Handle updates admin users
func (h UpdateAdminUserHandler) Handle(params adminuserop.UpdateAdminUserParams) middleware.Responder {
	payload := params.AdminUser
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			adminUserID, err := uuid.FromString(params.AdminUserID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", params.AdminUserID.String()), zap.Error(err))
			}

			// Don't allow Admin Users to deactivate themselves
			if adminUserID == appCtx.Session().AdminUserID && payload.Active != nil {
				return adminuserop.NewUpdateAdminUserForbidden(), nil
			}

			updatedAdminUser, verrs, err := h.AdminUserUpdater.UpdateAdminUser(appCtx, adminUserID, payload)

			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err), zap.Error(verrs))
				return adminuserop.NewUpdateAdminUserInternalServerError(), err
			}

			// Log if the account was enabled or disabled (POAM requirement)
			if payload.Active != nil {
				_, err = audit.CaptureAccountStatus(appCtx, updatedAdminUser, *payload.Active, params.HTTPRequest)
				if err != nil {
					appCtx.Logger().Error("Error capturing account status audit record in UpdateAdminUserHandler", zap.Error(err))
				}
			}

			_, err = audit.Capture(appCtx, updatedAdminUser, payload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForAdminUserModel(*updatedAdminUser)

			return adminuserop.NewUpdateAdminUserOK().WithPayload(returnPayload), nil
		})
}
