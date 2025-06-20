package adminapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
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
	}

	if o.RejectedOn != nil {
		payload.RejectedOn = *handlers.FmtDateTime(*o.RejectedOn)
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

var rejectedOfficeUserFilterConverters = map[string]func(string) func(*pop.Query){
	"search": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstSearch, lastSearch, emailSearch := fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.first_name ILIKE ? AND office_users.status = 'REJECTED' OR office_users.email ILIKE ? AND office_users.status = 'REJECTED' OR office_users.last_name ILIKE ? AND office_users.status = 'REJECTED'", firstSearch, lastSearch, emailSearch)
		}
	},
	"emails": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			emailSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.email ILIKE ? AND office_users.status = 'REJECTED'", emailSearch)
		}
	},
	"firstName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.first_name ILIKE ? AND office_users.status = 'REJECTED'", firstNameSearch)
		}
	},
	"lastName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			lastNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.last_name ILIKE ? AND office_users.status = 'REJECTED'", lastNameSearch)
		}
	},
	"offices": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			officeSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("transportation_offices.name ILIKE ? AND office_users.status = 'REJECTED'", officeSearch)
		}
	},
	"rejectionReason": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			rejectionSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.rejection_reason ILIKE ? AND office_users.status = 'REJECTED'", rejectionSearch)
		}
	},
	"rejectedOn": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			trimAllZero, trimDayZero, trimMonthZero, noTrim := fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content)
			query.Where("(TO_CHAR(office_users.rejected_on, 'FMMM/FMDD/YYYY') ILIKE ? OR TO_CHAR(office_users.rejected_on, 'MM/FMDD/YYYY') ILIKE ? OR TO_CHAR(office_users.rejected_on, 'FMMM/DD/YYYY') ILIKE ? OR TO_CHAR(office_users.rejected_on, 'MM/DD/YYYY') ILIKE ?) AND office_users.status = 'REJECTED'", trimAllZero, trimDayZero, trimMonthZero, noTrim)
		}
	},
	"roles": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			roleSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("roles.role_name ILIKE ? AND office_users.status = 'REJECTED'", roleSearch)
		}
	},
}

// Handle retrieves a list of rejected office users
func (h IndexRejectedOfficeUsersHandler) Handle(params rejected_office_users.IndexRejectedOfficeUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			var filtersMap map[string]string
			if params.Filter != nil && *params.Filter != "" {
				err := json.Unmarshal([]byte(*params.Filter), &filtersMap)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), errors.New("invalid filter format")), err
				}
			}

			var filterFuncs []func(*pop.Query)
			for key, filterFunc := range rejectedOfficeUserFilterConverters {
				if filterValue, exists := filtersMap[key]; exists {
					filterFuncs = append(filterFuncs, filterFunc(filterValue))
				}
			}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			officeUsers, count, err := h.RejectedOfficeUserListFetcher.FetchRejectedOfficeUsersList(appCtx, filterFuncs, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, officeUser := range officeUsers {
				payload[i] = payloadForRejectedOfficeUserModel(officeUser)
			}

			return rejected_office_users.NewIndexRejectedOfficeUsersOK().WithContentRange(fmt.Sprintf("rejected office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, count)).WithPayload(payload), nil
		})
}

// GetRejectedOfficeUserHandler returns a list of office users via GET /rejected_office_users/{officeUserId}
type GetRejectedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RejectedOfficeUserFetcher
	services.RoleFetcher
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

			roles, err := h.RoleFetcher.FetchRolesForUser(appCtx, *rejectedOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return rejected_office_users.NewGetRejectedOfficeUserBadRequest(), err
			}

			rejectedOfficeUser.User.Roles = roles

			payload := payloadForRejectedOfficeUserModel(rejectedOfficeUser)

			return rejected_office_users.NewGetRejectedOfficeUserOK().WithPayload(payload), nil
		})
}
