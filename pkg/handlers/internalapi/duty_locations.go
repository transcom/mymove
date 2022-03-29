package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"

	locationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
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
		ID:                     handlers.FmtUUID(location.ID),
		CreatedAt:              handlers.FmtDateTime(location.CreatedAt),
		UpdatedAt:              handlers.FmtDateTime(location.UpdatedAt),
		Name:                   swag.String(location.Name),
		Affiliation:            location.Affiliation,
		AddressID:              handlers.FmtUUID(location.AddressID),
		Address:                payloads.Address(&location.Address),
		TransportationOfficeID: handlers.FmtUUIDPtr(location.TransportationOfficeID),
	}
	payload.TransportationOffice = payloads.TransportationOffice(location.TransportationOffice)

	return &payload
}

// SearchDutyLocationsHandler returns a list of all issues
type SearchDutyLocationsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of locations based on the search query
func (h SearchDutyLocationsHandler) Handle(params locationop.SearchDutyLocationsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			locations, err := models.FindDutyLocations(appCtx.DB(), params.Search)
			if err != nil {
				dutyLocationErr := apperror.NewNotFoundError(uuid.Nil, "Finding duty locations")
				appCtx.Logger().Error(dutyLocationErr.Error(), zap.Error(err))
				return locationop.NewSearchDutyLocationsInternalServerError(), dutyLocationErr

			}

			locationPayloads := make(internalmessages.DutyLocationsPayload, len(locations))
			for i, location := range locations {
				locationPayload := payloadForDutyLocationModel(location)
				locationPayloads[i] = locationPayload
			}
			return locationop.NewSearchDutyLocationsOK().WithPayload(locationPayloads), nil
		})
}
