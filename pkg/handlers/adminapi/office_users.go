package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
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
		CreatedAt: *handlers.FmtDateTime(r.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(r.UpdatedAt),
	}
}

func payloadForPrivilege(p models.Privilege) *adminmessages.Privilege {
	return &adminmessages.Privilege{
		ID:            *handlers.FmtUUID(p.ID),
		PrivilegeType: *handlers.FmtString(string(p.PrivilegeType)),
		PrivilegeName: *handlers.FmtString(string(p.PrivilegeName)),
		CreatedAt:     *handlers.FmtDateTime(p.CreatedAt),
		UpdatedAt:     *handlers.FmtDateTime(p.UpdatedAt),
	}
}

func payloadForTransportationOfficeAssignment(toa models.TransportationOfficeAssignment) *adminmessages.TransportationOfficeAssignment {
	return &adminmessages.TransportationOfficeAssignment{
		OfficeUserID:           *handlers.FmtUUID(toa.ID),
		TransportationOfficeID: *handlers.FmtUUID(toa.TransportationOfficeID),
		PrimaryOffice:          *handlers.FmtBool(toa.PrimaryOffice),
		CreatedAt:              *handlers.FmtDateTime(toa.CreatedAt),
		UpdatedAt:              *handlers.FmtDateTime(toa.UpdatedAt),
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
	for _, role := range user.Roles {
		payload.Roles = append(payload.Roles, payloadForRole(role))
	}
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
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var officeUserFilterConverters = map[string]func(string) []services.QueryFilter{
	"search": func(content string) []services.QueryFilter {
		nameSearch := fmt.Sprintf("%s%%", content)
		return []services.QueryFilter{
			query.NewQueryFilter("email", "ILIKE", fmt.Sprintf("%%%s%%", content)),
			query.NewQueryFilter("first_name", "ILIKE", nameSearch),
			query.NewQueryFilter("last_name", "ILIKE", nameSearch),
		}
	},
}

// Handle retrieves a list of office users
func (h IndexOfficeUsersHandler) Handle(params officeuserop.IndexOfficeUsersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, officeUserFilterConverters)

			// Add a filter for approved status
			queryFilters = append(queryFilters, query.NewQueryFilter("status", "=", "APPROVED"))

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			queryAssociations := query.NewQueryAssociationsPreload([]services.QueryAssociation{
				query.NewQueryAssociation("User.Roles"),
				query.NewQueryAssociation("User.Privileges"),
			})

			var officeUsers models.OfficeUsers
			err := h.ListFetcher.FetchRecordList(appCtx, &officeUsers, queryFilters, queryAssociations, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficeUsersCount, err := h.ListFetcher.FetchRecordCount(appCtx, &officeUsers, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficeUsersCount := len(officeUsers)

			payload := make(adminmessages.OfficeUsers, queriedOfficeUsersCount)

			for i, s := range officeUsers {
				payload[i] = payloadForOfficeUserModel(s)
			}

			return officeuserop.NewIndexOfficeUsersOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficeUsersCount, totalOfficeUsersCount)).WithPayload(payload), nil
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
	services.RoleAssociater
	services.UserPrivilegeAssociator
	services.TransportaionOfficeAssignmentUpdater
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

			roles, err := h.RoleAssociater.FetchRolesForUser(appCtx, *createdOfficeUser.UserID)
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

			updatedTransportationOfficeAssignments, err := transportationOfficeAssignmentsPayloadToModel(payload.TransportationOfficeAssignments)
			if err != nil {
				appCtx.Logger().Error("UUID parsing error for transportation office assignments", zap.Error(err))
				return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
			}

			transportationOfficeAssignments, err :=
				h.TransportaionOfficeAssignmentUpdater.UpdateTransportaionOfficeAssignments(appCtx, createdOfficeUser.ID, updatedTransportationOfficeAssignments)
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
	services.TransportaionOfficeAssignmentUpdater
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

			updatedOfficeUser, verrs, err := h.OfficeUserUpdater.UpdateOfficeUser(appCtx, officeUserID, payload, primaryTransportationOfficeID)

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
			}

			if len(payload.TransportationOfficeAssignments) > 0 {

				transportationOfficeAssignmentsFromPayload, err := transportationOfficeAssignmentsPayloadToModel(payload.TransportationOfficeAssignments)
				if err != nil {
					appCtx.Logger().Error("UUID parsing error for transportation office assignments", zap.Error(err))
					return officeuserop.NewCreateOfficeUserUnprocessableEntity(), err
				}

				updatedTransportationOfficeAssignments, err :=
					h.TransportaionOfficeAssignmentUpdater.UpdateTransportaionOfficeAssignments(appCtx, updatedOfficeUser.ID, transportationOfficeAssignmentsFromPayload)
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

func privilegesPayloadToModel(payload []*adminmessages.OfficeUserPrivilege) []models.PrivilegeType {
	var rt []models.PrivilegeType
	for _, privilege := range payload {
		if privilege.PrivilegeType != nil {
			rt = append(rt, models.PrivilegeType(*privilege.PrivilegeType))
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
			PrimaryOffice:          *toa.PrimaryOffice,
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

	return transportationOfficeID, apperror.NewBadDataError("Could not identify primary transportaion office from list of assignments")
}
