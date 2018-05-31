package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	backupop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForBackupContactModel(contact models.BackupContact) internalmessages.ServiceMemberBackupContactPayload {
	contactPayload := internalmessages.ServiceMemberBackupContactPayload{
		ID:         fmtUUID(contact.ID),
		UpdatedAt:  fmtDateTime(contact.UpdatedAt),
		CreatedAt:  fmtDateTime(contact.CreatedAt),
		Name:       &contact.Name,
		Email:      &contact.Email,
		Telephone:  contact.Phone,
		Permission: contact.Permission,
	}
	return contactPayload
}

// CreateBackupContactHandler creates a new backup contact
type CreateBackupContactHandler HandlerContext

// Handle ... creates a new BackupContact from a request payload
func (h CreateBackupContactHandler) Handle(params backupop.CreateServiceMemberBackupContactParams) middleware.Responder {
	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
	serviceMember, err := models.FetchServiceMember(h.db, session, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	newContact, verrs, err := serviceMember.CreateBackupContact(h.db,
		*params.CreateBackupContactPayload.Name,
		*params.CreateBackupContactPayload.Email,
		params.CreateBackupContactPayload.Telephone,
		params.CreateBackupContactPayload.Permission)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	contactPayload := payloadForBackupContactModel(newContact)
	return backupop.NewCreateServiceMemberBackupContactCreated().WithPayload(&contactPayload)
}

// IndexBackupContactsHandler returns a list of all backup contacts for a service member
type IndexBackupContactsHandler HandlerContext

// Handle retrieves a list of all moves in the system belonging to the logged in user
func (h IndexBackupContactsHandler) Handle(params backupop.IndexServiceMemberBackupContactsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	/* #nosec UUID is pattern matched by swagger which checks the format */
	serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
	serviceMember, err := models.FetchServiceMember(h.db, session, serviceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	contacts := serviceMember.BackupContacts

	contactPayloads := make(internalmessages.IndexServiceMemberBackupContactsPayload, len(contacts))
	for i, contact := range contacts {
		contactPayload := payloadForBackupContactModel(contact)
		contactPayloads[i] = &contactPayload
	}

	return backupop.NewIndexServiceMemberBackupContactsOK().WithPayload(contactPayloads)
}

// ShowBackupContactHandler returns a backup contact for a user and backup contact ID
type ShowBackupContactHandler HandlerContext

// Handle retrieves a backup contact in the system belonging to the logged in user given backup contact ID
func (h ShowBackupContactHandler) Handle(params backupop.ShowServiceMemberBackupContactParams) middleware.Responder {
	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	contactID, _ := uuid.FromString(params.BackupContactID.String())
	contact, err := models.FetchBackupContact(h.db, session, contactID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	contactPayload := payloadForBackupContactModel(contact)
	return backupop.NewShowServiceMemberBackupContactOK().WithPayload(&contactPayload)
}

// UpdateBackupContactHandler updates a backup contact with a new one
type UpdateBackupContactHandler HandlerContext

// Handle ... updates a BackupContact from a request payload
func (h UpdateBackupContactHandler) Handle(params backupop.UpdateServiceMemberBackupContactParams) middleware.Responder {
	// User should always be populated by middleware
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	/* #nosec UUID is pattern matched by swagger which checks the format */
	contactID, _ := uuid.FromString(params.BackupContactID.String())
	contact, err := models.FetchBackupContact(h.db, session, contactID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	contact.Name = *params.UpdateServiceMemberBackupContactPayload.Name
	contact.Email = *params.UpdateServiceMemberBackupContactPayload.Email
	contact.Phone = params.UpdateServiceMemberBackupContactPayload.Telephone
	contact.Permission = params.UpdateServiceMemberBackupContactPayload.Permission

	if verrs, err := h.db.ValidateAndUpdate(&contact); verrs.HasAny() || err != nil {
		return responseForVErrors(h.logger, verrs, err)
	}

	contactPayload := payloadForBackupContactModel(contact)
	return backupop.NewUpdateServiceMemberBackupContactCreated().WithPayload(&contactPayload)
}
