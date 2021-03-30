package internalapi

import (
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/cli"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForServiceMemberModel(storer storage.FileStorer, serviceMember models.ServiceMember, requiresAccessCode bool) *internalmessages.ServiceMemberPayload {
	orders := make([]*internalmessages.Orders, len(serviceMember.Orders))
	for i, order := range serviceMember.Orders {
		orderPayload, _ := payloadForOrdersModel(storer, order)
		orders[i] = orderPayload
	}

	contactPayloads := make(internalmessages.IndexServiceMemberBackupContactsPayload, len(serviceMember.BackupContacts))
	for i, contact := range serviceMember.BackupContacts {
		contactPayload := payloadForBackupContactModel(contact)
		contactPayloads[i] = &contactPayload
	}

	// if an existing service member, set requires access code to what they're already set
	if requiresAccessCode != serviceMember.RequiresAccessCode {
		requiresAccessCode = serviceMember.RequiresAccessCode
	}

	var weightAllotment *internalmessages.WeightAllotment
	if serviceMember.Rank != nil {
		weightAllotment = payloadForWeightAllotmentModel(models.GetWeightAllotment(*serviceMember.Rank))
	}

	serviceMemberPayload := internalmessages.ServiceMemberPayload{
		ID:                   handlers.FmtUUID(serviceMember.ID),
		CreatedAt:            handlers.FmtDateTime(serviceMember.CreatedAt),
		UpdatedAt:            handlers.FmtDateTime(serviceMember.UpdatedAt),
		UserID:               handlers.FmtUUID(serviceMember.UserID),
		Edipi:                serviceMember.Edipi,
		Orders:               orders,
		Affiliation:          (*internalmessages.Affiliation)(serviceMember.Affiliation),
		Rank:                 (*internalmessages.ServiceMemberRank)(serviceMember.Rank),
		FirstName:            serviceMember.FirstName,
		MiddleName:           serviceMember.MiddleName,
		LastName:             serviceMember.LastName,
		Suffix:               serviceMember.Suffix,
		Telephone:            serviceMember.Telephone,
		SecondaryTelephone:   serviceMember.SecondaryTelephone,
		PhoneIsPreferred:     serviceMember.PhoneIsPreferred,
		PersonalEmail:        serviceMember.PersonalEmail,
		EmailIsPreferred:     serviceMember.EmailIsPreferred,
		ResidentialAddress:   payloads.Address(serviceMember.ResidentialAddress),
		BackupMailingAddress: payloads.Address(serviceMember.BackupMailingAddress),
		BackupContacts:       contactPayloads,
		IsProfileComplete:    handlers.FmtBool(serviceMember.IsProfileComplete()),
		CurrentStation:       payloadForDutyStationModel(serviceMember.DutyStation),
		RequiresAccessCode:   requiresAccessCode,
		WeightAllotment:      weightAllotment,
	}
	return &serviceMemberPayload
}

// CreateServiceMemberHandler creates a new service member via POST /serviceMember
type CreateServiceMemberHandler struct {
	handlers.HandlerContext
}

