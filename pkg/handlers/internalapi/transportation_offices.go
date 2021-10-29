package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ShowDutyStationTransportationOfficeHandler returns the transportation office for a duty station ID
type ShowDutyStationTransportationOfficeHandler struct {
	handlers.HandlerContext
}

// Handle retrieves the transportation office in the system for a given duty station ID
func (h ShowDutyStationTransportationOfficeHandler) Handle(params transportationofficeop.ShowDutyStationTransportationOfficeParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	dutyStationID, _ := uuid.FromString(params.DutyStationID.String())
	transportationOffice, err := models.FetchDutyStationTransportationOffice(appCtx.DB(), dutyStationID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	transportationOfficePayload := payloads.TransportationOffice(transportationOffice)

	return transportationofficeop.NewShowDutyStationTransportationOfficeOK().WithPayload(transportationOfficePayload)
}
