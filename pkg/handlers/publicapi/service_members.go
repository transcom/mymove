package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForServiceMemberModel(serviceMember *models.ServiceMember) *apimessages.ServiceMember {
	contactPayloads := make(apimessages.IndexServiceMemberBackupContacts, len(serviceMember.BackupContacts))
	for i, contact := range serviceMember.BackupContacts {
		contactPayload := payloadForBackupContactModel(contact)
		contactPayloads[i] = &contactPayload
	}

	serviceMemberPayload := apimessages.ServiceMember{
		FirstName:              serviceMember.FirstName,
		MiddleName:             serviceMember.MiddleName,
		LastName:               serviceMember.LastName,
		Suffix:                 serviceMember.Suffix,
		Edipi:                  serviceMember.Edipi,
		Affiliation:            (*apimessages.Affiliation)(serviceMember.Affiliation),
		Rank:                   (*apimessages.ServiceMemberRank)(serviceMember.Rank),
		Telephone:              serviceMember.Telephone,
		SecondaryTelephone:     serviceMember.SecondaryTelephone,
		PersonalEmail:          serviceMember.PersonalEmail,
		PhoneIsPreferred:       serviceMember.PhoneIsPreferred,
		TextMessageIsPreferred: serviceMember.TextMessageIsPreferred,
		EmailIsPreferred:       serviceMember.EmailIsPreferred,
		BackupContacts:         contactPayloads,
	}

	return &serviceMemberPayload
}