// Handle ... creates a new ServiceMember from a request payload
func (h CreateServiceMemberHandler) Handle(params servicememberop.CreateServiceMemberParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	// User should always be populated by middleware
	session, logger := h.SessionAndLoggerFromContext(ctx)

	residentialAddress := addressModelFromPayload(params.CreateServiceMemberPayload.ResidentialAddress)
	backupMailingAddress := addressModelFromPayload(params.CreateServiceMemberPayload.BackupMailingAddress)

	var stationID *uuid.UUID
	var station models.DutyStation
	if params.CreateServiceMemberPayload.CurrentStationID != nil {
		id, err := uuid.FromString(params.CreateServiceMemberPayload.CurrentStationID.String())
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}
		s, err := models.FetchDutyStation(h.DB(), id)
		if err != nil {
			return handlers.ResponseForError(logger, err)
		}
		stationID = &id
		station = s
	}

	// Create a new serviceMember for an authenticated user
	newServiceMember := models.ServiceMember{
		UserID:               session.UserID,
		Edipi:                params.CreateServiceMemberPayload.Edipi,
		Affiliation:          (*models.ServiceMemberAffiliation)(params.CreateServiceMemberPayload.Affiliation),
		Rank:                 (*models.ServiceMemberRank)(params.CreateServiceMemberPayload.Rank),
		FirstName:            params.CreateServiceMemberPayload.FirstName,
		MiddleName:           params.CreateServiceMemberPayload.MiddleName,
		LastName:             params.CreateServiceMemberPayload.LastName,
		Suffix:               params.CreateServiceMemberPayload.Suffix,
		Telephone:            params.CreateServiceMemberPayload.Telephone,
		SecondaryTelephone:   params.CreateServiceMemberPayload.SecondaryTelephone,
		PersonalEmail:        params.CreateServiceMemberPayload.PersonalEmail,
		PhoneIsPreferred:     params.CreateServiceMemberPayload.PhoneIsPreferred,
		EmailIsPreferred:     params.CreateServiceMemberPayload.EmailIsPreferred,
		ResidentialAddress:   residentialAddress,
		BackupMailingAddress: backupMailingAddress,
		DutyStation:          station,
		RequiresAccessCode:   h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode),
		DutyStationID:        stationID,
	}
	smVerrs, err := models.SaveServiceMember(h.DB(), &newServiceMember)
	if smVerrs.HasAny() || err != nil {
		return handlers.ResponseForError(logger, err)
	}
	// Update session info
	session.ServiceMemberID = newServiceMember.ID

	if newServiceMember.FirstName != nil {
		session.FirstName = *(newServiceMember.FirstName)
	}
	if newServiceMember.MiddleName != nil {
		session.Middle = *(newServiceMember.MiddleName)
	}
	if newServiceMember.LastName != nil {
		session.LastName = *(newServiceMember.LastName)
	}
	// And return
	serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), newServiceMember, h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode))
	responder := servicememberop.NewCreateServiceMemberCreated().WithPayload(serviceMemberPayload)
	sessionManager := h.SessionManager(session)
	return handlers.NewCookieUpdateResponder(params.HTTPRequest, logger, responder, sessionManager, session)
}

// ShowServiceMemberHandler returns a serviceMember for a user and service member ID
type ShowServiceMemberHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h ShowServiceMemberHandler) Handle(params servicememberop.ShowServiceMemberParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	// User should always be populated by middleware
	session, logger := h.SessionAndLoggerFromContext(ctx)

	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

	serviceMember, err := models.FetchServiceMemberForUser(h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), serviceMember, h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode))
	return servicememberop.NewShowServiceMemberOK().WithPayload(serviceMemberPayload)
}

// PatchServiceMemberHandler patches a serviceMember via PATCH /serviceMembers/{serviceMemberId}
type PatchServiceMemberHandler struct {
	handlers.HandlerContext
}

// Check to see if a move is in draft state. If there are no orders, then the
// move still counts as in draft state.
func (h PatchServiceMemberHandler) isDraftMove(serviceMember *models.ServiceMember) bool {
	if serviceMember.Orders == nil || len(serviceMember.Orders) <= 0 {
		return true
	}

	move := serviceMember.Orders[0].Moves[0]

	return move.Status == models.MoveStatusDRAFT
}

// Handle ... patches a new ServiceMember from a request payload
func (h PatchServiceMemberHandler) Handle(params servicememberop.PatchServiceMemberParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromContext(ctx)

	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

	var err error
	var serviceMember models.ServiceMember
	var verrs *validate.Errors

	serviceMember, err = models.FetchServiceMemberForUser(h.DB(), session, serviceMemberID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := params.PatchServiceMemberPayload

	if verrs, err = h.patchServiceMemberWithPayload(&serviceMember, payload); verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	if verrs, err = models.SaveServiceMember(h.DB(), &serviceMember); verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	if len(serviceMember.Orders) != 0 {
		// Will have to be refactored once we support multiple moves/orders
		order, err := models.FetchOrderForUser(h.DB(), session, serviceMember.Orders[0].ID)

		if err != nil {
			return handlers.ResponseForError(logger, err)
		}

		serviceMemberRank := (*string)(serviceMember.Rank)
		if serviceMemberRank != order.Grade {
			order.Grade = serviceMemberRank
		}

		if serviceMember.DutyStation.ID != order.OriginDutyStation.ID {
			order.OriginDutyStation = &serviceMember.DutyStation
			order.OriginDutyStationID = &serviceMember.DutyStation.ID
		}

		verrs, err = h.DB().ValidateAndSave(&order)
		if verrs.HasAny() || err != nil {
			return handlers.ResponseForVErrors(logger, verrs, err)
		}
		serviceMember.Orders[0] = order
	}

	serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), serviceMember, h.HandlerContext.GetFeatureFlag(cli.FeatureFlagAccessCode))
	return servicememberop.NewPatchServiceMemberOK().WithPayload(serviceMemberPayload)
}

