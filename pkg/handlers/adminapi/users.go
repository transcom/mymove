package adminapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/services/audit"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForUserModel(o models.User) *adminmessages.User {
	return &adminmessages.User{
		ID:                     *handlers.FmtUUID(o.ID),
		LoginGovEmail:          handlers.FmtString(o.LoginGovEmail),
		Active:                 handlers.FmtBool(o.Active),
		CreatedAt:              handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:              handlers.FmtDateTime(o.UpdatedAt),
		CurrentAdminSessionID:  handlers.FmtString(o.CurrentAdminSessionID),
		CurrentMilSessionID:    handlers.FmtString(o.CurrentMilSessionID),
		CurrentOfficeSessionID: handlers.FmtString(o.CurrentOfficeSessionID),
	}
}

// GetUserHandler returns a user via GET /users/{userID}
type GetUserHandler struct {
	handlers.HandlerContext
	services.UserFetcher
	services.NewQueryFilter
}

// Handle retrieves a specific user
func (h GetUserHandler) Handle(params userop.GetUserParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	userID := uuid.FromStringOrNil(params.UserID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", userID)}
	user, err := h.UserFetcher.FetchUser(appCtx, queryFilters)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	payload := payloadForUserModel(user)
	return userop.NewGetUserOK().WithPayload(payload)
}

// IndexUsersHandler returns a list of users via GET /users
type IndexUsersHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var usersFilterConverters = map[string]func(string) []services.QueryFilter{
	"search": func(content string) []services.QueryFilter {
		if _, err := uuid.FromString(content); err != nil {
			return []services.QueryFilter{query.NewQueryFilter("login_gov_email", "=", content)}
		}
		return []services.QueryFilter{query.NewQueryFilter("id", "=", content)}
	},
}

// Handle lists all users
func (h IndexUsersHandler) Handle(params userop.IndexUsersParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, usersFilterConverters)

	ordering := query.NewQueryOrder(params.Sort, params.Order)
	pagination := h.NewPagination(params.Page, params.PerPage)

	var users models.Users
	err := h.ListFetcher.FetchRecordList(appCtx, &users, queryFilters, nil, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	totalUsersCount, err := h.ListFetcher.FetchRecordCount(appCtx, &users, queryFilters)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	queriedUsersCount := len(users)

	payload := make(adminmessages.Users, queriedUsersCount)

	for i, s := range users {
		payload[i] = payloadForUserModel(s)
	}

	return userop.NewIndexUsersOK().WithContentRange(fmt.Sprintf("users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedUsersCount, totalUsersCount)).WithPayload(payload)
}

// UpdateUserHandler is the handler for updating users.
type UpdateUserHandler struct {
	handlers.HandlerContext
	services.UserSessionRevocation
	services.UserUpdater
	services.NewQueryFilter
}

// Handle updates a user's Active status and/or their sessions
func (h UpdateUserHandler) Handle(params userop.UpdateUserParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	payload := params.User

	// Check that the uuid provided is valid and get user model
	userID, err := uuid.FromString(params.UserID.String())
	if err != nil {
		appCtx.Logger().Error("updateUserHandler Error", zap.Error(fmt.Errorf("Could not parse ID: %s", params.UserID.String())))
		return userop.NewUpdateUserUnprocessableEntity()
	}
	dbUser := models.User{}
	err = appCtx.DB().Find(&dbUser, userID)
	if err != nil {
		appCtx.Logger().Error("updateUserHandler Error", zap.Error(fmt.Errorf("No user found for ID: %s", params.UserID.String())))
		return userop.NewUpdateUserNotFound()
	}

	// Update all properties from the payload that are not related to revoking a session.
	// Currently, only updating the Active property is supported.
	// If you want to add support for additional properties, edit UpdateUser.
	// Also we need to retrieve the user's original status from the db to correctly create the update model
	user, err := payloads.UserModel(payload, userID, dbUser.Active)
	if err != nil {
		appCtx.Logger().Error("updateUserHandler Error", zap.Error(err))
		return userop.NewUpdateUserUnprocessableEntity()
	}

	_, verrs, err := h.UpdateUser(appCtx, userID, user)
	if verrs != nil || err != nil {
		appCtx.Logger().Error(fmt.Sprintf("Error updating user %s", params.UserID.String()), zap.Error(err))
	}
	// We don't return because we should still try to revoke sessions

	// If we've set the user's active status to false, we should also revoke their sessions.
	// Update the payload properties so session revocation is triggered.
	// Even if updating the active status fails, we can try to revoke their session.
	if !user.Active {
		revoke := true

		payload.RevokeAdminSession = &revoke
		payload.RevokeOfficeSession = &revoke
		payload.RevokeMilSession = &revoke
	}

	sessionStore := h.SessionManager(appCtx.Session()).Store
	updatedUser, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(appCtx, userID, payload, sessionStore)
	if revokeErr != nil || validationErrors != nil {
		appCtx.Logger().Error("Error revoking user session", zap.Error(revokeErr), zap.Error(verrs))
		return userop.NewUpdateUserInternalServerError()
	}

	// Log if the account was enabled or disabled (POAM requirement)
	if payload.Active != nil {
		_, err = audit.CaptureAccountStatus(appCtx, updatedUser, *payload.Active, params.HTTPRequest)
		if err != nil {
			appCtx.Logger().Error("Error capturing account status audit record in UpdateUserHandler", zap.Error(err))
		}
	}

	_, err = audit.Capture(appCtx, updatedUser, params.User, params.HTTPRequest)
	if err != nil {
		appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
	}

	returnPayload := payloadForUserModel(*updatedUser)
	return userop.NewUpdateUserOK().WithPayload(returnPayload)
}
