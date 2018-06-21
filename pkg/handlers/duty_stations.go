package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyStation) *internalmessages.DutyStationPayload {
	payload := internalmessages.DutyStationPayload{
		ID:          fmtUUID(station.ID),
		CreatedAt:   fmtDateTime(station.CreatedAt),
		UpdatedAt:   fmtDateTime(station.UpdatedAt),
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
type SearchDutyStationsHandler HandlerContext

// Handle returns a list of stations based on the search query
func (h SearchDutyStationsHandler) Handle(params stationop.SearchDutyStationsParams) middleware.Responder {
	var stations models.DutyStations
	var err error

	stations, err = models.FindDutyStations(h.db, params.Search)
	if err != nil {
		h.logger.Error("Finding duty stations", zap.Error(err))
		return stationop.NewSearchDutyStationsInternalServerError()

	}

	stationPayloads := make(internalmessages.DutyStationsPayload, len(stations))
	for i, station := range stations {
		stationPayload := payloadForDutyStationModel(station)
		stationPayloads[i] = stationPayload
	}
	return stationop.NewSearchDutyStationsOK().WithPayload(stationPayloads)
}
