package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	countryop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/countries"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// SearchCountriesHandler returns a list of countries
type SearchCountriesHandler struct {
	handlers.HandlerConfig
	services.CountrySearcher
}

// Handle returns a list of locations based on the search query
func (h SearchCountriesHandler) Handle(params countryop.SearchCountriesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			countries, err := h.CountrySearcher.SearchCountries(appCtx, params.Search)
			if err != nil {
				return nil, err
			}

			returnPayload := payloads.VCountries(countries)
			return countryop.NewSearchCountriesOK().WithPayload(returnPayload), nil
		})
}
