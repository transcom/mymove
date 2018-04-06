package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_stations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyStation) internalmessages.DutyStationPayload {
	stationPayload := internalmessages.DutyStationPayload{
		ID:        fmtUUID(station.ID),
		CreatedAt: fmtDateTime(station.CreatedAt),
		UpdatedAt: fmtDateTime(station.UpdatedAt),
		Name:      swag.String(station.Name),
		Branch:    station.Branch,
		Address:   payloadForAddressModel(&station.Address),
	}
	return stationPayload
}

// SearchDutyStationsHandler returns a list of all issues
type SearchDutyStationsHandler HandlerContext

// Handle returns a list of stations based on the search query
func (h SearchDutyStationsHandler) Handle(params stationop.SearchDutyStationsParams) middleware.Responder {
	var stations models.DutyStations
	var response middleware.Responder

	// Verify user is logged in
	_, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response = stationop.NewSearchDutyStationsUnauthorized()
		return response
	}

	// ILIKE does case-insensitive pattern matching, "%" matches any string
	searchQuery := fmt.Sprintf("%%%s%%", params.Search)
	query := h.db.Where("branch = ? AND name ILIKE ?", params.Branch, searchQuery)

	if err := query.Eager().All(&stations); err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			response = stationop.NewSearchDutyStationsNotFound()
		} else {
			h.logger.Error("DB Query", zap.Error(err))
			response = stationop.NewSearchDutyStationsInternalServerError()
		}
	} else {
		stationPayloads := make(internalmessages.DutyStationsPayload, len(stations))
		for i, station := range stations {
			stationPayload := payloadForDutyStationModel(station)
			stationPayloads[i] = &stationPayload
		}
		response = stationop.NewSearchDutyStationsOK().WithPayload(stationPayloads)
	}
	return response
}
