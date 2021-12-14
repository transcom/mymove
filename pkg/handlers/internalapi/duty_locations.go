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

// TODO: temporary placeholder while we migrate away from duty_station
func payloadForDutyStationModel(station models.DutyStation) *internalmessages.DutyLocationPayload {
	// If the station ID has no UUID then it isn't real data
	// Unlike other payloads the
	if station.ID == uuid.Nil {
		return nil
	}
	payload := internalmessages.DutyLocationPayload{
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

func payloadForDutyLocationModel(station models.DutyLocation) *internalmessages.DutyLocationPayload {
	// If the station ID has no UUID then it isn't real data
	// Unlike other payloads the
	if station.ID == uuid.Nil {
		return nil
	}
	payload := internalmessages.DutyLocationPayload{
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

// SearchDutyLocationsHandler returns a list of all issues
type SearchDutyLocationsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of stations based on the search query
func (h SearchDutyLocationsHandler) Handle(params stationop.SearchDutyLocationsParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	stations, err := models.FindDutyLocations(appCtx.DB(), params.Search)
	if err != nil {
		appCtx.Logger().Error("Finding duty stations", zap.Error(err))
		return stationop.NewSearchDutyLocationsInternalServerError()

	}

	stationPayloads := make(internalmessages.DutyLocationsPayload, len(stations))
	for i, station := range stations {
		stationPayload := payloadForDutyLocationModel(station)
		stationPayloads[i] = stationPayload
	}
	return stationop.NewSearchDutyLocationsOK().WithPayload(stationPayloads)
}
