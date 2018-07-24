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

func otherFakeAddressPayload() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: swag.String("The street where you live"),
		StreetAddress2: swag.String("Flat 3"),
		City:           swag.String("Ninoville"),
		State:          swag.String("AL"),
		PostalCode:     swag.String("32145"),
	}
}
