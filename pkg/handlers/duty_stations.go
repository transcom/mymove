package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyStation) internalmessages.DutyStationPayload {

	stationPayload := internalmessages.DutyStationPayload{
		ID:          fmtUUID(station.ID),
		CreatedAt:   fmtDateTime(station.CreatedAt),
		UpdatedAt:   fmtDateTime(station.UpdatedAt),
		Name:        swag.String(station.Name),
		Affiliation: &station.Affiliation,
		Address:     payloadForAddressModel(&station.Address),
	}
	return stationPayload
}

// SearchDutyStationsHandler returns a list of all issues
type SearchDutyStationsHandler HandlerContext

// Handle returns a list of stations based on the search query
func (h SearchDutyStationsHandler) Handle(params stationop.SearchDutyStationsParams) middleware.Responder {
	var stations models.DutyStations
	var response middleware.Responder
	var err error

	stations, err = models.FindDutyStations(h.db, params.Search, params.Affiliation)
	if err != nil {
		h.logger.Error("Finding duty stations", zap.Error(err))
		response = stationop.NewSearchDutyStationsInternalServerError()
	}

	stationPayloads := make(internalmessages.DutyStationsPayload, len(stations))
	for i, station := range stations {
		stationPayload := payloadForDutyStationModel(station)
		stationPayloads[i] = &stationPayload
	}
	response = stationop.NewSearchDutyStationsOK().WithPayload(stationPayloads)

	return response
}
