package adminapi

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForRole(r roles.Role) *adminmessages.Role {
	roleType := string(r.RoleType)
	roleName := string(r.RoleName)
	return &adminmessages.Role{
		ID:        handlers.FmtUUID(r.ID),
		RoleType:  &roleType,
		RoleName:  &roleName,
		CreatedAt: handlers.FmtDateTime(r.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(r.UpdatedAt),
	}
}

func payloadForOfficeUserModel(o models.OfficeUser) *adminmessages.OfficeUser {
	var user models.User
	if o.UserID != nil {
		user = o.User
	}
	payload := &adminmessages.OfficeUser{
		ID:                     handlers.FmtUUID(o.ID),
		FirstName:              handlers.FmtString(o.FirstName),
		MiddleInitials:         handlers.FmtStringPtr(o.MiddleInitials),
		LastName:               handlers.FmtString(o.LastName),
		Telephone:              handlers.FmtString(o.Telephone),
		Email:                  handlers.FmtString(o.Email),
		TransportationOfficeID: handlers.FmtUUID(o.TransportationOfficeID),
		Active:                 handlers.FmtBool(o.Active),
		CreatedAt:              handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:              handlers.FmtDateTime(o.UpdatedAt),
	}
	if o.UserID != nil {
		userIDFmt := handlers.FmtUUID(*o.UserID)
		if userIDFmt != nil {
			payload.UserID = *userIDFmt
		}
	}
	for _, role := range user.Roles {
		payload.Roles = append(payload.Roles, payloadForRole(role))
	}
	return payload
}

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerContext
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := h.generateQueryFilters(params.Filter, logger)

	pagination := h.NewPagination(params.Page, params.PerPage)
	associations := query.NewQueryAssociations([]services.QueryAssociation{})
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	var officeUsers models.OfficeUsers
	err := h.ListFetcher.FetchRecordList(&officeUsers, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalOfficeUsersCount, err := h.ListFetcher.FetchRecordCount(&officeUsers, queryFilters)
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

// GetOfficeUserHandler retrieves office user handler
type GetOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserFetcher
	services.NewQueryFilter
}

// Handle retrieves an office user
func (h GetOfficeUserHandler) Handle(params officeuserop.GetOfficeUserParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	officeUserID := params.OfficeUserID

	queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", officeUserID)}

	officeUser, err := h.OfficeUserFetcher.FetchOfficeUser(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	userError := h.DB().Load(&officeUser, "User")
	if userError != nil {
		return handlers.ResponseForError(logger, userError)
	}
	//todo: we want to move this query out of the handler and into querybuilder, if possible
	roleError := h.DB().Q().Join("users_roles", "users_roles.role_id = roles.id").
		Where("users_roles.deleted_at IS NULL AND users_roles.user_id = ?", (officeUser.User.ID)).
		All(&officeUser.User.Roles)
	if roleError != nil {
		return handlers.ResponseForError(logger, roleError)
	}
	payload := payloadForOfficeUserModel(officeUser)

	return officeuserop.NewGetOfficeUserOK().WithPayload(payload)
}

// CreateOfficeUserHandler creates an office user
type CreateOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserCreator
	services.NewQueryFilter
	services.UserRoleAssociator
}

// Handle creates an office user
func (h CreateOfficeUserHandler) Handle(params officeuserop.CreateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	transportationOfficeID, err := uuid.FromString(payload.TransportationOfficeID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", payload.TransportationOfficeID.String()), zap.Error(err))
		return officeuserop.NewCreateOfficeUserUnprocessableEntity()
	}

	if len(payload.Roles) == 0 {
		logger.Error("At least one office user role is required")
		return officeuserop.NewCreateOfficeUserUnprocessableEntity()
	}

	updatedRoles := rolesPayloadToModel(payload.Roles)
	if len(updatedRoles) == 0 {
		logger.Error("No roles were matched from payload")
		return officeuserop.NewCreateOfficeUserUnprocessableEntity()
	}

	officeUser := models.OfficeUser{
		LastName:               payload.LastName,
		FirstName:              payload.FirstName,
		Telephone:              payload.Telephone,
		Email:                  payload.Email,
		TransportationOfficeID: transportationOfficeID,
		Active:                 true,
	}

	transportationIDFilter := []services.QueryFilter{
		h.NewQueryFilter("id", "=", transportationOfficeID),
	}

	createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(&officeUser, transportationIDFilter)
	if verrs != nil {
		validationError := &adminmessages.ValidationError{
			InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
		}

		validationError.Title = handlers.FmtString(handlers.ValidationErrMessage)
		validationError.Detail = handlers.FmtString("The information you provided is invalid.")
		validationError.Instance = handlers.FmtUUID(h.GetTraceID())

		return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError)
	}

	if err != nil {
		logger.Error("Error saving user", zap.Error(err))
		return officeuserop.NewCreateOfficeUserInternalServerError()
	}

	_, err = h.UserRoleAssociator.UpdateUserRoles(*createdOfficeUser.UserID, updatedRoles)
	if err != nil {
		logger.Error("Error updating user roles", zap.Error(err))
		return officeuserop.NewUpdateOfficeUserInternalServerError()
	}

	_, err = audit.Capture(createdOfficeUser, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Error capturing audit record", zap.Error(err))
	}

	returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
	return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload)
}

