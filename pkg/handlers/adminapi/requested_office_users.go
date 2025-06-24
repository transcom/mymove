package adminapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
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
		Status:                 (*string)(o.Status),
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

func getOfficeGroupID() (apiKey, customerGroupID string) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return v.GetString(cli.OktaAPIKeyFlag), v.GetString(cli.OktaOfficeGroupIDFlag)
}

// fetchOrCreateOktaProfile send some requests to the Okta Users API
// checks if a user already exists and if not, an okta account is created with access to the office group/application
func fetchOrCreateOktaProfile(appCtx appcontext.AppContext, params requested_office_users.UpdateRequestedOfficeUserParams) (*models.CreatedOktaUser, error) {
	apiKey, officeGroupID := getOfficeGroupID()

	payload := params.Body
	oktaEmail := payload.Email
	oktaFirstName := payload.FirstName
	oktaLastName := payload.LastName
	oktaPhone := payload.Telephone
	oktaEdipi := payload.Edipi
	oktaGsaId := payload.OtherUniqueID

	provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
	if err != nil {
		return nil, err
	}

	if oktaEmail == nil {
		return nil, fmt.Errorf("required okta email is nil")
	}
	users, err := models.SearchForExistingOktaUsers(appCtx, provider, apiKey, *oktaEmail, &oktaEdipi, &oktaGsaId)
	if err != nil {
		return nil, err
	}

	// if we don't find an exact match, then we need to send back an error because there's something weird in Okta
	// that will require a HDT or manual fixing
	if len(users) == 1 {
		oktaUser := &users[0]
		groups, err := models.GetOktaUserGroups(appCtx, provider, apiKey, oktaUser.ID)
		if err != nil {
			return nil, err
		}

		// checking if user is already in office group
		found := false
		for _, group := range groups {
			if group.ID == officeGroupID {
				found = true
				break
			}
		}

		// if they are not already in the office group, then we need to add them
		// use case of this would be for customer users that are registering for office accounts
		if !found {
			err = models.AddOktaUserToGroup(appCtx, provider, apiKey, officeGroupID, oktaUser.ID)
			if err != nil {
				return nil, err
			}
		}

		return oktaUser, nil
	} else if len(users) > 1 {
		var errMsg error
		if oktaEdipi != "" {
			errMsg = fmt.Errorf("multiple Okta accounts found using email %s and EDIPI: %s", *oktaEmail, oktaEdipi)
		} else if oktaGsaId != "" {
			errMsg = fmt.Errorf("multiple Okta accounts found using email %s and GSA ID: %s", *oktaEmail, oktaGsaId)
		} else {
			errMsg = fmt.Errorf("multiple Okta accounts found with the given email %s, EDIPI %s, and/or GSA ID %s", *oktaEmail, oktaEdipi, oktaGsaId)
		}
		appCtx.Logger().Error("okta account fetch error", zap.Error(errMsg))
		return nil, errMsg
	}

	profile := models.OktaProfile{
		FirstName:   handlers.GetStringOrEmpty(oktaFirstName),
		LastName:    handlers.GetStringOrEmpty(oktaLastName),
		Email:       handlers.GetStringOrEmpty(oktaEmail),
		Login:       handlers.GetStringOrEmpty(oktaEmail),
		MobilePhone: handlers.GetStringOrEmpty(oktaPhone),
		CacEdipi:    oktaEdipi,
		GsaID:       &oktaGsaId,
	}
	oktaPayload := models.OktaUserPayload{
		Profile:  profile,
		GroupIds: []string{officeGroupID},
	}

	return models.CreateOktaUser(appCtx, provider, apiKey, oktaPayload)
}

// IndexRequestedOfficeUsersHandler returns a list of requested office users via GET /requested_office_users
type IndexRequestedOfficeUsersHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
	services.TransportationOfficesFetcher
	services.RoleFetcher
}

var requestedOfficeUserFilterConverters = map[string]func(string) func(*pop.Query){
	"search": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstSearch, lastSearch, emailSearch := fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content)
			query.Where("(office_users.first_name ILIKE ? OR office_users.last_name ILIKE ? OR office_users.email ILIKE ?) AND office_users.status = 'REQUESTED'", firstSearch, lastSearch, emailSearch)
		}
	},
	"email": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			emailSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.email ILIKE ? AND office_users.status = 'REQUESTED'", emailSearch)
		}
	},
	"firstName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.first_name ILIKE ? AND office_users.status = 'REQUESTED'", firstNameSearch)
		}
	},
	"lastName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			lastNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.last_name ILIKE ? AND office_users.status = 'REQUESTED'", lastNameSearch)
		}
	},
	"office": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			officeSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("transportation_offices.name ILIKE ? AND office_users.status = 'REQUESTED'", officeSearch)
		}
	},
	"requestedOn": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			trimAllZero, trimDayZero, trimMonthZero, noTrim := fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content)
			query.Where("(TO_CHAR(office_users.created_at, 'FMMM/FMDD/YYYY') ILIKE ? OR TO_CHAR(office_users.created_at, 'MM/FMDD/YYYY') ILIKE ? OR TO_CHAR(office_users.created_at, 'FMMM/DD/YYYY') ILIKE ? OR TO_CHAR(office_users.created_at, 'MM/DD/YYYY') ILIKE ?) AND office_users.status = 'REQUESTED'", trimAllZero, trimDayZero, trimMonthZero, noTrim)
		}
	},
	"roles": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			rolesSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("roles.role_name ILIKE ? AND office_users.status = 'REQUESTED'", rolesSearch)
		}
	},
}

