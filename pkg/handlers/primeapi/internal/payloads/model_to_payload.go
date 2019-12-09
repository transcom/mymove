package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func Address(a *models.Address) *primemessages.Address {
	if a == nil {
		return nil
	}
	return &primemessages.Address{
		ID:             strfmt.UUID(a.ID.String()),
		StreetAddress1: &a.StreetAddress1,
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           &a.City,
		State:          &a.State,
		PostalCode:     &a.PostalCode,
		Country:        a.Country,
	}
}

func Customer(serviceMember *models.ServiceMember) *primemessages.Customer {
	if serviceMember == nil {
		return nil
	}
	var agency *string
	if serviceMember.Affiliation != nil {
		agency = handlers.FmtString(string(*serviceMember.Affiliation))
	}
	var rank *string
	if serviceMember.Rank != nil {
		rank = handlers.FmtString(string(*serviceMember.Rank))
	}

	return &primemessages.Customer{
		ID:            strfmt.UUID(serviceMember.ID.String()),
		Agency:        agency,
		Email:         serviceMember.PersonalEmail,
		FirstName:     serviceMember.FirstName,
		Grade:         rank,
		LastName:      serviceMember.LastName,
		MiddleName:    serviceMember.MiddleName,
		PickupAddress: Address(serviceMember.ResidentialAddress),
		Suffix:        serviceMember.Suffix,
		Telephone:     serviceMember.Telephone,
	}
}