// UpdateOfficeUserHandler updates an office user
type UpdateOfficeUserHandler struct {
	handlers.HandlerContext
	services.OfficeUserUpdater
	services.NewQueryFilter
	services.UserRoleAssociator
}

// Handle updates an office user
func (h UpdateOfficeUserHandler) Handle(params officeuserop.UpdateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	officeUserID, err := uuid.FromString(params.OfficeUserID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.OfficeUserID.String()), zap.Error(err))
	}

	updatedOfficeUser, verrs, err := h.OfficeUserUpdater.UpdateOfficeUser(officeUserID, payload)

	if err != nil || verrs != nil {
		fmt.Printf("%#v", verrs)
		logger.Error("Error saving user", zap.Error(err))
		return officeuserop.NewUpdateOfficeUserInternalServerError()
	}
	if updatedOfficeUser.UserID != nil && payload.Roles != nil {
		updatedRoles := rolesPayloadToModel(payload.Roles)
		_, err = h.UserRoleAssociator.UpdateUserRoles(*updatedOfficeUser.UserID, updatedRoles)
		if err != nil {
			logger.Error("Error updating user roles", zap.Error(err))
			return officeuserop.NewUpdateOfficeUserInternalServerError()
		}
	}

	_, err = audit.Capture(updatedOfficeUser, payload, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Error capturing audit record", zap.Error(err))
	}

	returnPayload := payloadForOfficeUserModel(*updatedOfficeUser)

	return officeuserop.NewUpdateOfficeUserOK().WithPayload(returnPayload)
}

func rolesPayloadToModel(payload []*adminmessages.OfficeUserRole) []roles.RoleType {
	var rt []roles.RoleType
	for _, role := range payload {
		if role.RoleType != nil {
			rt = append(rt, roles.RoleType(*role.RoleType))
		}
	}
	return rt
}

// generateQueryFilters is helper to convert filter params from a json string
// of the form `{"search": "example1@example.com"}` to an array of services.QueryFilter
func (h IndexOfficeUsersHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
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
		nameSearch := fmt.Sprintf("%s%%", f.Search)
		queryFilters = append(queryFilters, query.NewQueryFilter("email", "ILIKE", fmt.Sprintf("%%%s%%", f.Search)))
		queryFilters = append(queryFilters, query.NewQueryFilter("first_name", "ILIKE", nameSearch))
		queryFilters = append(queryFilters, query.NewQueryFilter("last_name", "ILIKE", nameSearch))
	}
	return queryFilters
}
