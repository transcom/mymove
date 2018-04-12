package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForServiceMemberModel(user models.User, serviceMember models.ServiceMember) *internalmessages.ServiceMemberPayload {
	serviceMemberPayload := internalmessages.ServiceMemberPayload{
		ID:                        fmtUUID(serviceMember.ID),
		CreatedAt:                 fmtDateTime(serviceMember.CreatedAt),
		UpdatedAt:                 fmtDateTime(serviceMember.UpdatedAt),
		UserID:                    fmtUUID(user.ID),
		Edipi:                     serviceMember.Edipi,
		Branch:                    serviceMember.Branch,
		Rank:                      serviceMember.Rank,
		FirstName:                 serviceMember.FirstName,
		MiddleInitial:             serviceMember.MiddleInitial,
		LastName:                  serviceMember.LastName,
		Suffix:                    serviceMember.Suffix,
		Telephone:                 serviceMember.Telephone,
		SecondaryTelephone:        serviceMember.SecondaryTelephone,
		PersonalEmail:             serviceMember.PersonalEmail,
		PhoneIsPreferred:          serviceMember.PhoneIsPreferred,
		SecondaryPhoneIsPreferred: serviceMember.SecondaryPhoneIsPreferred,
		EmailIsPreferred:          serviceMember.EmailIsPreferred,
		ResidentialAddress:        payloadForAddressModel(serviceMember.ResidentialAddress),
		BackupMailingAddress:      payloadForAddressModel(serviceMember.BackupMailingAddress),
		IsProfileComplete:         fmtBool(serviceMember.IsProfileComplete()),
	}
	return &serviceMemberPayload
}

// CreateServiceMemberHandler creates a new service member via POST /serviceMember
type CreateServiceMemberHandler HandlerContext

// Handle ... creates a new ServiceMember from a request payload
func (h CreateServiceMemberHandler) Handle(params servicememberop.CreateServiceMemberParams) middleware.Responder {
	var response middleware.Responder
	residentialAddress := models.AddressModelFromPayload(params.CreateServiceMemberPayload.ResidentialAddress)
	backupMailingAddress := models.AddressModelFromPayload(params.CreateServiceMemberPayload.BackupMailingAddress)

	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	// Create a new serviceMember for an authenticated user
	newServiceMember := models.ServiceMember{
		UserID:                    user.ID,
		Edipi:                     params.CreateServiceMemberPayload.Edipi,
		Branch:                    params.CreateServiceMemberPayload.Branch,
		Rank:                      params.CreateServiceMemberPayload.Rank,
		FirstName:                 params.CreateServiceMemberPayload.FirstName,
		MiddleInitial:             params.CreateServiceMemberPayload.MiddleInitial,
		LastName:                  params.CreateServiceMemberPayload.LastName,
		Suffix:                    params.CreateServiceMemberPayload.Suffix,
		Telephone:                 params.CreateServiceMemberPayload.Telephone,
		SecondaryTelephone:        params.CreateServiceMemberPayload.SecondaryTelephone,
		PersonalEmail:             params.CreateServiceMemberPayload.PersonalEmail,
		PhoneIsPreferred:          params.CreateServiceMemberPayload.PhoneIsPreferred,
		SecondaryPhoneIsPreferred: params.CreateServiceMemberPayload.SecondaryPhoneIsPreferred,
		EmailIsPreferred:          params.CreateServiceMemberPayload.EmailIsPreferred,
		ResidentialAddress:        residentialAddress,
		BackupMailingAddress:      backupMailingAddress,
	}
	verrs, err := models.CreateServiceMemberWithAddresses(h.db, &newServiceMember)
	if verrs.HasAny() {
		h.logger.Error("DB Validation", zap.Error(verrs))
		response = servicememberop.NewCreateServiceMemberBadRequest()
	} else if err != nil {
		h.logger.Error("DB Insertion", zap.Error(err))
		response = servicememberop.NewCreateServiceMemberBadRequest()
	} else {
		servicememberPayload := payloadForServiceMemberModel(user, newServiceMember)
		response = servicememberop.NewCreateServiceMemberCreated().WithPayload(servicememberPayload)
	}
	return response
}

// ShowServiceMemberHandler returns a serviceMember for a user and service member ID
type ShowServiceMemberHandler HandlerContext

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h ShowServiceMemberHandler) Handle(params servicememberop.ShowServiceMemberParams) middleware.Responder {
	var response middleware.Responder
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	serviceMemberID, err := uuid.FromString(params.ServiceMemberID.String())
	if err != nil {
		response = servicememberop.NewShowServiceMemberBadRequest()
		return response
	}

	serviceMemberResult, err := models.GetServiceMemberForUser(h.db, user.ID, serviceMemberID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = servicememberop.NewShowServiceMemberInternalServerError()
	} else if !serviceMemberResult.IsValid() {
		switch errCode := serviceMemberResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound:
			response = servicememberop.NewShowServiceMemberNotFound()
		case models.FetchErrorForbidden:
			response = servicememberop.NewShowServiceMemberForbidden()
		default:
			response = servicememberop.NewShowServiceMemberInternalServerError()
		}
		return response

	} else {
		serviceMemberPayload := payloadForServiceMemberModel(user, serviceMemberResult.ServiceMember())
		response = servicememberop.NewShowServiceMemberOK().WithPayload(serviceMemberPayload)
	}
	return response
}

// PatchServiceMemberHandler patches a serviceMember via PATCH /serviceMembers/{serviceMemberId}
type PatchServiceMemberHandler HandlerContext

// Handle ... patches a new ServiceMember from a request payload
func (h PatchServiceMemberHandler) Handle(params servicememberop.PatchServiceMemberParams) middleware.Responder {
	var response middleware.Responder
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	serviceMemberID, err := uuid.FromString(params.ServiceMemberID.String())
	if err != nil {
		response = servicememberop.NewPatchServiceMemberBadRequest()
	}

	// Validate that this serviceMember belongs to the current user
	serviceMemberResult, err := models.GetServiceMemberForUser(h.db, user.ID, serviceMemberID)
	if err != nil {
		h.logger.Error("DB Error checking on serviceMember validity", zap.Error(err))
		response = servicememberop.NewPatchServiceMemberInternalServerError()
	} else if !serviceMemberResult.IsValid() {
		switch errCode := serviceMemberResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound:
			response = servicememberop.NewPatchServiceMemberNotFound()
		case models.FetchErrorForbidden:
			response = servicememberop.NewPatchServiceMemberForbidden()
		default:
			response = servicememberop.NewPatchServiceMemberInternalServerError()
		}
		return response
	} else { // The given serviceMember does belong to the current user.
		serviceMember := serviceMemberResult.ServiceMember()
		payload := params.PatchServiceMemberPayload

		verrs, err := serviceMember.PatchServiceMemberWithPayload(h.db, payload)

		if verrs.HasAny() {
			response = servicememberop.NewPatchServiceMemberBadRequest()
		} else if err != nil {
			response = servicememberop.NewPatchServiceMemberInternalServerError()
		} else {
			serviceMemberPayload := payloadForServiceMemberModel(user, serviceMember)
			response = servicememberop.NewPatchServiceMemberCreated().WithPayload(serviceMemberPayload)
		}
	}
	return response
}
