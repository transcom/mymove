package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	addressop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

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

// ShowAddressHandler returns an address
type ShowAddressHandler struct {
	handlers.HandlerContext
}

// Handle returns a address given an addressId
func (h ShowAddressHandler) Handle(params addressop.ShowAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	addressID, err := uuid.FromString(params.AddressID.String())

	address := models.FetchAddressByID(h.DB(), &addressID)
	if err != nil {
		logger.Error("Finding address", zap.Error(err))
	}

	addressPayload := payloadForAddressModel(address)
	return addressop.NewShowAddressOK().WithPayload(addressPayload)
}
