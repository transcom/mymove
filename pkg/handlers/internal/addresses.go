package internal

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

/*
 * --------------------------------------------
 * The code below is for the INTERNAL REST API.
 * --------------------------------------------
 */

func addressModelFromPayload(rawAddress *internalmessages.Address) *models.Address {
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

func payloadForAddressModel(a *models.Address) *internalmessages.Address {
	if a == nil {
		return nil
	}
	return &internalmessages.Address{
		ID:             strfmt.UUID(a.ID.String()),
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}

func updateAddressWithPayload(a *models.Address, payload *internalmessages.Address) {
	a.StreetAddress1 = *payload.StreetAddress1
	a.StreetAddress2 = payload.StreetAddress2
	a.StreetAddress3 = payload.StreetAddress3
	a.City = *payload.City
	a.State = *payload.State
	a.PostalCode = *payload.PostalCode
	a.Country = payload.Country
}

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

func publicPayloadForAddressModel(a *models.Address) *apimessages.Address {
	if a == nil {
		return nil
	}
	return &apimessages.Address{
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          *swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}
