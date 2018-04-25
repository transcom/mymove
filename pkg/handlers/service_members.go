package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
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
		ID:                      fmtUUID(serviceMember.ID),
		CreatedAt:               fmtDateTime(serviceMember.CreatedAt),
		UpdatedAt:               fmtDateTime(serviceMember.UpdatedAt),
		UserID:                  fmtUUID(user.ID),
		Edipi:                   serviceMember.Edipi,
		Affiliation:             serviceMember.Affiliation,
		Rank:                    serviceMember.Rank,
		FirstName:               serviceMember.FirstName,
		MiddleName:              serviceMember.MiddleName,
		LastName:                serviceMember.LastName,
		Suffix:                  serviceMember.Suffix,
		Telephone:               serviceMember.Telephone,
		SecondaryTelephone:      serviceMember.SecondaryTelephone,
		PhoneIsPreferred:        serviceMember.PhoneIsPreferred,
		TextMessageIsPreferred:  serviceMember.TextMessageIsPreferred,
		EmailIsPreferred:        serviceMember.EmailIsPreferred,
		ResidentialAddress:      payloadForAddressModel(serviceMember.ResidentialAddress),
		BackupMailingAddress:    payloadForAddressModel(serviceMember.BackupMailingAddress),
		HasSocialSecurityNumber: fmtBool(serviceMember.SocialSecurityNumberID != nil),
		IsProfileComplete:       fmtBool(serviceMember.IsProfileComplete()),
	}
	return &serviceMemberPayload
}

func patchServiceMemberWithPayload(serviceMember *models.ServiceMember, payload *internalmessages.PatchServiceMemberPayload) (*validate.Errors, error) {
	if payload.Edipi != nil {
		serviceMember.Edipi = payload.Edipi
	}
	if payload.Affiliation != nil {
		serviceMember.Affiliation = payload.Affiliation
	}
	if payload.Rank != nil {
		serviceMember.Rank = payload.Rank
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
		serviceMember.PersonalEmail = swag.String(payload.PersonalEmail.String())
	}
	if payload.PhoneIsPreferred != nil {
		serviceMember.PhoneIsPreferred = payload.PhoneIsPreferred
	}
	if payload.TextMessageIsPreferred != nil {
		serviceMember.TextMessageIsPreferred = payload.TextMessageIsPreferred
	}
	if payload.EmailIsPreferred != nil {
		serviceMember.EmailIsPreferred = payload.EmailIsPreferred
	}
	if payload.SocialSecurityNumber != nil {
		if serviceMember.SocialSecurityNumber == nil {
			newSsn := models.SocialSecurityNumber{}
			serviceMember.SocialSecurityNumber = &newSsn
		}

		if verrs, err := serviceMember.SocialSecurityNumber.SetEncryptedHash(payload.SocialSecurityNumber.String()); verrs.HasAny() || err != nil {
			return verrs, err
		}
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

// CreateServiceMemberHandler creates a new service member via POST /serviceMember
type CreateServiceMemberHandler HandlerContext

// Handle ... creates a new ServiceMember from a request payload
func (h CreateServiceMemberHandler) Handle(params servicememberop.CreateServiceMemberParams) middleware.Responder {
	var response middleware.Responder
	residentialAddress := addressModelFromPayload(params.CreateServiceMemberPayload.ResidentialAddress)
	backupMailingAddress := addressModelFromPayload(params.CreateServiceMemberPayload.BackupMailingAddress)

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
		UserID:                 user.ID,
		Edipi:                  params.CreateServiceMemberPayload.Edipi,
		Affiliation:            params.CreateServiceMemberPayload.Affiliation,
		Rank:                   params.CreateServiceMemberPayload.Rank,
		FirstName:              params.CreateServiceMemberPayload.FirstName,
		MiddleName:             params.CreateServiceMemberPayload.MiddleName,
		LastName:               params.CreateServiceMemberPayload.LastName,
		Suffix:                 params.CreateServiceMemberPayload.Suffix,
		Telephone:              params.CreateServiceMemberPayload.Telephone,
		SecondaryTelephone:     params.CreateServiceMemberPayload.SecondaryTelephone,
		PersonalEmail:          stringFromEmail(params.CreateServiceMemberPayload.PersonalEmail),
		PhoneIsPreferred:       params.CreateServiceMemberPayload.PhoneIsPreferred,
		TextMessageIsPreferred: params.CreateServiceMemberPayload.TextMessageIsPreferred,
		EmailIsPreferred:       params.CreateServiceMemberPayload.EmailIsPreferred,
		ResidentialAddress:     residentialAddress,
		BackupMailingAddress:   backupMailingAddress,
		SocialSecurityNumber:   ssn,
	}
	smVerrs, err := models.SaveServiceMember(h.db, &newServiceMember)
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
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
	serviceMember, err := models.FetchServiceMember(h.db, user, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	serviceMemberPayload := payloadForServiceMemberModel(user, serviceMember)
	return servicememberop.NewShowServiceMemberOK().WithPayload(serviceMemberPayload)
}

// PatchServiceMemberHandler patches a serviceMember via PATCH /serviceMembers/{serviceMemberId}
type PatchServiceMemberHandler HandlerContext

// Handle ... patches a new ServiceMember from a request payload
func (h PatchServiceMemberHandler) Handle(params servicememberop.PatchServiceMemberParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
	serviceMember, err := models.FetchServiceMember(h.db, user, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.PatchServiceMemberPayload
	if verrs, err := patchServiceMemberWithPayload(&serviceMember, payload); verrs.HasAny() || err != nil {
		return responseForVErrors(h.logger, verrs, err)
	}
	if verrs, err := models.SaveServiceMember(h.db, &serviceMember); verrs.HasAny() || err != nil {
		return responseForVErrors(h.logger, verrs, err)
	}

	serviceMemberPayload := payloadForServiceMemberModel(user, serviceMember)
	return servicememberop.NewPatchServiceMemberOK().WithPayload(serviceMemberPayload)
}
