package adminapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForRole(r roles.Role) *adminmessages.Role {
	roleType := string(r.RoleType)
	roleName := string(r.RoleName)
	sort := int32(r.Sort)
	return &adminmessages.Role{
		ID:        handlers.FmtUUID(r.ID),
		RoleType:  &roleType,
		RoleName:  &roleName,
		Sort:      sort,
		CreatedAt: *handlers.FmtDateTime(r.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(r.UpdatedAt),
	}
}

func payloadForPrivilege(p roles.Privilege) *adminmessages.Privilege {
	sort := int32(p.Sort)
	return &adminmessages.Privilege{
		ID:            *handlers.FmtUUID(p.ID),
		PrivilegeType: *handlers.FmtString(string(p.PrivilegeType)),
		PrivilegeName: *handlers.FmtString(string(p.PrivilegeName)),
		Sort:          sort,
		CreatedAt:     *handlers.FmtDateTime(p.CreatedAt),
		UpdatedAt:     *handlers.FmtDateTime(p.UpdatedAt),
	}
}

func payloadForRolePrivilege(role roles.Role) *adminmessages.Role {
	sort := int32(role.Sort)
	r := &adminmessages.Role{
		ID:        handlers.FmtUUID(role.ID),
		RoleType:  handlers.FmtString(string(role.RoleType)),
		RoleName:  handlers.FmtString(string(role.RoleName)),
		Sort:      sort,
		CreatedAt: *handlers.FmtDateTime(role.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(role.UpdatedAt),
	}

	for _, rp := range role.RolePrivileges {
		privSort := int32(rp.Privilege.Sort)
		r.Privileges = append(r.Privileges, &adminmessages.Privilege{
			ID:            *handlers.FmtUUID(rp.PrivilegeID),
			PrivilegeType: *handlers.FmtString(string(rp.Privilege.PrivilegeType)),
			PrivilegeName: *handlers.FmtString(string(rp.Privilege.PrivilegeName)),
			Sort:          privSort,
			CreatedAt:     *handlers.FmtDateTime(rp.Privilege.CreatedAt),
			UpdatedAt:     *handlers.FmtDateTime(rp.Privilege.UpdatedAt),
		})
	}
	return r
}

func payloadForTransportationOfficeAssignment(toa models.TransportationOfficeAssignment) *adminmessages.TransportationOfficeAssignment {
	return &adminmessages.TransportationOfficeAssignment{
		OfficeUserID:           *handlers.FmtUUID(toa.ID),
		TransportationOfficeID: *handlers.FmtUUID(toa.TransportationOfficeID),
		PrimaryOffice:          *handlers.FmtBool(*toa.PrimaryOffice),
		CreatedAt:              *handlers.FmtDateTime(toa.CreatedAt),
		UpdatedAt:              *handlers.FmtDateTime(toa.UpdatedAt),
	}
}

// Ensures the payload does not have duplicate roles in the roles array. That would cause the Admin UI to show duplicate roles for a user.
func nonDuplicateRolesList(roles roles.Roles) []*adminmessages.Role {
	var rolesList []*adminmessages.Role
	seenRoles := make(map[string]bool)

	for _, role := range roles {
		roleName := string(role.RoleName)

		if !seenRoles[roleName] {
			rolesList = append(rolesList, payloadForRole(role))
			seenRoles[roleName] = true
		}
	}

	return rolesList
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
		Status:                 (*string)(o.Status),
		CreatedAt:              *handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:              *handlers.FmtDateTime(o.UpdatedAt),
	}
	if o.UserID != nil {
		userIDFmt := handlers.FmtUUID(*o.UserID)
		if userIDFmt != nil {
			payload.UserID = *userIDFmt
		}
	}

	payload.Roles = nonDuplicateRolesList(user.Roles)

	for _, privilege := range user.Privileges {
		payload.Privileges = append(payload.Privileges, payloadForPrivilege(privilege))
	}
	for _, transportationAssignment := range o.TransportationOfficeAssignments {
		payload.TransportationOfficeAssignments = append(payload.TransportationOfficeAssignments, payloadForTransportationOfficeAssignment(transportationAssignment))
	}
	return payload
}

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexOfficeUsersHandler struct {
	handlers.HandlerConfig
	services.OfficeUserListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var officeUserFilterConverters = map[string]func(string) func(*pop.Query){
	"search": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstSearch, lastSearch, emailSearch := fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content), fmt.Sprintf("%%%s%%", content)
			query.Where("(office_users.first_name ILIKE ? OR office_users.last_name ILIKE ? OR office_users.email ILIKE ?) AND office_users.status = 'APPROVED'", firstSearch, lastSearch, emailSearch)
		}
	},
	"email": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			emailSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.email ILIKE ? AND office_users.status = 'APPROVED'", emailSearch)
		}
	},
	"phone": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			phoneSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.telephone ILIKE ? AND office_users.status = 'APPROVED'", phoneSearch)
		}
	},
	"firstName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			firstNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.first_name ILIKE ? AND office_users.status = 'APPROVED'", firstNameSearch)
		}
	},
	"lastName": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			lastNameSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("office_users.last_name ILIKE ? AND office_users.status = 'APPROVED'", lastNameSearch)
		}
	},
	"office": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			officeSearch := fmt.Sprintf("%%%s%%", content)
			query.Where("transportation_offices.name ILIKE ? AND office_users.status = 'APPROVED'", officeSearch)
		}
	},
	"active": func(content string) func(*pop.Query) {
		return func(query *pop.Query) {
			query.Where("office_users.active = ? AND office_users.status = 'APPROVED'", content)
		}
	},
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
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
			for key, filterFunc := range officeUserFilterConverters {
				if filterValue, exists := filtersMap[key]; exists {
					filterFuncs = append(filterFuncs, filterFunc(filterValue))
				}
			}

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			officeUsers, count, err := h.OfficeUserListFetcher.FetchOfficeUsersList(appCtx, filterFuncs, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, s := range officeUsers {
				payload[i] = payloadForOfficeUserModel(s)
			}

			return officeuserop.NewIndexOfficeUsersOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, count)).WithPayload(payload), nil
		})
}

