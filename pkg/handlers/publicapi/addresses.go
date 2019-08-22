package publicapi

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
)

func addressModelFromPayload(rawAddress *apimessages.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	return &models.Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		StreetAddress3: rawAddress.StreetAddress3,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		PostalCode:     *rawAddress.PostalCode,
		Country:        rawAddress.Country,
	}
}

func payloadForAddressModel(a *models.Address) *apimessages.Address {
	if a == nil {
		return nil
	}
	if a.ID == uuid.Nil {
		return nil
	}

	return &apimessages.Address{
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}
