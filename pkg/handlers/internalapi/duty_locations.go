package internalapi

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	locationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/duty_locations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyLocationModel(location models.DutyLocation) *internalmessages.DutyLocationPayload {
	// If the location ID has no UUID then it isn't real data
	// Unlike other payloads the
	if location.ID == uuid.Nil {
		return nil
	}
	payload := internalmessages.DutyLocationPayload{
		ID:                         handlers.FmtUUID(location.ID),
		CreatedAt:                  handlers.FmtDateTime(location.CreatedAt),
		UpdatedAt:                  handlers.FmtDateTime(location.UpdatedAt),
		Name:                       models.StringPointer(location.Name),
		Affiliation:                location.Affiliation,
		AddressID:                  handlers.FmtUUID(location.AddressID),
		Address:                    payloads.Address(&location.Address),
		TransportationOfficeID:     handlers.FmtUUIDPtr(location.TransportationOfficeID),
		ProvidesServicesCounseling: location.ProvidesServicesCounseling,
	}
	payload.TransportationOffice = payloads.TransportationOffice(location.TransportationOffice)

	return &payload
}

// SearchDutyLocationsHandler returns a list of all issues
type SearchDutyLocationsHandler struct {
	handlers.HandlerConfig
}

// Handle returns a list of locations based on the search query
func (h SearchDutyLocationsHandler) Handle(params locationop.SearchDutyLocationsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - Alaska - Determines if AK can be included/excluded **/
			isAlaskaEnabled := false
			featureFlagName := "enable_alaska"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			statesToExclude := make([]string, 0)
			if !isAlaskaEnabled {
				// HI locations will also be excluded as part of AK embargo
				statesToExclude = append(statesToExclude, "AK")
				statesToExclude = append(statesToExclude, "HI")
			}

			locations, err := models.FindDutyLocationsExcludingStates(appCtx.DB(), params.Search, statesToExclude)
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