// GetOfficeUserHandler retrieves office user handler
type GetOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserFetcherPop
	services.NewQueryFilter
}

// Handle retrieves an office user
func (h GetOfficeUserHandler) Handle(params officeuserop.GetOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			officeUserID := uuid.FromStringOrNil(params.OfficeUserID.String())
			officeUser, err := h.OfficeUserFetcherPop.FetchOfficeUserByID(appCtx, officeUserID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			userError := appCtx.DB().Load(&officeUser, "User")
			if userError != nil {
				return handlers.ResponseForError(appCtx.Logger(), userError), userError
			}

			//todo: we want to move this query out of the handler and into querybuilder, if possible
			roleError := appCtx.DB().Q().Join("users_roles", "users_roles.role_id = roles.id").
				Where("users_roles.deleted_at IS NULL AND users_roles.user_id = ?", (officeUser.User.ID)).
				All(&officeUser.User.Roles)
			if roleError != nil {
				return handlers.ResponseForError(appCtx.Logger(), roleError), roleError
			}

			privilegeError := appCtx.DB().Q().Join("users_privileges", "users_privileges.privilege_id = privileges.id").
				Where("users_privileges.deleted_at IS NULL AND users_privileges.user_id = ?", (officeUser.User.ID)).
				All(&officeUser.User.Privileges)
			if privilegeError != nil {
				return handlers.ResponseForError(appCtx.Logger(), privilegeError), privilegeError
			}

			transportationOfficeAssignmentError := appCtx.DB().Q().EagerPreload("TransportationOffice").
				Join("transportation_offices", "transportation_office_assignments.transportation_office_id = transportation_offices.id").
				Where("transportation_office_assignments.id = ?", (officeUser.ID)).
				All(&officeUser.TransportationOfficeAssignments)
			if transportationOfficeAssignmentError != nil {
				return handlers.ResponseForError(appCtx.Logger(), transportationOfficeAssignmentError), transportationOfficeAssignmentError
			}

			payload := payloadForOfficeUserModel(officeUser)

			return officeuserop.NewGetOfficeUserOK().WithPayload(payload), nil
		})
}

