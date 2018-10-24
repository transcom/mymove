package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"go.uber.org/dig"
	"go.uber.org/zap"

	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ShowLoggedInUserHandler returns the logged in user
type showLoggedInUserHandler struct {
	db                 *pop.Connection
	logger             *zap.Logger
	storer             storage.FileStorer
	fetchServiceMember services.FetchServiceMember
}

// ShowLoggedInUserHandlerParams contains dependencies for NewShowLoggerInUserProvider
type ShowLoggedInUserHandlerParams struct {
	dig.In
	Db                 *pop.Connection
	Logger             *zap.Logger
	Storer             storage.FileStorer
	FetchServiceMember services.FetchServiceMember
}

// NewShowLoggedInUserHandler is a DI provider to generate the new Handler
func NewShowLoggedInUserHandler(params ShowLoggedInUserHandlerParams) userop.ShowLoggedInUserHandler {
	return &showLoggedInUserHandler{params.Db,
		params.Logger,
		params.Storer,
		params.FetchServiceMember}
}

// Handle returns the logged in user
func (h *showLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)
	if !session.IsServiceMember() {
		userPayload := internalmessages.LoggedInUserPayload{
			ID: handlers.FmtUUID(session.UserID),
		}
		return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
	}

	// Load Servicemember and first level associations
	serviceMember, err := h.fetchServiceMember.Execute(session, session.ServiceMemberID)
	if err != nil {
		h.logger.Error("Error retrieving service_member", zap.Error(err))
		return userop.NewShowLoggedInUserUnauthorized()
	}

	// Load duty station and transportation office association
	if serviceMember.DutyStationID != nil {
		// Fetch associations on duty station
		dutyStation, err := models.FetchDutyStation(h.db, *serviceMember.DutyStationID)
		if err != nil {
			return handlers.ResponseForError(h.logger, err)
		}
		serviceMember.DutyStation = dutyStation

		// Fetch duty station transportation office
		transportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, *serviceMember.DutyStationID)
		if err != nil {
			if errors.Cause(err) != models.ErrFetchNotFound {
				// The absence of an office shouldn't render the entire request a 404
				return handlers.ResponseForError(h.logger, err)
			}
			// We might not have Transportation Office data for a Duty Station, and that's ok
			if errors.Cause(err) != models.ErrFetchNotFound {
				return handlers.ResponseForError(h.logger, err)
			}
		}
		serviceMember.DutyStation.TransportationOffice = transportationOffice
	}

	// Load the latest orders associations and new duty station transport office
	if len(serviceMember.Orders) > 0 {
		orders, err := models.FetchOrderForUser(h.db, session, serviceMember.Orders[0].ID)
		if err != nil {
			return handlers.ResponseForError(h.logger, err)
		}
		serviceMember.Orders[0] = orders

		newDutyStationTransportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, orders.NewDutyStationID)
		if err != nil {
			if errors.Cause(err) != models.ErrFetchNotFound {
				// The absence of an office shouldn't render the entire request a 404
				return handlers.ResponseForError(h.logger, err)
			}
			// We might not have Transportation Office data for a Duty Station, and that's ok
			if errors.Cause(err) != models.ErrFetchNotFound {
				return handlers.ResponseForError(h.logger, err)
			}
		}
		serviceMember.Orders[0].NewDutyStation.TransportationOffice = newDutyStationTransportationOffice

		// Load associations on PPM if they exist
		if len(serviceMember.Orders[0].Moves) > 0 {
			if len(serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves) > 0 {
				// TODO: load advances on all ppms for the latest order's move
				ppm, err := models.FetchPersonallyProcuredMove(h.db, session, serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].ID)
				if err != nil {
					return handlers.ResponseForError(h.logger, err)
				}
				serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].Advance = ppm.Advance
			}
		}
	}

	userPayload := internalmessages.LoggedInUserPayload{
		ID:            handlers.FmtUUID(session.UserID),
		ServiceMember: payloadForServiceMemberModel(h.storer, *serviceMember),
	}
	return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
}
