package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
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

// IndexRequestedOfficeUsersHandler returns a list of office users via GET /office_users
type IndexRequestedOfficeUsersHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of office users
func (h IndexRequestedOfficeUsersHandler) Handle(params requested_office_users.IndexRequestedOfficeUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Default behavior for this handler is going to be returning counts for each of the component services as categories
			queryFilters := []services.QueryFilter{query.NewQueryFilter("status", "=", "REQUESTED")}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

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

			return officeuserop.NewIndexOfficeUsersOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalOfficeUsersCount)).WithPayload(payload), nil
		})
}