// Handle retrieves a list of requested office users
func (h IndexRequestedOfficeUsersHandler) Handle(params requested_office_users.IndexRequestedOfficeUsersParams) middleware.Responder {
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
			for key, filterFunc := range requestedOfficeUserFilterConverters {
				if filterValue, exists := filtersMap[key]; exists {
					filterFuncs = append(filterFuncs, filterFunc(filterValue))
				}
			}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			officeUsers, count, err := h.RequestedOfficeUserListFetcher.FetchRequestedOfficeUsersList(appCtx, filterFuncs, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, s := range officeUsers {
				payload[i] = payloadForRequestedOfficeUserModel(s)
			}

			return requested_office_users.NewIndexRequestedOfficeUsersOK().WithContentRange(fmt.Sprintf("requested office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, count)).WithPayload(payload), nil
		})
}

// GetRequestedOfficeUserHandler returns a list of office users via GET /requested_office_users/{officeUserId}
type GetRequestedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserFetcher
	services.RoleFetcher
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

			roles, err := h.RoleFetcher.FetchRolesForUser(appCtx, *requestedOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return requested_office_users.NewGetRequestedOfficeUserInternalServerError(), err
			}

			requestedOfficeUser.User.Roles = roles

			payload := payloadForRequestedOfficeUserModel(requestedOfficeUser)

			return requested_office_users.NewGetRequestedOfficeUserOK().WithPayload(payload), nil
		})
}

// UpdateRequestedOfficeUserHandler updates a requested office user via PATCH /requested_office_users/{officeUserId}
type UpdateRequestedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.RequestedOfficeUserUpdater
	services.UserRoleAssociator
	services.RoleFetcher
}

// Handle updates a single requested office user
// this endpoint will be used when an admin is approving/rejecting the user without updates
// as well as approving/rejecting the user with updates
func (h UpdateRequestedOfficeUserHandler) Handle(params requested_office_users.UpdateRequestedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			requestedOfficeUserID, err := uuid.FromString(params.OfficeUserID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing error for %s", params.OfficeUserID.String()), zap.Error(err))
				return requested_office_users.NewUpdateRequestedOfficeUserUnprocessableEntity(), err
			}

			body := params.Body
			updatedRoles := rolesPayloadToModel(body.Roles)
			if len(updatedRoles) == 0 {
				err := apperror.NewBadDataError("No roles were matched from payload")
				appCtx.Logger().Error(err.Error())
				return requested_office_users.NewUpdateRequestedOfficeUserUnprocessableEntity(), err
			}

			var requestedOfficeUser *models.OfficeUser
			transactionError := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
				// handle Okta account creation only if the user is being approved
				// ignore this if we are in our dev environment
				if params.Body.Status == "APPROVED" && appCtx.Session().IDToken != "devlocal" {
					var err error
					_, err = fetchOrCreateOktaProfile(txAppCtx, params)
					if err != nil {
						txAppCtx.Logger().Error("error fetching/creating Okta profile", zap.Error(err))
						return fmt.Errorf("failed to create Okta account: %w", err)
					}
					txAppCtx.Logger().Info("Okta account successfully fetched or created")
				}

				var verrs *validate.Errors
				var err error
				requestedOfficeUser, verrs, err = h.RequestedOfficeUserUpdater.UpdateRequestedOfficeUser(txAppCtx, requestedOfficeUserID, params.Body)
				if err != nil {
					txAppCtx.Logger().Error("Error updating RequestedOfficeUser", zap.Error(err))
					return err
				}
				if verrs.HasAny() {
					txAppCtx.Logger().Error("Validation errors updating RequestedOfficeUser", zap.String("errors", verrs.Error()))
					return verrs
				}

				if requestedOfficeUser.UserID != nil && body.Roles != nil {
					_, verrs, err = h.UserRoleAssociator.UpdateUserRoles(txAppCtx, *requestedOfficeUser.UserID, updatedRoles)
					if verrs.HasAny() {
						txAppCtx.Logger().Error("Validation errors updating user roles", zap.String("errors", verrs.Error()))
						return verrs
					}
					if err != nil {
						txAppCtx.Logger().Error("Error updating user roles", zap.Error(err))
						return err
					}
				}

				roles, err := h.RoleFetcher.FetchRolesForUser(txAppCtx, *requestedOfficeUser.UserID)
				if err != nil {
					txAppCtx.Logger().Error("Error fetching user roles", zap.Error(err))
					return err
				}
				requestedOfficeUser.User.Roles = roles

				// send email notification if request was rejected
				if params.Body.Status == "REJECTED" {
					err = h.NotificationSender().SendNotification(txAppCtx, notifications.NewOfficeAccountRejected(requestedOfficeUser.ID))
					if err != nil {
						txAppCtx.Logger().Error("Error sending rejection email", zap.Error(err))
						return apperror.NewBadDataError("problem sending email to rejected office user")
					}
				}

				return nil
			})

			if transactionError != nil {
				appCtx.Logger().Error("Transaction error while updating requested office user", zap.Error(transactionError))
				return requested_office_users.NewUpdateRequestedOfficeUserInternalServerError(), transactionError
			}

			payload := payloadForRequestedOfficeUserModel(*requestedOfficeUser)

			return requested_office_users.NewUpdateRequestedOfficeUserOK().WithPayload(payload), nil
		})
}
