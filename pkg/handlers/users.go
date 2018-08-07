package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if !session.IsServiceMember() {
		userPayload := internalmessages.LoggedInUserPayload{
			ID: fmtUUID(session.UserID),
		}
		return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)
	}
	// Load Servicemember and first level associations
	serviceMember, err := models.FetchServiceMember(h.db, session, session.ServiceMemberID)
	if err != nil {
		h.logger.Error("Error retrieving service_member", zap.Error(err))
		return userop.NewShowLoggedInUserUnauthorized()
	}

	// Load duty station and transportation office association
	if serviceMember.DutyStationID != nil {
		// Fetch associations on duty station
		dutyStation, err := models.FetchDutyStation(h.db, *serviceMember.DutyStationID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		// Fetch duty station transportation office
		transportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, *serviceMember.DutyStationID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		serviceMember.DutyStation = dutyStation
		serviceMember.DutyStation.TransportationOffice = transportationOffice
	}

	// Load the latest orders associations and new duty station transport office
	if len(serviceMember.Orders) > 0 {
		orders, err := models.FetchOrderForUser(h.db, session, serviceMember.Orders[0].ID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		newDutyStationTransportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, orders.NewDutyStationID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		serviceMember.Orders[0] = orders
		serviceMember.Orders[0].NewDutyStation.TransportationOffice = newDutyStationTransportationOffice

		// Load associations on PPM if they exist
		if len(serviceMember.Orders[0].Moves) > 0 {
			if len(serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves) > 0 {
				// TODO: load advances on all ppms for the latest order's move
				ppm, err := models.FetchPersonallyProcuredMove(h.db, session, serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].ID)
				if err != nil {
					return responseForError(h.logger, err)
				}
				serviceMember.Orders[0].Moves[0].PersonallyProcuredMoves[0].Advance = ppm.Advance
			}
		}
	}

	userPayload := internalmessages.LoggedInUserPayload{
		ID:            fmtUUID(session.UserID),
		ServiceMember: payloadForServiceMemberModel(h.storage, serviceMember),
	}
	return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)

}
