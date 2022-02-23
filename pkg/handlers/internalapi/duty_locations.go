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

func payloadForDutyLocationModel(location models.DutyLocation) *internalmessages.DutyLocationPayload {
	// If the location ID has no UUID then it isn't real data
	// Unlike other payloads the
	if location.ID == uuid.Nil {
		return nil
	}
	payload := internalmessages.DutyLocationPayload{
		ID:          handlers.FmtUUID(location.ID),
		CreatedAt:   handlers.FmtDateTime(location.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(location.UpdatedAt),
		Name:        swag.String(location.Name),
		Affiliation: location.Affiliation,
		AddressID:   handlers.FmtUUID(location.AddressID),
		Address:     payloads.Address(&location.Address),
	}

	return &payload
}

// SearchDutyLocationsHandler returns a list of all issues
type SearchDutyLocationsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of locations based on the search query
func (h SearchDutyLocationsHandler) Handle(params stationop.SearchDutyLocationsParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	locations, err := models.FindDutyLocations(appCtx.DB(), params.Search)
	if err != nil {
		appCtx.Logger().Error("Finding duty locations", zap.Error(err))
		return stationop.NewSearchDutyLocationsInternalServerError()

	}

	locationPayloads := make(internalmessages.DutyLocationsPayload, len(locations))
	for i, location := range locations {
		locationPayload := payloadForDutyLocationModel(location)
		locationPayloads[i] = locationPayload
	}
	return stationop.NewSearchDutyLocationsOK().WithPayload(locationPayloads)
}
