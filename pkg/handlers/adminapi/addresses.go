package adminapi

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForAddressModel(a *models.Address) *adminmessages.Address {
	if a == nil {
		return nil
	}
	return &adminmessages.Address{
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}
