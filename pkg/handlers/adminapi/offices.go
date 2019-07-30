package adminapi

import (
	"github.com/go-openapi/runtime/middleware"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForOfficeModel(o models.TransportationOffice) *adminmessages.TransportationOffice {
	return &adminmessages.TransportationOffice{
		ID:         handlers.FmtUUID(o.ID),
		Name:       handlers.FmtString(o.Name),
		Address:    payloadForAddressModel(&o.Address),
		Gbloc:      o.Gbloc,
		PhoneLines: payloadForPhoneLines(o.PhoneLines),
		Latitude:   o.Latitude,
		Longitude:  o.Longitude,
	}
}

// IndexOfficesHandler returns a list of office users via GET /office_users
type IndexOfficesHandler struct {
	handlers.HandlerContext
	services.NewQueryFilter
}

// Handle retrieves a list of office users
func (h IndexOfficesHandler) Handle(params officeop.IndexOfficesParams) middleware.Responder {
	return officeop.NewIndexOfficesOK().WithPayload(adminmessages.TransportationOffices{})
}
