package primeapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	addressop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/addresses"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestCountrySearchHandler() {
	suite.Run("success", func() {
		countrySearcher := address.NewCountrySearcher()
		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		params := addressop.SearchCountriesParams{
			HTTPRequest: req,
			Search:      models.StringPointer("us"),
		}

		handler := SearchCountriesHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			CountrySearcher: countrySearcher}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.SearchCountriesOK{}, response)
		responsePayload := response.(*addressop.SearchCountriesOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.True(len(responsePayload.Payload) == 2)
	})

	suite.Run("success - return all", func() {
		countrySearcher := address.NewCountrySearcher()
		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		params := addressop.SearchCountriesParams{
			HTTPRequest: req,
			Search:      nil,
		}

		handler := SearchCountriesHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			CountrySearcher: countrySearcher}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.SearchCountriesOK{}, response)
		responsePayload := response.(*addressop.SearchCountriesOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.True(len(responsePayload.Payload) == 274)
	})

	suite.Run("failure", func() {
		mockCountrySearcher := mocks.CountrySearcher{}
		mockCountrySearcher.On("SearchCountries",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*string"),
		).Return(nil, apperror.QueryError{})

		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		params := addressop.SearchCountriesParams{
			HTTPRequest: req,
			Search:      nil,
		}

		handler := SearchCountriesHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			CountrySearcher: &mockCountrySearcher,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.SearchCountriesInternalServerError{}, response)
	})
}
