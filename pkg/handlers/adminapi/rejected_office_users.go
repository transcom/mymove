package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/rejected_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForRejectedOfficeUserModel(o models.OfficeUser) *adminmessages.OfficeUser {
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
		Status:                 (*string)(o.Status),
		Edipi:                  handlers.FmtStringPtr(o.EDIPI),
		OtherUniqueID:          handlers.FmtStringPtr(o.OtherUniqueID),
		RejectionReason:        handlers.FmtStringPtr(o.RejectionReason),
		CreatedAt:              *handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:              *handlers.FmtDateTime(o.UpdatedAt),
		RejectedOn:             *handlers.FmtDateTime(o.RejectedOn),
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

// IndexRejectedOfficeUsersHandler returns a list of rejected office users via GET /rejected_office_users
type IndexRejectedOfficeUsersHandler struct {
	handlers.HandlerConfig
	services.RejectedOfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var rejectedOfficeUserFilterConverters = map[string]func(string) []services.QueryFilter{
	"search": func(content string) []services.QueryFilter {
		nameSearch := fmt.Sprintf("%s%%", content)
		return []services.QueryFilter{
			query.NewQueryFilter("email", "ILIKE", fmt.Sprintf("%%%s%%", content)),
			query.NewQueryFilter("first_name", "ILIKE", nameSearch),
			query.NewQueryFilter("last_name", "ILIKE", nameSearch),
		}
	},
}

// Handle retrieves a list of rejected office users
func (h IndexRejectedOfficeUsersHandler) Handle(params rejected_office_users.IndexRejectedOfficeUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// adding in filters for when a search or filtering is done
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, rejectedOfficeUserFilterConverters)

			// We only want users that are in a REJECTED status
			queryFilters = append(queryFilters, query.NewQueryFilter("status", "=", "REJECTED"))

			// adding in pagination for the UI
			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			// need to also get the user's roles
			queryAssociations := query.NewQueryAssociationsPreload([]services.QueryAssociation{
				query.NewQueryAssociation("User.Roles"),
			})

			officeUsers, err := h.RejectedOfficeUserListFetcher.FetchRejectedOfficeUsersList(appCtx, queryFilters, queryAssociations, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficeUsersCount, err := h.RejectedOfficeUserListFetcher.FetchRejectedOfficeUsersCount(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, s := range officeUsers {
				payload[i] = payloadForRejectedOfficeUserModel(s)
			}

			return rejected_office_users.NewIndexRejectedOfficeUsersOK().WithContentRange(fmt.Sprintf("rejected office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalOfficeUsersCount)).WithPayload(payload), nil
		})
}

// GetRejectedOfficeUserHandler returns a list of office users via GET /rejected_office_users/{officeUserId}
type GetRejectedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RejectedOfficeUserFetcher
	services.RoleAssociater
	services.NewQueryFilter
}

// Handle retrieves a single rejected office user
func (h GetRejectedOfficeUserHandler) Handle(params rejected_office_users.GetRejectedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			rejectedOfficeUserID := params.OfficeUserID

			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", rejectedOfficeUserID)}
			rejectedOfficeUser, err := h.RejectedOfficeUserFetcher.FetchRejectedOfficeUser(appCtx, queryFilters)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			roles, err := h.RoleAssociater.FetchRolesForUser(appCtx, *rejectedOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return rejected_office_users.NewGetRejectedOfficeUserBadRequest(), err
			}

			rejectedOfficeUser.User.Roles = roles

			payload := payloadForRejectedOfficeUserModel(rejectedOfficeUser)

			return rejected_office_users.NewGetRejectedOfficeUserOK().WithPayload(payload), nil
		})
}
