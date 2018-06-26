package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	// "github.com/transcom/mymove/pkg/storage"
)

// func payloadForUserModel(storer storage.FileStorer, user *models.User, serviceMember *models.ServiceMember) *internalmessages.LoggedInUserPayload {
// 	var smPayload *internalmessages.ServiceMemberPayload

// 	if serviceMember != nil {
// 		smPayload = payloadForServiceMemberModel(storer, *serviceMember)
// 	}

// 	userPayload := internalmessages.LoggedInUserPayload{
// 		ID:            fmtUUID(user.ID),
// 		CreatedAt:     fmtDateTime(user.CreatedAt),
// 		ServiceMember: smPayload,
// 		UpdatedAt:     fmtDateTime(user.UpdatedAt),
// 	}
// 	return &userPayload
// }

// ShowLoggedInUserHandler returns the logged in user
type ShowLoggedInUserHandler HandlerContext

// Handle returns the logged in user
func (h ShowLoggedInUserHandler) Handle(params userop.ShowLoggedInUserParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// var user *models.User
	// var response middleware.Responder

	// serviceMember, err := models.GetFullServiceMemberProfile(h.db, session)
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
		dutyStation, err := models.FetchDutyStation(h.db, *serviceMember.DutyStationID)
		if err != nil {
			// h.logger.Error("Error retrieving duty station", zap.Error(err))
			// return userop.NewShowLoggedInUserUnauthorized()
			return responseForError(h.logger, err)

		}
		transportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, dutyStation.ID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		serviceMember.DutyStation = dutyStation
		serviceMember.DutyStation.TransportationOffice = transportationOffice
	}

	// Load latest orders association and new duty station transport office
	if len(serviceMember.Orders) > 0 {
		orders, err := models.FetchOrder(h.db, session, serviceMember.Orders[0].ID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		newDutyStationTransportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, orders.NewDutyStationID)
		if err != nil {
			return responseForError(h.logger, err)
		}
		serviceMember.Orders[0] = orders
		serviceMember.Orders[0].NewDutyStation.TransportationOffice = newDutyStationTransportationOffice
	}

	// Load associations on PPM
	// fmt.Println("orders duty station transport office", serviceMember.Orders[0].NewDutyStation.TransportationOffice.PhoneLines[0].Number)
	// fmt.Println("Ppms", serviceMember.Orders[1].Moves[0].PersonallyProcuredMoves[0])
	fmt.Println("Ppms", serviceMember.Orders[1].Moves)

	// TODO: Fetch everything
	// Check if service member fields are getting poulated - res address, bakcup address, backup contacts
	// if so, populate orders on a eager load in orders
	// same for duty station
	// same for ppms

	userPayload := internalmessages.LoggedInUserPayload{
		ID:            fmtUUID(session.UserID),
		ServiceMember: payloadForServiceMemberModel(h.storage, serviceMember),
	}
	// fmt.Println("User payload", userPayload)
	// fmt.Println("Is service member", session.IsServiceMember())
	// fmt.Println("ServiceMember", serviceMember.BackupContacts[0].Name)

	return userop.NewShowLoggedInUserOK().WithPayload(&userPayload)

}
