package ghcapi

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	addressop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/addresses"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type GetLocationByZipCityStateHandler struct {
	handlers.HandlerConfig
	services.VLocation
}

func (h GetLocationByZipCityStateHandler) Handle(params addressop.GetLocationByZipCityStateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeApp() && appCtx.Session().OfficeUserID == uuid.Nil {
				noOfficeUserIDErr := apperror.NewSessionError("No office user ID")
				return addressop.NewGetLocationByZipCityStateForbidden(), noOfficeUserIDErr
			}

			/** Feature Flag - Alaska - Determines if AK can be included/excluded **/
			isAlaskaEnabled := false
			akFeatureFlagName := "enable_alaska"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, akFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", akFeatureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			/** Feature Flag - Hawaii - Determines if HI can be included/excluded **/
			isHawaiiEnabled := false
			hiFeatureFlagName := "enable_hawaii"
			flag, err = h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, hiFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", hiFeatureFlagName), zap.Error(err))
			} else {
				isHawaiiEnabled = flag.Match
			}

			// build states to exlude filter list
			statesToExclude := make([]string, 0)
			if !isAlaskaEnabled {
				statesToExclude = append(statesToExclude, "AK")
			}
			if !isHawaiiEnabled {
				statesToExclude = append(statesToExclude, "HI")
			}

			locationList, err := h.GetLocationsByZipCityState(appCtx, params.Search, statesToExclude)
			if err != nil {
				appCtx.Logger().Error("Error searching for Zip/City/State: ", zap.Error(err))
				return addressop.NewGetLocationByZipCityStateInternalServerError(), err
			}

			returnPayload := payloads.VLocations(*locationList)
			return addressop.NewGetLocationByZipCityStateOK().WithPayload(returnPayload), nil
		})
}

// SearchCountriesHandler returns a list of countries
type SearchCountriesHandler struct {
	handlers.HandlerConfig
	services.CountrySearcher
}

// Handle returns a list of locations based on the search query
func (h SearchCountriesHandler) Handle(params addressop.SearchCountriesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeApp() && appCtx.Session().OfficeUserID == uuid.Nil {
				noOfficeUserIDErr := apperror.NewSessionError("No office user ID")
				return addressop.NewSearchCountriesForbidden(), noOfficeUserIDErr
			}

			/** Feature Flag - Determines if countries should return US or all countries **/
			isOconusCityFinder := false
			oconusFeatureFlagName := "oconus_city_finder"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, oconusFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", oconusFeatureFlagName), zap.Error(err))
			} else {
				isOconusCityFinder = flag.Match
			}

			var countries models.Countries

			if !isOconusCityFinder {
				country, err := models.FetchCountryByCode(appCtx.DB(), "US")

				if err != nil {
					appCtx.Logger().Error("Error searching for countries: ", zap.Error(err))
					return addressop.NewSearchCountriesInternalServerError(), err
				}

				countries = append(countries, country)
			} else {
				countries, err = h.CountrySearcher.SearchCountries(appCtx, params.Search)
				if err != nil {
					appCtx.Logger().Error("Error searching for countries: ", zap.Error(err))
					return addressop.NewSearchCountriesInternalServerError(), err
				}
			}

			returnPayload := payloads.Countries(countries)
			return addressop.NewSearchCountriesOK().WithPayload(returnPayload), nil
		})
}
