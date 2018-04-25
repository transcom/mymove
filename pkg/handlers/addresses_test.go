package handlers

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func fakeAddressPayload() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: swag.String("An address"),
		StreetAddress2: swag.String("Apt. 2"),
		StreetAddress3: swag.String("address line 3"),
		City:           swag.String("Happytown"),
		State:          swag.String("AL"),
		PostalCode:     swag.String("01234"),
	}
}
