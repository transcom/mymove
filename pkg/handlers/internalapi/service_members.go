package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForServiceMemberModel(storer storage.FileStorer, serviceMember models.ServiceMember) *internalmessages.ServiceMemberPayload {
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
		CurrentLocation:      payloadForDutyLocationModel(serviceMember.DutyLocation),
		WeightAllotment:      weightAllotment,
	}
	return &serviceMemberPayload
}

// CreateServiceMemberHandler creates a new service member via POST /serviceMember
type CreateServiceMemberHandler struct {
	handlers.HandlerConfig
}

// Handle ... creates a new ServiceMember from a request payload
func (h CreateServiceMemberHandler) Handle(params servicememberop.CreateServiceMemberParams) middleware.Responder {
	// User should always be populated by middleware
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			residentialAddress := addressModelFromPayload(params.CreateServiceMemberPayload.ResidentialAddress)
			backupMailingAddress := addressModelFromPayload(params.CreateServiceMemberPayload.BackupMailingAddress)

			var dutyLocationID *uuid.UUID
			var dutyLocation models.DutyLocation
			if params.CreateServiceMemberPayload.CurrentLocationID != nil {
				id, err := uuid.FromString(params.CreateServiceMemberPayload.CurrentLocationID.String())
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				s, err := models.FetchDutyLocation(appCtx.DB(), id)
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
				dutyLocationID = &id
				dutyLocation = s
			}

			// Create a new serviceMember for an authenticated user
			newServiceMember := models.ServiceMember{
				UserID:               appCtx.Session().UserID,
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
				DutyLocation:         dutyLocation,
				DutyLocationID:       dutyLocationID,
			}
			smVerrs, err := models.SaveServiceMember(appCtx, &newServiceMember)
			if smVerrs.HasAny() || err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// Update session info
			appCtx.Session().ServiceMemberID = newServiceMember.ID

			if newServiceMember.FirstName != nil {
				appCtx.Session().FirstName = *(newServiceMember.FirstName)
			}
			if newServiceMember.MiddleName != nil {
				appCtx.Session().Middle = *(newServiceMember.MiddleName)
			}
			if newServiceMember.LastName != nil {
				appCtx.Session().LastName = *(newServiceMember.LastName)
			}
			// And return
			serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), newServiceMember)
			responder := servicememberop.NewCreateServiceMemberCreated().WithPayload(serviceMemberPayload)
			sessionManager := h.SessionManager(appCtx.Session())
			return handlers.NewCookieUpdateResponder(params.HTTPRequest, responder, sessionManager, appCtx.Session()), nil
		})
}

// ShowServiceMemberHandler returns a serviceMember for a user and service member ID
type ShowServiceMemberHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h ShowServiceMemberHandler) Handle(params servicememberop.ShowServiceMemberParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), serviceMember)
			return servicememberop.NewShowServiceMemberOK().WithPayload(serviceMemberPayload), nil
		})
}

// PatchServiceMemberHandler patches a serviceMember via PATCH /serviceMembers/{serviceMemberId}
type PatchServiceMemberHandler struct {
	handlers.HandlerConfig
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
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

			var err error
			var serviceMember models.ServiceMember
			var verrs *validate.Errors

			serviceMember, err = models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := params.PatchServiceMemberPayload

			if verrs, err = h.patchServiceMemberWithPayload(appCtx, &serviceMember, payload); verrs.HasAny() || err != nil {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			if verrs, err = models.SaveServiceMember(appCtx, &serviceMember); verrs.HasAny() || err != nil {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			if len(serviceMember.Orders) != 0 && h.isDraftMove(&serviceMember) {
				// Will have to be refactored once we support multiple moves/orders
				order, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), serviceMember.Orders[0].ID)

				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}

				serviceMemberRank := (*string)(serviceMember.Rank)
				if serviceMemberRank != order.Grade {
					order.Grade = serviceMemberRank
				}

				if serviceMember.DutyLocation.ID != order.OriginDutyLocation.ID {
					order.OriginDutyLocation = &serviceMember.DutyLocation
					order.OriginDutyLocationID = &serviceMember.DutyLocation.ID
				}

				verrs, err = appCtx.DB().ValidateAndSave(&order)
				if verrs.HasAny() || err != nil {
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
				serviceMember.Orders[0] = order
			}

			serviceMemberPayload := payloadForServiceMemberModel(h.FileStorer(), serviceMember)
			return servicememberop.NewPatchServiceMemberOK().WithPayload(serviceMemberPayload), nil
		})
}

func (h PatchServiceMemberHandler) patchServiceMemberWithPayload(appCtx appcontext.AppContext, serviceMember *models.ServiceMember, payload *internalmessages.PatchServiceMemberPayload) (*validate.Errors, error) {
	if h.isDraftMove(serviceMember) {
		if payload.CurrentLocationID != nil {
			dutyLocationID, err := uuid.FromString(payload.CurrentLocationID.String())
			if err != nil {
				return validate.NewErrors(), err
			}
			// Fetch the model partially as a validation on the ID
			dutyLocation, err := models.FetchDutyLocation(appCtx.DB(), dutyLocationID)
			if err != nil {
				return validate.NewErrors(), err
			}
			serviceMember.DutyLocation = dutyLocation
			serviceMember.DutyLocationID = &dutyLocationID
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
	handlers.HandlerConfig
}

// Handle retrieves orders for a logged in service member
func (h ShowServiceMemberOrdersHandler) Handle(params servicememberop.ShowServiceMemberOrdersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			order, err := models.FetchLatestOrder(appCtx.Session(), appCtx.DB())
			if err != nil {
				return servicememberop.NewShowServiceMemberOrdersNotFound(), err
			}

			orderPayload, err := payloadForOrdersModel(h.FileStorer(), order)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return servicememberop.NewShowServiceMemberOrdersOK().WithPayload(orderPayload), nil
		})
}