func (h PatchServiceMemberHandler) patchServiceMemberWithPayload(serviceMember *models.ServiceMember, payload *internalmessages.PatchServiceMemberPayload) (*validate.Errors, error) {
	if h.isDraftMove(serviceMember) {
		if payload.CurrentStationID != nil {
			stationID, err := uuid.FromString(payload.CurrentStationID.String())
			if err != nil {
				return validate.NewErrors(), err
			}
			// Fetch the model partially as a validation on the ID
			station, err := models.FetchDutyStation(h.DB(), stationID)
			if err != nil {
				return validate.NewErrors(), err
			}
			serviceMember.DutyStation = station
			serviceMember.DutyStationID = &stationID
		}

		if payload.Affiliation != nil {
			serviceMember.Affiliation = (*models.ServiceMemberAffiliation)(payload.Affiliation)
		}

		if payload.Rank != nil {
			serviceMember.Rank = (*models.ServiceMemberRank)(payload.Rank)
		}
	}
	if payload.Edipi != nil {
		serviceMember.Edipi = payload.Edipi
	}

	if payload.FirstName != nil {
		serviceMember.FirstName = payload.FirstName
	}
	if payload.MiddleName != nil {
		serviceMember.MiddleName = payload.MiddleName
	}
	if payload.LastName != nil {
		serviceMember.LastName = payload.LastName
	}
	if payload.Suffix != nil {
		serviceMember.Suffix = payload.Suffix
	}
	if payload.Telephone != nil {
		serviceMember.Telephone = payload.Telephone
	}
	if payload.SecondaryTelephone != nil {
		serviceMember.SecondaryTelephone = payload.SecondaryTelephone
	}
	if payload.PersonalEmail != nil {
		serviceMember.PersonalEmail = payload.PersonalEmail
	}
	if payload.PhoneIsPreferred != nil {
		serviceMember.PhoneIsPreferred = payload.PhoneIsPreferred
	}
	if payload.EmailIsPreferred != nil {
		serviceMember.EmailIsPreferred = payload.EmailIsPreferred
	}

	if payload.ResidentialAddress != nil {
		if serviceMember.ResidentialAddress == nil {
			serviceMember.ResidentialAddress = addressModelFromPayload(payload.ResidentialAddress)
		} else {
			updateAddressWithPayload(serviceMember.ResidentialAddress, payload.ResidentialAddress)
		}
	}
	if payload.BackupMailingAddress != nil {
		if serviceMember.BackupMailingAddress == nil {
			serviceMember.BackupMailingAddress = addressModelFromPayload(payload.BackupMailingAddress)
		} else {
			updateAddressWithPayload(serviceMember.BackupMailingAddress, payload.BackupMailingAddress)
		}
	}

	return validate.NewErrors(), nil
}

// ShowServiceMemberOrdersHandler returns latest orders for a logged in serviceMember
type ShowServiceMemberOrdersHandler struct {
	handlers.HandlerContext
}

// Handle retrieves orders for a logged in service member
func (h ShowServiceMemberOrdersHandler) Handle(params servicememberop.ShowServiceMemberOrdersParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	session, logger := h.SessionAndLoggerFromContext(ctx)

	serviceMember, err := models.FetchServiceMemberForUser(h.DB(), session, session.ServiceMemberID)
	if err != nil {
		return servicememberop.NewShowServiceMemberOrdersNotFound()
	}

	order, err := serviceMember.FetchLatestOrder(session, h.DB())
	if err != nil {
		return servicememberop.NewShowServiceMemberOrdersNotFound()
	}

	orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return servicememberop.NewShowServiceMemberOrdersOK().WithPayload(orderPayload)
}
