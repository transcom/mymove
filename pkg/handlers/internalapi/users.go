package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

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
	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsServiceMember() {
		var officeUser models.OfficeUser
		var err error
		if session.OfficeUserID != uuid.Nil {
			officeUser, err = h.officeUserFetcherPop.FetchOfficeUserByID(session.OfficeUserID)
			if err != nil {
				logger.Error("Error retrieving office_user", zap.Error(err))
				return userop.NewIsLoggedInUserInternalServerError()
			}
		}

		userPayload := internalmessages.LoggedInUserPayload{
			ID:         handlers.FmtUUID(session.UserID),
			FirstName:  session.FirstName,
			Email:      session.Email,
			OfficeUser: payloads.OfficeUser(&officeUser),
		}
		decoratePayloadWithRoles(session, &userPayload)
		return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
	}

	// Load Servicemember and first level associations
	serviceMember, err := models.FetchServiceMemberForUser(ctx, h.DB(), session, session.ServiceMemberID)

	if err != nil {
		logger.Error("Error retrieving service_member", zap.Error(err))
		return userop.NewShowLoggedInUserUnauthorized()
	}

	// Load duty station and transportation office association
	if serviceMember.DutyStationID != nil {
		// Fetch associations on duty station
		dutyStation, dutyStationErr := models.FetchDutyStation(h.DB(), *serviceMember.DutyStationID)
		if dutyStationErr != nil {
			return handlers.ResponseForError(logger, dutyStationErr)
		}
		serviceMember.DutyStation = dutyStation

		// Fetch duty station transportation office
		transportationOffice, tspErr := models.FetchDutyStationTransportationOffice(h.DB(), *serviceMember.DutyStationID)
		if tspErr != nil {
			if errors.Cause(tspErr) != models.ErrFetchNotFound {
				// The absence of an office shouldn't render the entire request a 404
				return handlers.ResponseForError(logger, tspErr)
			}
			// We might not have Transportation Office data for a Duty Station, and that's ok
			if errors.Cause(tspErr) != models.ErrFetchNotFound {
				return handlers.ResponseForError(logger, tspErr)
			}
		}
		serviceMember.DutyStation.TransportationOffice = transportationOffice
	}

	// Load the latest orders associations and new duty station transport office
	if len(serviceMember.Orders) > 0 {
		orders, orderErr := models.FetchOrderForUser(h.DB(), session, serviceMember.Orders[0].ID)
		if orderErr != nil {
			return handlers.ResponseForError(logger, orderErr)
		}

		serviceMember.Orders[0] = orders

		newDutyStationTransportationOffice, dutyStationErr := models.FetchDutyStationTransportationOffice(h.DB(), orders.NewDutyStationID)
		if dutyStationErr != nil {
			if errors.Cause(err) != models.ErrFetchNotFound {
				// The absence of an office shouldn't render the entire request a 404
				return handlers.ResponseForError(logger, dutyStationErr)
			}
		}
		serviceMember.Orders[0].NewDutyStation.TransportationOffice = newDutyStationTransportationOffice

		// Load associations on PPM if they exist
		if len(serviceMember.Orders[0].Moves) > 0 {
			if len(serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves) > 0 {
				// TODO: load advances on all ppms for the latest order's move
				ppm, ppmErr := models.FetchPersonallyProcuredMove(h.DB(), session, serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].ID)
				if ppmErr != nil {
					return handlers.ResponseForError(logger, ppmErr)
				}
				serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].Advance = ppm.Advance
			}
		}
	}

	requiresAccessCode := h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode)

	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	userPayload := internalmessages.LoggedInUserPayload{
		ID:            handlers.FmtUUID(session.UserID),
		ServiceMember: payloadForServiceMemberModel(h.FileStorer(), serviceMember, requiresAccessCode),
		FirstName:     session.FirstName,
		Email:         session.Email,
	}
	decoratePayloadWithRoles(session, &userPayload)
	return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
}
