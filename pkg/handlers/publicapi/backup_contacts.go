package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForBackupContactModel(contact models.BackupContact) apimessages.ServiceMemberBackupContact {
	contactPayload := apimessages.ServiceMemberBackupContact{
		ID:              handlers.FmtUUID(contact.ID),
		ServiceMemberID: *handlers.FmtUUID(contact.ServiceMemberID),
		UpdatedAt:       handlers.FmtDateTime(contact.UpdatedAt),
		CreatedAt:       handlers.FmtDateTime(contact.CreatedAt),
		Name:            &contact.Name,
		Email:           &contact.Email,
		Telephone:       contact.Phone,
		Permission:      apimessages.BackupContactPermission(contact.Permission),
	}
	return contactPayload
}
