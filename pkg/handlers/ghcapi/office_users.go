package ghcapi

import (
	"database/sql"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	officeuserop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/office_users"
	rpop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/role_privileges"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
)

// RequestOfficeUserHandler allows for the creation of "requested" status office users
type RequestOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserCreator
	services.NewQueryFilter
	services.UserRoleAssociator
	services.RoleFetcher
	services.UserPrivilegeAssociator
	services.PrivilegeFetcher
	services.TransportationOfficeAssignmentUpdater
}

// Convert internal role model to ghc role model
func payloadForRole(r roles.Role) *ghcmessages.Role {
	roleType := string(r.RoleType)
	roleName := string(r.RoleName)
	sort := int32(r.Sort)
	return &ghcmessages.Role{
		ID:        handlers.FmtUUID(r.ID),
		RoleType:  &roleType,
		RoleName:  &roleName,
		Sort:      sort,
		CreatedAt: *handlers.FmtDateTime(r.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(r.UpdatedAt),
	}
}

// Convert internal privilege model to ghc privilege model
func payloadForPrivilege(p roles.Privilege) *ghcmessages.Privilege {
	privilegeType := string(p.PrivilegeType)
	privilegeName := string(p.PrivilegeName)
	sort := int32(p.Sort)
	return &ghcmessages.Privilege{
		ID:            handlers.FmtUUID(p.ID),
		PrivilegeType: &privilegeType,
		PrivilegeName: &privilegeName,
		Sort:          sort,
		CreatedAt:     *handlers.FmtDateTime(p.CreatedAt),
		UpdatedAt:     *handlers.FmtDateTime(p.UpdatedAt),
	}
}

// Convert ghc role models to internal role models
func rolesPayloadToModel(payload []*ghcmessages.OfficeUserRole) []roles.RoleType {
	var rt []roles.RoleType
	for _, role := range payload {
		if role.RoleType != nil {
			rt = append(rt, roles.RoleType(*role.RoleType))
		}
	}
	return rt
}

// Convert ghc privilege models to internal privilege models
func privilegesPayloadToModel(payload []*ghcmessages.OfficeUserPrivilege) []roles.PrivilegeType {
	var pt []roles.PrivilegeType
	for _, role := range payload {
		if role.PrivilegeType != nil {
			pt = append(pt, roles.PrivilegeType(*role.PrivilegeType))
		}
	}
	return pt
}

func payloadForRolePrivilege(role roles.Role) *ghcmessages.Role {
	sort := int32(role.Sort)
	r := &ghcmessages.Role{
		ID:        handlers.FmtUUID(role.ID),
		RoleType:  handlers.FmtString(string(role.RoleType)),
		RoleName:  handlers.FmtString(string(role.RoleName)),
		Sort:      sort,
		CreatedAt: *handlers.FmtDateTime(role.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(role.UpdatedAt),
	}
	for _, rp := range role.RolePrivileges {
		privSort := int32(rp.Privilege.Sort)
		r.Privileges = append(r.Privileges, &ghcmessages.Privilege{
			ID:            handlers.FmtUUID(rp.PrivilegeID),
			PrivilegeType: handlers.FmtString(string(rp.Privilege.PrivilegeType)),
			PrivilegeName: handlers.FmtString(string(rp.Privilege.PrivilegeName)),
			Sort:          privSort,
			CreatedAt:     *handlers.FmtDateTime(rp.Privilege.CreatedAt),
			UpdatedAt:     *handlers.FmtDateTime(rp.Privilege.UpdatedAt),
		})
	}
	return r
}

// Convert internal office user model to ghc office user model
func payloadForOfficeUserModel(o models.OfficeUser) *ghcmessages.OfficeUser {
	var user models.User
	if o.UserID != nil {
		user = o.User
	}
	payload := &ghcmessages.OfficeUser{
		ID:                              handlers.FmtUUID(o.ID),
		FirstName:                       handlers.FmtString(o.FirstName),
		MiddleInitials:                  handlers.FmtStringPtr(o.MiddleInitials),
		LastName:                        handlers.FmtString(o.LastName),
		Telephone:                       handlers.FmtString(o.Telephone),
		Email:                           handlers.FmtString(o.Email),
		Edipi:                           handlers.FmtStringPtr(o.EDIPI),
		OtherUniqueID:                   handlers.FmtStringPtr(o.OtherUniqueID),
		TransportationOfficeID:          handlers.FmtUUID(o.TransportationOfficeID),
		TransportationOfficeAssignments: payloadForTransportationOfficeAssignments(o.TransportationOfficeAssignments),
		Active:                          handlers.FmtBool(o.Active),
		Status:                          (*string)(o.Status),
		CreatedAt:                       *handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:                       *handlers.FmtDateTime(o.UpdatedAt),
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
	return payload
}

func payloadForTransportationOfficeAssignments(toas models.TransportationOfficeAssignments) []*ghcmessages.TransportationOfficeAssignment {
	var payload []*ghcmessages.TransportationOfficeAssignment
	for _, toa := range toas {
		payload = append(payload, &ghcmessages.TransportationOfficeAssignment{
			OfficeUserID:           handlers.FmtUUID(toa.ID),
			TransportationOfficeID: handlers.FmtUUID(toa.TransportationOfficeID),
			PrimaryOffice:          handlers.FmtBool(*toa.PrimaryOffice),
		})
	}
	return payload
}

// Handle creates the office user with a status of requested
func (h RequestOfficeUserHandler) Handle(params officeuserop.CreateRequestedOfficeUserParams) middleware.Responder {
	payload := params.OfficeUser
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			transportationOfficeID, err := uuid.FromString(payload.TransportationOfficeID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", payload.TransportationOfficeID.String()), zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity(), err
			}

			if len(payload.Roles) == 0 {
				err = apperror.NewBadDataError("At least one office user role is required")
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity(), err
			}

			updatedRoles := rolesPayloadToModel(payload.Roles)
			if len(updatedRoles) == 0 {
				err = apperror.NewBadDataError("No roles were matched from payload")
				appCtx.Logger().Error(err.Error())
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity(), err
			}

			updatedPrivileges := privilegesPayloadToModel(payload.Privileges)

			// Enforce identification rule for this payload
			if payload.Edipi == nil && payload.OtherUniqueID == nil {
				err = apperror.NewBadDataError("Either an EDIPI or Other Unique ID must be provided")
				appCtx.Logger().Error(err.Error())
				payload := payloadForValidationError(
					"Identification parameter error",
					err.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity().WithPayload(payload), err
			}

			status := models.OfficeUserStatusREQUESTED
			// By default set status to "REQUESTED", as is the purpose of this endpoint
			officeUser := models.OfficeUser{
				LastName:               payload.LastName,
				FirstName:              payload.FirstName,
				MiddleInitials:         payload.MiddleInitials,
				Telephone:              payload.Telephone,
				Email:                  payload.Email,
				EDIPI:                  payload.Edipi,
				OtherUniqueID:          payload.OtherUniqueID,
				TransportationOfficeID: transportationOfficeID,
				Active:                 false,
				Status:                 &status,
			}

			transportationIDFilter := []services.QueryFilter{
				h.NewQueryFilter("id", "=", transportationOfficeID),
			}

			createdOfficeUser, verrs, err := h.OfficeUserCreator.CreateOfficeUser(appCtx, &officeUser, transportationIDFilter)
			if verrs != nil && len(verrs.Errors) > 0 {
				payload := payloadForValidationError(
					"Office user creation",
					"Validation error",
					h.GetTraceIDFromRequest(params.HTTPRequest),
					verrs)
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity().WithPayload(payload), verrs
			}

			if err != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			_, verrs, err = h.UserRoleAssociator.UpdateUserRoles(appCtx, *createdOfficeUser.UserID, updatedRoles)
			if verrs.HasAny() {
				validationError := &ghcmessages.ValidationError{
					InvalidFields: handlers.NewValidationErrorsResponse(verrs).Errors,
				}

				validationError.Title = handlers.FmtString(handlers.ValidationErrMessage)
				validationError.Detail = handlers.FmtString("The information you provided is invalid.")
				validationError.Instance = handlers.FmtUUID(h.GetTraceIDFromRequest(params.HTTPRequest))

				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity().WithPayload(validationError), verrs
			}
			if err != nil {
				appCtx.Logger().Error("Error updating user roles", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			_, err = h.UserPrivilegeAssociator.UpdateUserPrivileges(appCtx, *createdOfficeUser.UserID, updatedPrivileges)

			if err != nil {
				appCtx.Logger().Error("Error updating user privileges", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			roles, err := h.RoleFetcher.FetchRolesForUser(appCtx, *createdOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			privileges, err := h.UserPrivilegeAssociator.FetchPrivilegesForUser(appCtx, *createdOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user privileges", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			createdOfficeUser.User.Roles = roles
			createdOfficeUser.User.Privileges = privileges

			transportationOfficeAssignments := models.TransportationOfficeAssignments{
				{
					TransportationOfficeID: transportationOfficeID,
					PrimaryOffice:          models.BoolPointer(true),
				},
			}

			createdTransportationOfficeAssignments, err :=
				h.TransportationOfficeAssignmentUpdater.UpdateTransportationOfficeAssignments(appCtx, createdOfficeUser.ID, transportationOfficeAssignments)
			if err != nil {
				appCtx.Logger().Error("Error updating office user's transportation office assignments", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserUnprocessableEntity(), err
			}

			createdOfficeUser.TransportationOfficeAssignments = createdTransportationOfficeAssignments

			_, err = audit.Capture(appCtx, createdOfficeUser, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForOfficeUserModel(*createdOfficeUser)
			return officeuserop.NewCreateRequestedOfficeUserCreated().WithPayload(returnPayload), nil
		})
}

// UpdateOfficeUserHandler updates an office user
type UpdateOfficeUserHandler struct {
	handlers.HandlerConfig
	services.OfficeUserUpdater
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

			if officeUserID != appCtx.Session().OfficeUserID {
				err := apperror.NewForbiddenError("Office User ID does not match session office user ID")
				appCtx.Logger().Error(err.Error(), zap.Error(err))
				return officeuserop.NewUpdateOfficeUserUnauthorized(), err
			}

			officeUserDB, err := models.FetchOfficeUserByID(appCtx.DB(), officeUserID)

			if err != nil {
				appCtx.Logger().Error(err.Error(), zap.Error(err))
				return officeuserop.NewUpdateOfficeUserNotFound(), err
			}

			newOfficeUser := payloads.OfficeUserModelFromUpdate(payload, officeUserDB)

			updatedOfficeUser, verrs, err := h.OfficeUserUpdater.UpdateOfficeUser(appCtx, officeUserID, newOfficeUser, uuid.Nil)

			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err), zap.Error(verrs))
				return officeuserop.NewUpdateOfficeUserInternalServerError(), err
			}

			_, err = audit.Capture(appCtx, updatedOfficeUser, payload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForOfficeUserModel(*updatedOfficeUser)

			return officeuserop.NewUpdateOfficeUserOK().WithPayload(returnPayload), nil
		})
}

// GetRolesPrivilegesHandler retrieves a list of unique role to privilege mappings via GET /office_users/roles-privileges
type GetRolesPrivilegesHandler struct {
	handlers.HandlerConfig
	services.RoleFetcher
}

func (h GetRolesPrivilegesHandler) Handle(params rpop.GetRolesPrivilegesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			rolesWithRolePrivs, err := h.RoleFetcher.FetchRolesPrivileges(appCtx)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return rpop.NewGetRolesPrivilegesNotFound(), err
			} else if err != nil {
				appCtx.Logger().Error(err.Error())
				return rpop.NewGetRolesPrivilegesInternalServerError(), err
			}

			payload := make([]*ghcmessages.Role, len(rolesWithRolePrivs))
			for i, rwrp := range rolesWithRolePrivs {
				payload[i] = payloadForRolePrivilege(rwrp)
			}

			return rpop.NewGetRolesPrivilegesOK().WithPayload(payload), nil
		})
}
