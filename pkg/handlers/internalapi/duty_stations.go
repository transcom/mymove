package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	stationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

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
