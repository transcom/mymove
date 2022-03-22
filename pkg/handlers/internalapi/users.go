package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/cli"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler struct {
	handlers.HandlerContext
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

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			if !appCtx.Session().IsServiceMember() {
				var officeUser models.OfficeUser
				var err error
				if appCtx.Session().OfficeUserID != uuid.Nil {
					officeUser, err = h.officeUserFetcherPop.FetchOfficeUserByID(appCtx, appCtx.Session().OfficeUserID)
					if err != nil {
						appCtx.Logger().Error("Error retrieving office_user", zap.Error(err))
						return userop.NewIsLoggedInUserInternalServerError()
					}
				}

				userPayload := internalmessages.LoggedInUserPayload{
					ID:         handlers.FmtUUID(appCtx.Session().UserID),
					FirstName:  appCtx.Session().FirstName,
					Email:      appCtx.Session().Email,
					OfficeUser: payloads.OfficeUser(&officeUser),
				}
				decoratePayloadWithRoles(appCtx.Session(), &userPayload)
				return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
			}

			// Load Servicemember and first level associations
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), appCtx.Session().ServiceMemberID)

			if err != nil {
				appCtx.Logger().Error("Error retrieving service_member", zap.Error(err))
				return userop.NewShowLoggedInUserUnauthorized()
			}

			// Load duty station and transportation office association
			if serviceMember.DutyLocationID != nil {
				// Fetch associations on duty station
				dutyLocation, dutyLocationErr := models.FetchDutyLocation(appCtx.DB(), *serviceMember.DutyLocationID)
				if dutyLocationErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), dutyLocationErr)
				}
				serviceMember.DutyLocation = dutyLocation

				// Fetch duty station transportation office
				transportationOffice, tspErr := models.FetchDutyLocationTransportationOffice(appCtx.DB(), *serviceMember.DutyLocationID)
				if tspErr != nil {
					if errors.Cause(tspErr) != models.ErrFetchNotFound {
						// The absence of an office shouldn't render the entire request a 404
						return handlers.ResponseForError(appCtx.Logger(), tspErr)
					}
					// We might not have Transportation Office data for a Duty Station, and that's ok
					if errors.Cause(tspErr) != models.ErrFetchNotFound {
						return handlers.ResponseForError(appCtx.Logger(), tspErr)
					}
				}
				serviceMember.DutyLocation.TransportationOffice = transportationOffice
			}

			// Load the latest orders associations and new duty station transport office
			if len(serviceMember.Orders) > 0 {
				orders, orderErr := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), serviceMember.Orders[0].ID)
				if orderErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), orderErr)
				}

				serviceMember.Orders[0] = orders

				newDutyLocationTransportationOffice, dutyLocationErr := models.FetchDutyLocationTransportationOffice(appCtx.DB(), orders.NewDutyLocationID)
				if dutyLocationErr != nil {
					if errors.Cause(dutyLocationErr) != models.ErrFetchNotFound {
						// The absence of an office shouldn't render the entire request a 404
						return handlers.ResponseForError(appCtx.Logger(), dutyLocationErr)
					}
				}
				serviceMember.Orders[0].NewDutyLocation.TransportationOffice = newDutyLocationTransportationOffice

				// Load associations on PPM if they exist
				if len(serviceMember.Orders[0].Moves) > 0 {
					if len(serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves) > 0 {
						// TODO: load advances on all ppms for the latest order's move
						ppm, ppmErr := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].ID)
						if ppmErr != nil {
							return handlers.ResponseForError(appCtx.Logger(), ppmErr)
						}
						serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].Advance = ppm.Advance
					}

					// Check if move is valid and not hidden
					// If the move is hidden, return an error
					if !(*serviceMember.Orders[0].Moves[0].Show) {
						return userop.NewShowLoggedInUserUnauthorized()
					}
				}
			}

			requiresAccessCode := h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}
			userPayload := internalmessages.LoggedInUserPayload{
				ID:            handlers.FmtUUID(appCtx.Session().UserID),
				ServiceMember: payloadForServiceMemberModel(h.FileStorer(), serviceMember, requiresAccessCode),
				FirstName:     appCtx.Session().FirstName,
				Email:         appCtx.Session().Email,
			}
			decoratePayloadWithRoles(appCtx.Session(), &userPayload)
			return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
		})
}
