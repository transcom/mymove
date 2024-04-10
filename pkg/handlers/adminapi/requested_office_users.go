package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForRequestedOfficeUserModel(o models.OfficeUser) *adminmessages.OfficeUser {
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
		Status:                 handlers.FmtStringPtr(o.Status),
		Edipi:                  handlers.FmtStringPtr(o.EDIPI),
		OtherUniqueID:          handlers.FmtStringPtr(o.OtherUniqueID),
		RejectionReason:        handlers.FmtStringPtr(o.RejectionReason),
		CreatedAt:              *handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:              *handlers.FmtDateTime(o.UpdatedAt),
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

// IndexRequestedOfficeUsersHandler returns a list of requested office users via GET /requested_office_users
type IndexRequestedOfficeUsersHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var requestedOfficeUserFilterConverters = map[string]func(string) []services.QueryFilter{
	"search": func(content string) []services.QueryFilter {
		nameSearch := fmt.Sprintf("%s%%", content)
		return []services.QueryFilter{
			query.NewQueryFilter("email", "ILIKE", fmt.Sprintf("%%%s%%", content)),
			query.NewQueryFilter("first_name", "ILIKE", nameSearch),
			query.NewQueryFilter("last_name", "ILIKE", nameSearch),
		}
	},
}

// Handle retrieves a list of requested office users
func (h IndexRequestedOfficeUsersHandler) Handle(params requested_office_users.IndexRequestedOfficeUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// adding in filters for when a search or filtering is done
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, requestedOfficeUserFilterConverters)

			// We only want users that are in a REQUESTED status
			queryFilters = append(queryFilters, query.NewQueryFilter("status", "=", "REQUESTED"))

			// adding in pagination for the UI
			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			// need to also get the user's roles
			queryAssociations := query.NewQueryAssociationsPreload([]services.QueryAssociation{
				query.NewQueryAssociation("User.Roles"),
			})

			officeUsers, err := h.RequestedOfficeUserListFetcher.FetchRequestedOfficeUsersList(appCtx, queryFilters, queryAssociations, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficeUsersCount, err := h.RequestedOfficeUserListFetcher.FetchRequestedOfficeUsersCount(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, s := range officeUsers {
				payload[i] = payloadForRequestedOfficeUserModel(s)
			}

			return requested_office_users.NewIndexRequestedOfficeUsersOK().WithContentRange(fmt.Sprintf("requested office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalOfficeUsersCount)).WithPayload(payload), nil
		})
}

// GetRequestedOfficeUserHandler returns a list of office users via GET /requested_office_users/{officeUserId}
type GetRequestedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserFetcher
	services.RoleAssociater
	services.NewQueryFilter
}

// Handle retrieves a single requested office user
func (h GetRequestedOfficeUserHandler) Handle(params requested_office_users.GetRequestedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			requestedOfficeUserID := params.OfficeUserID

			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", requestedOfficeUserID)}

			requestedOfficeUser, err := h.RequestedOfficeUserFetcher.FetchRequestedOfficeUser(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			roles, err := h.RoleAssociater.FetchRoles(appCtx, *requestedOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return requested_office_users.NewUpdateRequestedOfficeUserInternalServerError(), err
			}

			requestedOfficeUser.User.Roles = roles

			payload := payloadForRequestedOfficeUserModel(requestedOfficeUser)

			return requested_office_users.NewGetRequestedOfficeUserOK().WithPayload(payload), nil
		})
}

// GetRequestedOfficeUserHandler returns a list of office users via GET /requested_office_users/{officeUserId}
type UpdateRequestedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserUpdater
	services.UserRoleAssociator
	services.RoleAssociater
}

// Handle updates a single requested office user
// this endpoint will be used when an admin is approving/rejecting the user without updates
// as well as approving/rejecting the user with updates
func (h UpdateRequestedOfficeUserHandler) Handle(params requested_office_users.UpdateRequestedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			requestedOfficeUserID, err := uuid.FromString(params.OfficeUserID.String())
			body := params.Body
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", params.OfficeUserID.String()), zap.Error(err))
			}

			// roles are associated with users and not office_users, so we need to handle this logic separately
			updatedRoles := rolesPayloadToModel(body.Roles)
			if len(updatedRoles) == 0 {
				err = apperror.NewBadDataError("No roles were matched from payload")
				appCtx.Logger().Error(err.Error())
				return requested_office_users.NewUpdateRequestedOfficeUserUnprocessableEntity(), err
			}

			requestedOfficeUser, verrs, err := h.RequestedOfficeUserUpdater.UpdateRequestedOfficeUser(appCtx, requestedOfficeUserID, params.Body)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			if verrs != nil {
				appCtx.Logger().Error(err.Error())
				return requested_office_users.NewUpdateRequestedOfficeUserUnprocessableEntity(), verrs
			}

			if requestedOfficeUser.UserID != nil && body.Roles != nil {
				_, err = h.UserRoleAssociator.UpdateUserRoles(appCtx, *requestedOfficeUser.UserID, updatedRoles)
				if err != nil {
					appCtx.Logger().Error("Error updating user roles", zap.Error(err))
					return requested_office_users.NewUpdateRequestedOfficeUserInternalServerError(), err
				}
			}

			roles, err := h.RoleAssociater.FetchRoles(appCtx, *requestedOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return requested_office_users.NewUpdateRequestedOfficeUserInternalServerError(), err
			}

			requestedOfficeUser.User.Roles = roles

			// send the email to the user if their request was rejected
			if params.Body.Status == "REJECTED" {
				err = h.NotificationSender().SendNotification(appCtx,
					notifications.NewOfficeAccountRejected(requestedOfficeUser.ID),
				)
				if err != nil {
					err = apperror.NewBadDataError("problem sending email to rejected office user")
					appCtx.Logger().Error(err.Error())
					return requested_office_users.NewUpdateRequestedOfficeUserUnprocessableEntity(), err
				}
			}

			payload := payloadForRequestedOfficeUserModel(*requestedOfficeUser)

			return requested_office_users.NewUpdateRequestedOfficeUserOK().WithPayload(payload), nil
		})
}
