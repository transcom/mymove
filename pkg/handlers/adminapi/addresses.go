package adminapi

import (
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

// Country payload
func Country(country *models.Country) *string {
	if country == nil {
		return nil
	}
	return &country.Country
}

func payloadForAddressModel(a *models.Address) *adminmessages.Address {
	if a == nil {
		return nil
	}
	return &adminmessages.Address{
		StreetAddress1: models.StringPointer(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           models.StringPointer(a.City),
		State:          models.StringPointer(a.State),
		PostalCode:     models.StringPointer(a.PostalCode),
		Country:        Country(a.Country),
		County:         &a.County,
	}
}
