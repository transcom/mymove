package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ShowDutyLocationTransportationOfficeHandler returns the transportation office for a duty location ID
type ShowDutyLocationTransportationOfficeHandler struct {
	handlers.HandlerContext
}

// Handle retrieves the transportation office in the system for a given duty location ID
func (h ShowDutyLocationTransportationOfficeHandler) Handle(params transportationofficeop.ShowDutyLocationTransportationOfficeParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			dutyLocationID, _ := uuid.FromString(params.DutyLocationID.String())
			transportationOffice, err := models.FetchDutyLocationTransportationOffice(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}
			transportationOfficePayload := payloads.TransportationOffice(transportationOffice)

			return transportationofficeop.NewShowDutyLocationTransportationOfficeOK().WithPayload(transportationOfficePayload)
		})
}
