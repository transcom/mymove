package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"

	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyLocation) *internalmessages.DutyStationPayload {
	// If the station ID has no UUID then it isn't real data
	// Unlike other payloads the
	if station.ID == uuid.Nil {
		return nil
	}
	payload := internalmessages.DutyStationPayload{
		ID:          handlers.FmtUUID(station.ID),
		CreatedAt:   handlers.FmtDateTime(station.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(station.UpdatedAt),
		Name:        swag.String(station.Name),
		Affiliation: station.Affiliation,
		AddressID:   handlers.FmtUUID(station.AddressID),
		Address:     payloads.Address(&station.Address),
	}

	payload.TransportationOffice = payloads.TransportationOffice(station.TransportationOffice)

	return &payload
}

// SearchDutyStationsHandler returns a list of all issues
type SearchDutyStationsHandler struct {
	handlers.HandlerContext
}

// TODO: temporary workaround until this file gets deleted entirely
// Handle returns a list of stations based on the search query
func (h SearchDutyStationsHandler) Handle(params stationop.SearchDutyLocationsParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	locations, err := models.FindDutyLocations(appCtx.DB(), params.Search)
	if err != nil {
		appCtx.Logger().Error("Finding duty stations", zap.Error(err))
		return stationop.NewSearchDutyLocationsInternalServerError()

	}

	locationPayloads := make(internalmessages.DutyLocationsPayload, len(locations))
	for i, location := range locations {
		locationPayload := payloadForDutyLocationModel(location)
		locationPayloads[i] = locationPayload
	}
	return stationop.NewSearchDutyLocationsOK().WithPayload(locationPayloads)
}
