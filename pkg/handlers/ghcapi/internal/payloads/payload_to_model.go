package payloads

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerToServiceMember transforms UpdateCustomerPayload to ServiceMember model
func CustomerToServiceMember(payload ghcmessages.UpdateCustomerPayload) models.ServiceMember {

	var address = models.Address{
		ID:             uuid.FromStringOrNil(payload.CurrentAddress.ID.String()),
		StreetAddress1: *payload.CurrentAddress.StreetAddress1,
		StreetAddress2: payload.CurrentAddress.StreetAddress2,
		StreetAddress3: payload.CurrentAddress.StreetAddress3,
		City:           *payload.CurrentAddress.City,
		State:          *payload.CurrentAddress.State,
		PostalCode:     *payload.CurrentAddress.PostalCode,
		Country:        payload.CurrentAddress.Country,
	}

	var backupContact = models.BackupContact{
		Email: *payload.BackupContact.Email,
		Name:  *payload.BackupContact.Name,
		Phone: payload.BackupContact.Phone,
	}

	var backupContacts []models.BackupContact
	backupContacts = append(backupContacts, backupContact)

	return models.ServiceMember{
		ResidentialAddress: &address,
		BackupContacts:     backupContacts,
		FirstName:          &payload.FirstName,
		LastName:           &payload.LastName,
		Suffix:             payload.Suffix,
		MiddleName:         payload.MiddleName,
		PersonalEmail:      payload.Email,
		Telephone:          payload.Phone,
	}
}