// CreateOfficeUserHandler creates an office user
type CreateOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserCreator
	services.NewQueryFilter
	services.UserRoleAssociator
	services.RoleFetcher
	services.UserPrivilegeAssociator
	services.TransportationOfficeAssignmentUpdater
}

// Handle creates an office user
func (h CreateOfficeUserHandler) Handle(params officeuserop.CreateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if len(payload.TransportationOfficeAssignments) == 0 {
				err := apperror.NewBadDataError("At least one transportation office is required")
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			primaryTransportationOfficeID, err := getPrimaryTransportationOfficeIDFromPayload(payload.TransportationOfficeAssignments)

			if err != nil {
				appCtx.Logger().Error("Error identifying primary transportation office", zap.Error(err))
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			if len(payload.Roles) == 0 {
				err = apperror.NewBadDataError("At least one office user role is required")
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			updatedRoles := rolesPayloadToModel(payload.Roles)
			if len(updatedRoles) == 0 {
				err = apperror.NewBadDataError("No roles were matched from payload")
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			verrs, err := h.UserPrivilegeAssociator.VerifyUserPrivilegeAllowed(appCtx, payload.Roles, payload.Privileges)

			if err != nil {
				appCtx.Logger().Error("Error verifying user privileges allowed", zap.Error(err))
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			if verrs.HasAny() {
				validationError := &adminmessages.ValidationError{
					InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors, ClientError: adminmessages.ClientError{
						Title:    handlers.FmtString(handlers.ValidationErrMessage),
						Detail:   handlers.FmtString("Selected office user role is not authorized for supplied privilege"),
						Instance: handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest)),
					},
				}

				return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
			}

			// if the user is being manually created, then we know they will already be approved
			officeUserStatus := models.OfficeUserStatusAPPROVED

			officeUser := models.OfficeUser{
				LastName:               payload.LastName,
				FirstName:              payload.FirstName,
				Telephone:              payload.Telephone,
				Email:                  payload.Email,
				TransportationOfficeID: primaryTransportationOfficeID,
				Active:                 true,
				Status:                 &officeUserStatus,
			}

			primaryTransportationIDFilter := []services.QueryFilter{
				h.NewQueryFilter("id", "=", primaryTransportationOfficeID),
			}

			createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(appCtx, &officeUser, primaryTransportationIDFilter)
			if verrs != nil {
				validationError := &adminmessages.ValidationError{
					InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
				}

				validationError.Title = handlers.FmtString(handlers.ValidationErrMessage)
				validationError.Detail = handlers.FmtString("The information you provided is invalid.")
				validationError.Instance = handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest))

				return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
			}

			if err != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err))
				return officeuserop.NewCreateOfficeUserInternalServerError(), err
			}

			_, verrs, err = h.UserRoleAssociator.UpdateUserRoles(appCtx, *createdOfficeUser.UserID, updatedRoles)
			if verrs.HasAny() {
				validationError := &adminmessages.ValidationError{
					InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
				}

				validationError.Title = handlers.FmtString(handlers.ValidationErrMessage)
				validationError.Detail = handlers.FmtString("The information you provided is invalid.")
				validationError.Instance = handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest))

				return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
			}
			if err != nil {
				appCtx.Logger().Error("Error updating user roles", zap.Error(err))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			roles, err := h.RoleFetcher.FetchRolesForUser(appCtx, *createdOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			createdOfficeUser.User.Roles = roles

			updatedPrivileges := privilegesPayloadToModel(payload.Privileges)
			_, err = h.UserPrivilegeAssociator.UpdateUserPrivileges(appCtx, *createdOfficeUser.UserID, updatedPrivileges)
			if err != nil {
				appCtx.Logger().Error("Error updating user privileges", zap.Error(err))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			privileges, err := h.UserPrivilegeAssociator.FetchPrivilegesForUser(appCtx, *createdOfficeUser.UserID)

			if err != nil {
				appCtx.Logger().Error("Error fetching user privileges", zap.Error(err))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			createdOfficeUser.User.Privileges = privileges

			updatedTransportationOfficeAssignments, err := transportationOfficeAssignmentsPayloadToModel(payload.TransportationOfficeAssignments)
			if err != nil {
				appCtx.Logger().Error("UUID parsing error for transportation office assignments", zap.Error(err))
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			transportationOfficeAssignments, err :=
				h.TransportationOfficeAssignmentUpdater.UpdateTransportationOfficeAssignments(appCtx, createdOfficeUser.ID, updatedTransportationOfficeAssignments)
			if err != nil {
				appCtx.Logger().Error("Error updating office user's transportation office assignments", zap.Error(err))
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			createdOfficeUser.TransportationOfficeAssignments = transportationOfficeAssignments

			_, err = audit.Capture(appCtx, createdOfficeUser, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
			return officeuserop.NewCreateOfficeUserCreated().WithPayload(returnPayload), nil
		})
}

// UpdateOfficeUserHandler updates an office user
type UpdateOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserUpdater
	services.NewQueryFilter
	services.UserRoleAssociator
	services.UserPrivilegeAssociator
	services.UserSessionRevocation
	services.TransportationOfficeAssignmentUpdater
	services.RoleFetcher
}

// Handle updates an office user
func (h UpdateOfficeUserHandler) Handle(params officeuserop.UpdateOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			officeUserID, err := uuid.FromString(params.OfficeUserID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", params.OfficeUserID.String()), zap.Error(err))
			}

			var primaryTransportationOfficeID uuid.UUID
			if len(payload.TransportationOfficeAssignments) > 0 {
				primaryTransportationOfficeID, err = getPrimaryTransportationOfficeIDFromPayload(payload.TransportationOfficeAssignments)

				if err != nil {
					appCtx.Logger().Error("Error identifying primary transportation office", zap.Error(err))
					return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
				}
			}

			verrs, err := h.UserPrivilegeAssociator.VerifyUserPrivilegeAllowed(appCtx, payload.Roles, payload.Privileges)

			if err != nil {
				appCtx.Logger().Error("Error verifying user privileges allowed", zap.Error(err))
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			if verrs.HasAny() {
				validationError := &adminmessages.ValidationError{
					InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors, ClientError: adminmessages.ClientError{
						Title:    handlers.FmtString(handlers.ValidationErrMessage),
						Detail:   handlers.FmtString("Selected office user role is not authorized for supplied privilege"),
						Instance: handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest)),
					},
				}

				return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
			}
			officeUserDB, err := models.FetchOfficeUserByID(appCtx.DB(), officeUserID)

			if err != nil {
				appCtx.Logger().Error("Error fetching office user", zap.Error(err))
				return officeuserop.NewUpdateOfficeUserNotFound(), err
			}

			newOfficeUser := payloads.OfficeUserModelFromUpdate(payload, officeUserDB)

			updatedOfficeUser, verrs, err := h.OfficeUserUpdater.UpdateOfficeUser(appCtx, officeUserID, newOfficeUser, primaryTransportationOfficeID)

			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err), zap.Error(verrs))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			if updatedOfficeUser.UserID != nil && payload.Roles != nil {
				updatedRoles := rolesPayloadToModel(payload.Roles)
				_, verrs, err = h.UserRoleAssociator.UpdateUserRoles(appCtx, *updatedOfficeUser.UserID, updatedRoles)
				if verrs.HasAny() {
					validationError := &adminmessages.ValidationError{
						InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
					}

					validationError.Title = handlers.FmtString(handlers.ValidationErrMessage)
					validationError.Detail = handlers.FmtString("The information you provided is invalid.")
					validationError.Instance = handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest))

					return officeuserop.NewCreateOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
				}
				if err != nil {
					appCtx.Logger().Error("Error updating user roles", zap.Error(err))
					return officeuserop.NewUpdateOfficeUserInternalServerError(), err
				}

				boolean := true
				revokeOfficeSessionPayload := adminmessages.UserUpdate{
					RevokeOfficeSession: &boolean,
				}

				_, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(
					appCtx,
					*updatedOfficeUser.UserID,
					&revokeOfficeSessionPayload,
					h.SessionManagers(),
				)

				if revokeErr != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(revokeErr))
					return userop.NewUpdateUserInternalServerError(), revokeErr
				}

				if validationErrors != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(verrs))
					return userop.NewUpdateUserInternalServerError(), validationErrors
				}

				roles, err := h.RoleFetcher.FetchRolesForUser(appCtx, *updatedOfficeUser.UserID)

				if err != nil {
					appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
					return officeuserop.NewUpdateOfficeUserInternalServerError(), err
				}

				updatedOfficeUser.User.Roles = roles
			}

			if updatedOfficeUser.UserID != nil && payload.Privileges != nil {
				updatedPrivileges := privilegesPayloadToModel(payload.Privileges)
				_, err = h.UserPrivilegeAssociator.UpdateUserPrivileges(appCtx, *updatedOfficeUser.UserID, updatedPrivileges)
				if err != nil {
					appCtx.Logger().Error("Error updating user privileges", zap.Error(err))
					return officeuserop.NewUpdateOfficeUserInternalServerError(), err
				}

				boolean := true
				revokeOfficeSessionPayload := adminmessages.UserUpdate{
					RevokeOfficeSession: &boolean,
				}

				_, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(
					appCtx,
					*updatedOfficeUser.UserID,
					&revokeOfficeSessionPayload,
					h.SessionManagers(),
				)

				if revokeErr != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(revokeErr))
					return userop.NewUpdateUserInternalServerError(), revokeErr
				}

				if validationErrors != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(verrs))
					return userop.NewUpdateUserInternalServerError(), validationErrors
				}

				privileges, err := h.UserPrivilegeAssociator.FetchPrivilegesForUser(appCtx, *updatedOfficeUser.UserID)

				if err != nil {
					appCtx.Logger().Error("Error fetching user privileges", zap.Error(err))
					return officeuserop.NewUpdateOfficeUserInternalServerError(), err
				}

				updatedOfficeUser.User.Privileges = privileges
			}

			if len(payload.TransportationOfficeAssignments) > 0 {

				transportationOfficeAssignmentsFromPayload, err := transportationOfficeAssignmentsPayloadToModel(payload.TransportationOfficeAssignments)
				if err != nil {
					appCtx.Logger().Error("UUID parsing error for transportation office assignments", zap.Error(err))
					return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
				}

				updatedTransportationOfficeAssignments, err :=
					h.TransportationOfficeAssignmentUpdater.UpdateTransportationOfficeAssignments(appCtx, updatedOfficeUser.ID, transportationOfficeAssignmentsFromPayload)
				if err != nil {
					appCtx.Logger().Error("Error updating office user's transportation office assignments", zap.Error(err))
					return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
				}

				updatedOfficeUser.TransportationOfficeAssignments = updatedTransportationOfficeAssignments

				boolean := true
				revokeOfficeSessionPayload := adminmessages.UserUpdate{
					RevokeOfficeSession: &boolean,
				}

				_, validationErrors, revokeErr := h.UserSessionRevocation.RevokeUserSession(
					appCtx,
					*updatedOfficeUser.UserID,
					&revokeOfficeSessionPayload,
					h.SessionManagers(),
				)

				if revokeErr != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(revokeErr))
					return userop.NewUpdateUserInternalServerError(), revokeErr
				}

				if validationErrors != nil {
					err = apperror.NewInternalServerError("Error revoking user session")
					appCtx.Logger().Error(err.Error(), zap.Error(verrs))
					return userop.NewUpdateUserInternalServerError(), validationErrors
				}
			}

			// Log if the account was enabled or disabled (POAM requirement)
			if payload.Active != nil {
				_, err = audit.CaptureAccountStatus(appCtx, updatedOfficeUser, *payload.Active, params.HTTPRequest)
				if err != nil {
					appCtx.Logger().Error("Error capturing account status audit record in UpdateOfficeUserHandler", zap.Error(err))
				}
			}

			_, err = audit.Capture(appCtx, updatedOfficeUser, payload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForOfficeUserModel(*updatedOfficeUser)

			return officeuserop.NewUpdateOfficeUserOK().WithPayload(returnPayload), nil
		})
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

