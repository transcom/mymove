package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyStation) *internalmessages.DutyStationPayload {
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
		Affiliation: &station.Affiliation,
		Address:     payloadForAddressModel(&station.Address),
	}

	if station.TransportationOfficeID != nil {
		payload.TransportationOffice = payloadForTransportationOfficeModel(station.TransportationOffice)
	}

	return &payload
}

// SearchDutyStationsHandler returns a list of all issues
type SearchDutyStationsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of stations based on the search query
func (h SearchDutyStationsHandler) Handle(params stationop.SearchDutyStationsParams) middleware.Responder {
	var stations models.DutyStations
	var err error

	stations, err = models.FindDutyStations(h.DB(), params.Search)
	if err != nil {
		h.Logger().Error("Finding duty stations", zap.Error(err))
		return stationop.NewSearchDutyStationsInternalServerError()

	}

	stationPayloads := make(internalmessages.DutyStationsPayload, len(stations))
	for i, station := range stations {
		stationPayload := payloadForDutyStationModel(station)
		stationPayloads[i] = stationPayload
	}
	return stationop.NewSearchDutyStationsOK().WithPayload(stationPayloads)
}
