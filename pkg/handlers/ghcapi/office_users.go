package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	officeuserop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
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
	services.RoleAssociater
	services.TransportaionOfficeAssignmentUpdater
}

// Convert internal role model to ghc role model
func payloadForRole(r roles.Role) *ghcmessages.Role {
	roleType := string(r.RoleType)
	roleName := string(r.RoleName)
	return &ghcmessages.Role{
		ID:        handlers.FmtUUID(r.ID),
		RoleType:  &roleType,
		RoleName:  &roleName,
		CreatedAt: *handlers.FmtDateTime(r.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(r.UpdatedAt),
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

// Convert internal office user model to ghc office user model
func payloadForOfficeUserModel(o models.OfficeUser) *ghcmessages.OfficeUser {
	var user models.User
	if o.UserID != nil {
		user = o.User
	}
	payload := &ghcmessages.OfficeUser{
		ID:                     handlers.FmtUUID(o.ID),
		FirstName:              handlers.FmtString(o.FirstName),
		MiddleInitials:         handlers.FmtStringPtr(o.MiddleInitials),
		LastName:               handlers.FmtString(o.LastName),
		Telephone:              handlers.FmtString(o.Telephone),
		Email:                  handlers.FmtString(o.Email),
		Edipi:                  handlers.FmtStringPtr(o.EDIPI),
		OtherUniqueID:          handlers.FmtStringPtr(o.OtherUniqueID),
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

			roles, err := h.RoleAssociater.FetchRolesForUser(appCtx, *createdOfficeUser.UserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching user roles", zap.Error(err))
				return officeuserop.NewCreateRequestedOfficeUserInternalServerError(), err
			}

			createdOfficeUser.User.Roles = roles

			transportationOfficeAssignments := models.TransportationOfficeAssignments{
				{
					TransportationOfficeID: transportationOfficeID,
					PrimaryOffice:          true,
				},
			}

			createdTransportationOfficeAssignments, err :=
				h.TransportaionOfficeAssignmentUpdater.UpdateTransportaionOfficeAssignments(appCtx, createdOfficeUser.ID, transportationOfficeAssignments)
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
