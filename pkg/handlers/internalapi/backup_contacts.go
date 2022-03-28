package internalapi

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	backupop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForBackupContactModel(contact models.BackupContact) internalmessages.ServiceMemberBackupContactPayload {
	permission := internalmessages.NewBackupContactPermission(internalmessages.BackupContactPermission(contact.Permission))
	contactPayload := internalmessages.ServiceMemberBackupContactPayload{
		ID:              handlers.FmtUUID(contact.ID),
		ServiceMemberID: *handlers.FmtUUID(contact.ServiceMemberID),
		UpdatedAt:       handlers.FmtDateTime(contact.UpdatedAt),
		CreatedAt:       handlers.FmtDateTime(contact.CreatedAt),
		Name:            &contact.Name,
		Email:           &contact.Email,
		Telephone:       contact.Phone,
		Permission:      permission,
	}
	return contactPayload
}

// CreateBackupContactHandler creates a new backup contact
type CreateBackupContactHandler struct {
	handlers.HandlerContext
}

// Handle ... creates a new BackupContact from a request payload
func (h CreateBackupContactHandler) Handle(params backupop.CreateServiceMemberBackupContactParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			if params.CreateBackupContactPayload.Permission == nil {
				err = errors.New("missing required field: Permission")
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newContact, verrs, err := serviceMember.CreateBackupContact(appCtx.DB(),
				*params.CreateBackupContactPayload.Name,
				*params.CreateBackupContactPayload.Email,
				params.CreateBackupContactPayload.Telephone,
				models.BackupContactPermission(*params.CreateBackupContactPayload.Permission))
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

			contactPayload := payloadForBackupContactModel(newContact)
			return backupop.NewCreateServiceMemberBackupContactCreated().WithPayload(&contactPayload), nil
		})
}

// IndexBackupContactsHandler returns a list of all backup contacts for a service member
type IndexBackupContactsHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a list of all moves in the system belonging to the logged in user
func (h IndexBackupContactsHandler) Handle(params backupop.IndexServiceMemberBackupContactsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceMemberID, _ := uuid.FromString(params.ServiceMemberID.String())
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			contacts := serviceMember.BackupContacts

			contactPayloads := make(internalmessages.IndexServiceMemberBackupContactsPayload, len(contacts))
			for i, contact := range contacts {
				contactPayload := payloadForBackupContactModel(contact)
				contactPayloads[i] = &contactPayload
			}

			return backupop.NewIndexServiceMemberBackupContactsOK().WithPayload(contactPayloads), nil
		})
}

// ShowBackupContactHandler returns a backup contact for a user and backup contact ID
type ShowBackupContactHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a backup contact in the system belonging to the logged in user given backup contact ID
func (h ShowBackupContactHandler) Handle(params backupop.ShowServiceMemberBackupContactParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			contactID, _ := uuid.FromString(params.BackupContactID.String())
			contact, err := models.FetchBackupContact(appCtx.DB(), appCtx.Session(), contactID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			contactPayload := payloadForBackupContactModel(contact)
			return backupop.NewShowServiceMemberBackupContactOK().WithPayload(&contactPayload)
		})
}

// UpdateBackupContactHandler updates a backup contact with a new one
type UpdateBackupContactHandler struct {
	handlers.HandlerContext
}

// Handle ... updates a BackupContact from a request payload
func (h UpdateBackupContactHandler) Handle(params backupop.UpdateServiceMemberBackupContactParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			contactID, _ := uuid.FromString(params.BackupContactID.String())
			contact, err := models.FetchBackupContact(appCtx.DB(), appCtx.Session(), contactID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			contact.Name = *params.UpdateServiceMemberBackupContactPayload.Name
			contact.Email = *params.UpdateServiceMemberBackupContactPayload.Email
			contact.Phone = params.UpdateServiceMemberBackupContactPayload.Telephone
			if params.UpdateServiceMemberBackupContactPayload.Permission != nil {
				contact.Permission = models.BackupContactPermission(*params.UpdateServiceMemberBackupContactPayload.Permission)
			}

			if verrs, err := appCtx.DB().ValidateAndUpdate(&contact); verrs.HasAny() || err != nil {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
			}

			contactPayload := payloadForBackupContactModel(contact)
			return backupop.NewUpdateServiceMemberBackupContactCreated().WithPayload(&contactPayload)
		})
}
