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
	"github.com/transcom/mymove/pkg/services"
)

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler struct {
	handlers.HandlerConfig
	officeUserFetcherPop services.OfficeUserFetcherPop
}

// decoratePayloadWithRoles will add session roles to the logged in user payload and return it
func decoratePayloadWithRoles(s *auth.Session, p *internalmessages.LoggedInUserPayload) {
	for _, role := range s.Roles {
		p.Roles = append(p.Roles, &internalmessages.Role{
			ID:        handlers.FmtUUID(s.UserID),
			RoleType:  handlers.FmtString(string(role.RoleType)),
			CreatedAt: handlers.FmtDateTime(role.CreatedAt),
			UpdatedAt: handlers.FmtDateTime(role.UpdatedAt),
		})
	}
}

// decoratePayloadWithPermissions will add session permissions to the logged in user payload and return it
func decoratePayloadWithPermissions(s *auth.Session, p *internalmessages.LoggedInUserPayload) {
	p.Permissions = []string{}
	p.Permissions = append(p.Permissions, s.Permissions...)
}

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsServiceMember() {
				var officeUser models.OfficeUser
				var err error
				if appCtx.Session().OfficeUserID != uuid.Nil {
					officeUser, err = h.officeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
					if err != nil {
						appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
						return userop.NewIsLoggedInUserInternalServerError(), err
					}
				}

				userPayload := internalmessages.LoggedInUserPayload{
					ID:         handlers.FmtUUID(appCtx.Session().UserID),
					FirstName:  appCtx.Session().FirstName,
					Email:      appCtx.Session().Email,
					OfficeUser: payloads.OfficeUser(&officeUser),
				}
				decoratePayloadWithRoles(appCtx.Session(), &userPayload)
				decoratePayloadWithPermissions(appCtx.Session(), &userPayload)

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
			decoratePayloadWithRoles(appCtx.Session(), &userPayload)
			decoratePayloadWithPermissions(appCtx.Session(), &userPayload)
			return userop.NewShowLoggedInUserOK().WithPayload(&userPayload), nil
		})
}