func privilegesPayloadToModel(payload []*adminmessages.OfficeUserPrivilege) []roles.PrivilegeType {
	var rt []roles.PrivilegeType
	for _, privilege := range payload {
		if privilege.PrivilegeType != nil {
			rt = append(rt, roles.PrivilegeType(*privilege.PrivilegeType))
		}
	}
	return rt
}

func transportationOfficeAssignmentsPayloadToModel(payload []*adminmessages.OfficeUserTransportationOfficeAssignment) (models.TransportationOfficeAssignments, error) {
	var toas models.TransportationOfficeAssignments
	for _, toa := range payload {
		transportationOfficeID, err := uuid.FromString(toa.TransportationOfficeID.String())

		if err != nil {
			return models.TransportationOfficeAssignments{}, err
		}

		model := &models.TransportationOfficeAssignment{
			TransportationOfficeID: transportationOfficeID,
			PrimaryOffice:          toa.PrimaryOffice,
		}

		toas = append(toas, *model)
	}
	return toas, nil
}

func getPrimaryTransportationOfficeIDFromPayload(payload []*adminmessages.OfficeUserTransportationOfficeAssignment) (uuid.UUID, error) {
	var transportationOfficeID uuid.UUID
	var err error

	if len(payload) == 1 {
		transportationOfficeID, err = uuid.FromString(payload[0].TransportationOfficeID.String())
		return transportationOfficeID, err
	}

	for _, toa := range payload {
		if toa.PrimaryOffice != nil && *toa.PrimaryOffice {
			transportationOfficeID, err = uuid.FromString(toa.TransportationOfficeID.String())
			return transportationOfficeID, err
		}
	}

	return transportationOfficeID, apperror.NewBadDataError("Could not identify primary transportation office from list of assignments")
}

