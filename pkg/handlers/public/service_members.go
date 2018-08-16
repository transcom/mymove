package public

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
)

func publicPayloadForServiceMemberModel(serviceMember *models.ServiceMember) *apimessages.ServiceMember {

	serviceMemberPayload := apimessages.ServiceMember{
		FirstName:              serviceMember.FirstName,
		MiddleName:             serviceMember.MiddleName,
		LastName:               serviceMember.LastName,
		Suffix:                 serviceMember.Suffix,
		Telephone:              serviceMember.Telephone,
		SecondaryTelephone:     serviceMember.SecondaryTelephone,
		PersonalEmail:          serviceMember.PersonalEmail,
		PhoneIsPreferred:       serviceMember.PhoneIsPreferred,
		TextMessageIsPreferred: serviceMember.TextMessageIsPreferred,
		EmailIsPreferred:       serviceMember.EmailIsPreferred,
	}
	return &serviceMemberPayload
}
