package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler struct {
	handlers.HandlerConfig
	officeUserFetcherPop services.OfficeUserFetcherPop
}

// decoratePayloadWithCurrentAndInactiveRoles will add session role to the logged in user payload and return it
func decoratePayloadWithCurrentAndInactiveRoles(s *auth.Session, p *internalmessages.LoggedInUserPayload, allUserRoles roles.Roles) {
	if s == nil || p == nil {
		return
	}
	p.ActiveRole = &internalmessages.Role{
		ID:        handlers.FmtUUID(s.ActiveRole.ID),
		RoleType:  handlers.FmtString(string(s.ActiveRole.RoleType)),
		CreatedAt: handlers.FmtDateTime(s.ActiveRole.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(s.ActiveRole.UpdatedAt),
	}
	for _, role := range allUserRoles {
		// Make sure we don't accidentally mark the current role as inactive
		if role.RoleType != s.ActiveRole.RoleType {
			p.InactiveRoles = append(p.InactiveRoles, &internalmessages.Role{
				ID:        handlers.FmtUUID(role.ID),
				RoleType:  handlers.FmtString(string(role.RoleType)),
				CreatedAt: handlers.FmtDateTime(role.CreatedAt),
				UpdatedAt: handlers.FmtDateTime(role.UpdatedAt),
			})
		}
	}
}

// decoratePayloadWithPermissions will add session permissions to the logged in user payload and return it
func decoratePayloadWithPermissions(s *auth.Session, p *internalmessages.LoggedInUserPayload) {
	p.Permissions = []string{}
	p.Permissions = append(p.Permissions, s.Permissions...)
}

func decoratePayloadWithPrivileges(appCtx appcontext.AppContext, p *internalmessages.LoggedInUserPayload) {
	privileges, _ := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)

	for _, privilege := range privileges {
		p.Privileges = append(p.Privileges, &internalmessages.Privilege{
			PrivilegeType: *handlers.FmtString(string(privilege.PrivilegeType)),
		})
	}
}