// DeleteOfficeUserHandler deletes an office user via DELETE /office_user/{officeUserId}
type DeleteOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserDeleter
}

func (h DeleteOfficeUserHandler) Handle(params officeuserop.DeleteOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// we only allow this to be called from the admin app
			if !appCtx.Session().IsAdminApp() {
				return officeuserop.NewDeleteOfficeUserUnauthorized(), nil
			}

			officeUserID, err := uuid.FromString(params.OfficeUserID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			err = h.OfficeUserDeleter.DeleteOfficeUser(appCtx, officeUserID)
			if err != nil {
				switch err.(type) {
				case apperror.NotFoundError:
					return officeuserop.NewDeleteOfficeUserNotFound(), err
				default:
					return officeuserop.NewDeleteOfficeUserInternalServerError(), err
				}
			}

			return officeuserop.NewDeleteOfficeUserNoContent(), nil
		})
}

// GetRolesPrivilegesHandler retrieves a list of unique role to privilege mappings via GET /office_users/roles-privileges
type GetRolesPrivilegesHandler struct {
	handlers.HandlerConfig
	services.RoleFetcher
}

func (h GetRolesPrivilegesHandler) Handle(params officeuserop.GetRolesPrivilegesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// we only allow this to be called from the admin app
			if !appCtx.Session().IsAdminApp() {
				return officeuserop.NewGetRolesPrivilegesUnauthorized(), nil
			}

			rolesWithRolePrivs, err := h.RoleFetcher.FetchRolesPrivileges(appCtx)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return officeuserop.NewGetRolesPrivilegesNotFound(), err
			} else if err != nil {
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewGetRolesPrivilegesInternalServerError(), err
			}

			payload := make([]*adminmessages.Role, len(rolesWithRolePrivs))
			for i, rwrp := range rolesWithRolePrivs {
				payload[i] = payloadForRolePrivilege(rwrp)
			}

			return officeuserop.NewGetRolesPrivilegesOK().WithPayload(payload), nil
		})
}
