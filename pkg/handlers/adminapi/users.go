package adminapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
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
		CurrentAdminSessionID:  handlers.FmtString(o.CurrentAdminSessionID),
		CurrentMilSessionID:    handlers.FmtString(o.CurrentMilSessionID),
		CurrentOfficeSessionID: handlers.FmtString(o.CurrentOfficeSessionID),
	}
}

// GetUserHandler returns an user via GET /users/{userID}
type GetUserHandler struct {
	handlers.HandlerContext
	services.UserFetcher
	services.NewQueryFilter
}

// Handle retrieves a specific user
func (h GetUserHandler) Handle(params userop.GetUserParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	userID := uuid.FromStringOrNil(params.UserID.String())
	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", userID)}
	user, err := h.UserFetcher.FetchUser(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
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

// Handle lists all users
func (h IndexUsersHandler) Handle(params userop.IndexUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := h.generateQueryFilters(params.Filter, logger)

	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	ordering := query.NewQueryOrder(params.Sort, params.Order)
	pagination := h.NewPagination(params.Page, params.PerPage)

	var users models.Users
	err := h.ListFetcher.FetchRecordList(&users, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalUsersCount, err := h.ListFetcher.FetchRecordCount(&users, queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedUsersCount := len(users)

	payload := make(adminmessages.Users, queriedUsersCount)

	for i, s := range users {
		payload[i] = payloadForUserModel(s)
	}

	return userop.NewIndexUsersOK().WithContentRange(fmt.Sprintf("users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedUsersCount, totalUsersCount)).WithPayload(payload)
}

func (h IndexUsersHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		Search string `json:"search"`
	}

	f := Filter{}
	var queryFilters []services.QueryFilter
	if filters == nil {
		return queryFilters
	}
	b := []byte(*filters)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fs := fmt.Sprintf("%v", filters)
		logger.Warn("unable to decode param", zap.Error(err),
			zap.String("filters", fs))
	}

	if f.Search != "" {
		_, err := uuid.FromString(f.Search)
		if err != nil {
			queryFilters = append(queryFilters, query.NewQueryFilter("login_gov_email", "=", f.Search))
		} else {
			queryFilters = append(queryFilters, query.NewQueryFilter("id", "=", f.Search))
		}
	}
	return queryFilters
}

// RevokeUserSessionHandler is the handler for creating users.
type RevokeUserSessionHandler struct {
	handlers.HandlerContext
	services.UserSessionRevocation
	services.NewQueryFilter
}

// Handle revokes a user session
func (h RevokeUserSessionHandler) Handle(params userop.RevokeUserSessionParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	payload := params.User

	userID, err := uuid.FromString(params.UserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.UserID.String()), zap.Error(err))
	}

	sessionStore := h.SessionManager(session).Store
	updatedUser, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(userID, payload, sessionStore)
	if revokeErr != nil || validationErrors != nil {
		fmt.Printf("%#v", validationErrors)
		logger.Error("Error revoking user session", zap.Error(revokeErr))
		return userop.NewRevokeUserSessionInternalServerError()
	}

	returnPayload := payloadForUserModel(*updatedUser)
	return userop.NewRevokeUserSessionOK().WithPayload(returnPayload)
}
