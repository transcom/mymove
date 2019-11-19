package payloads

import (
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
)

func AddressModel(address *primemessages.Address) *models.Address {
	if address == nil {
		return nil
	}
	return &models.Address{
		StreetAddress1: *address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           *address.City,
		State:          *address.State,
		PostalCode:     *address.PostalCode,
		Country:        address.Country,
	}
}