// Helper function explicitly for these handlers.
// We'll take the session UserID, fetch roles from the DB,
// and return a default if found as well as all roles.
// Returning no rows is ok!
// The app supports authenticating users with no role.
// The app also handles the authorization of role-specific routes and action
// This helper function is due to how verbose the err checking is and to keep the parent funcs clean
func getDefaultAndAllRoles(appCtx appcontext.AppContext, userID uuid.UUID) (*roles.Role, roles.Roles, error) {
	if userID != uuid.Nil {
		userRoles, err := roles.FetchRolesForUser(appCtx.DB(), appCtx.Session().UserID)
		if err != nil && errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			// An err is not thrown when empty in this check.
			// An err is only thrown when there is an actual database problem,
			// as the query FetchRolesForUser uses does not return SqlErrNoRows
			appCtx.Logger().Error("database error when fetching roles for user",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.Error(err),
			)
			return nil, nil, err
		} else {
			// They have roles and their session doesn't have one yet
			// Set a default
			defaultRole, err := userRoles.Default()
			if err != nil {
				// This err occurs when no roles. Let them proceed
				appCtx.Logger().Error("could not find a default role for the logged in user, proceeding without a role",
					zap.String("userID", appCtx.Session().UserID.String()),
					zap.Error(err))
			} else {
				return defaultRole, userRoles, nil
			}
		}
	}
	return nil, nil, nil
}

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Pull latest roles from the DB in case the administrators
			// changed something. Additionally set a default role if not
			// done yet.
			defaultRole, userRoles, err := getDefaultAndAllRoles(appCtx, appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retrieving user roles", zap.Error(err))
				return userop.NewIsLoggedInUserInternalServerError(), err
			}
			if (appCtx.Session().ActiveRole.RoleType == roles.Role{}.RoleType) && defaultRole != nil {
				appCtx.Session().ActiveRole = *defaultRole
			}

			if !appCtx.Session().IsServiceMember() {
				var officeUser models.OfficeUser
				var err error
				if appCtx.Session().OfficeUserID != uuid.Nil {
					officeUser, err = h.officeUserFetcherPop.FetchOfficeUserByIDWithTransportationOfficeAssignments(appCtx, appCtx.Session().OfficeUserID)
					if err != nil {
						appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
						return userop.NewIsLoggedInUserInternalServerError(), err
					}
				}
				if (appCtx.Session().ActiveOfficeID == uuid.Nil) && officeUser.PrimaryOffice().ID != uuid.Nil {
					appCtx.Session().ActiveOfficeID = officeUser.PrimaryOffice().ID
				}
				var activeOffice models.TransportationOffice
				if officeUser.PrimaryOffice().ID == appCtx.Session().ActiveOfficeID {
					appCtx.Session().ActiveOfficeID = officeUser.PrimaryOffice().ID
					activeOffice = officeUser.PrimaryOffice()
				}
				if officeUser.SecondaryOffice().ID == appCtx.Session().ActiveOfficeID {
					appCtx.Session().ActiveOfficeID = officeUser.SecondaryOffice().ID
					activeOffice = officeUser.SecondaryOffice()
				}

				userPayload := internalmessages.LoggedInUserPayload{
					ID:           handlers.FmtUUID(appCtx.Session().UserID),
					FirstName:    appCtx.Session().FirstName,
					Email:        appCtx.Session().Email,
					OfficeUser:   payloads.OfficeUser(&officeUser),
					ActiveOffice: payloads.TransportationOffice(activeOffice),
				}

				decoratePayloadWithCurrentAndInactiveRoles(appCtx.Session(), &userPayload, userRoles)
				decoratePayloadWithPermissions(appCtx.Session(), &userPayload)
				decoratePayloadWithPrivileges(appCtx, &userPayload)

				return userop.NewShowLoggedInUserOK().WithPayload(&userPayload), nil
			}

			// Load Servicemember and first level associations
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), appCtx.Session().ServiceMemberID)

			if err != nil {
				appCtx.Logger().Error("Error retrieving service_member", zap.Error(err))
				return userop.NewShowLoggedInUserUnauthorized(), err
			}

			// Load the latest orders associations and new duty location transport office
			if len(serviceMember.Orders) > 0 {
				orders, orderErr := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), serviceMember.Orders[0].ID)
				if orderErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), orderErr), orderErr
				}

				serviceMember.Orders[0] = orders

				newDutyLocationTransportationOffice, dutyLocationErr := models.FetchDutyLocationTransportationOffice(appCtx.DB(), orders.NewDutyLocationID)
				if dutyLocationErr != nil {
					if errors.Cause(dutyLocationErr) != models.ErrFetchNotFound {
						// The absence of an office shouldn't render the entire request a 404
						return handlers.ResponseForError(appCtx.Logger(), dutyLocationErr), dutyLocationErr
					}
				}
				serviceMember.Orders[0].NewDutyLocation.TransportationOffice = newDutyLocationTransportationOffice

				// Load associations on PPM if they exist
				if len(serviceMember.Orders[0].Moves) > 0 {

					// Check if move is valid and not hidden
					// If the move is hidden, return an error
					if !(*serviceMember.Orders[0].Moves[0].Show) {
						return userop.NewShowLoggedInUserUnauthorized(), apperror.NewForbiddenError("user unauthorized to access move")
					}
				}
			}

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			userPayload := internalmessages.LoggedInUserPayload{
				ID:            handlers.FmtUUID(appCtx.Session().UserID),
				ServiceMember: payloadForServiceMemberModel(h.FileStorer(), serviceMember),
				FirstName:     appCtx.Session().FirstName,
				Email:         appCtx.Session().Email,
			}

			decoratePayloadWithCurrentAndInactiveRoles(appCtx.Session(), &userPayload, userRoles)
			decoratePayloadWithPermissions(appCtx.Session(), &userPayload)
			decoratePayloadWithPrivileges(appCtx, &userPayload)
			return userop.NewShowLoggedInUserOK().WithPayload(&userPayload), nil
		})
}
