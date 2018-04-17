package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
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
		PersonalEmail:             (*strfmt.Email)(serviceMember.PersonalEmail),
		PhoneIsPreferred:          serviceMember.PhoneIsPreferred,
		SecondaryPhoneIsPreferred: serviceMember.SecondaryPhoneIsPreferred,
		EmailIsPreferred:          serviceMember.EmailIsPreferred,
		ResidentialAddress:        payloadForAddressModel(serviceMember.ResidentialAddress),
		BackupMailingAddress:      payloadForAddressModel(serviceMember.BackupMailingAddress),
		HasSocialSecurityNumber:   fmtBool(serviceMember.SocialSecurityNumberID != nil),
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

	ssnString := params.CreateServiceMemberPayload.SocialSecurityNumber
	var ssn *models.SocialSecurityNumber
	verrs := validate.NewErrors()
	if ssnString != nil {
		var err error
		ssn, verrs, err = models.BuildSocialSecurityNumber(ssnString.String())
		if err != nil {
			h.logger.Error("Unexpected error building SSN model", zap.Error(err))
			return servicememberop.NewCreateServiceMemberInternalServerError()
		}
		// if there are any validation errors, they will get rolled up with the rest of them.
	}

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
		PersonalEmail:             stringFromEmail(params.CreateServiceMemberPayload.PersonalEmail),
		PhoneIsPreferred:          params.CreateServiceMemberPayload.PhoneIsPreferred,
		SecondaryPhoneIsPreferred: params.CreateServiceMemberPayload.SecondaryPhoneIsPreferred,
		EmailIsPreferred:          params.CreateServiceMemberPayload.EmailIsPreferred,
		ResidentialAddress:        residentialAddress,
		BackupMailingAddress:      backupMailingAddress,
		SocialSecurityNumber:      ssn,
	}
	smVerrs, err := models.CreateServiceMember(h.db, &newServiceMember)
	verrs.Append(smVerrs)
	if verrs.HasAny() {
		h.logger.Info("DB Validation", zap.Error(verrs))
		response = servicememberop.NewCreateServiceMemberBadRequest()
	} else if err != nil {
		if err == models.ErrCreateViolatesUniqueConstraint {
			h.logger.Info("Attempted to create a second SM when one already exists")
			response = servicememberop.NewCreateServiceMemberBadRequest()
		} else {
			h.logger.Error("DB Insertion", zap.Error(err))
			response = servicememberop.NewCreateServiceMemberInternalServerError()
		}
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
	// swagger validates our UUID format.
	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())

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
			h.logger.Error("Unexpected error Fetching Service Member", zap.Error(err))
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
			h.logger.Error("Unexpected error Patching Service Member", zap.Error(err))
			response = servicememberop.NewPatchServiceMemberInternalServerError()
		} else {
			serviceMemberPayload := payloadForServiceMemberModel(user, serviceMember)
			response = servicememberop.NewPatchServiceMemberOK().WithPayload(serviceMemberPayload)
		}
	}

	return response
}
