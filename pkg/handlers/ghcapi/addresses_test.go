package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	addressop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/addresses"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestGetLocationByZipCityHandler() {
	suite.Run("successful zip city lookup", func() {
		zip := "90210"
		var fetchedVLocation models.VLocation
		err := suite.DB().Where("uspr_zip_id = $1", zip).First(&fetchedVLocation)

		suite.NoError(err)
		suite.Equal(zip, fetchedVLocation.UsprZipID)

		vLocationService := address.NewVLocation()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/addresses/zip_city_lookup/"+zip, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := addressop.GetLocationByZipCityStateParams{
			HTTPRequest: req,
			Search:      zip,
		}

		handler := GetLocationByZipCityStateHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			VLocation:     vLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityStateOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityStateOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(zip, responsePayload.Payload[0].PostalCode)
	})
}

func (suite *HandlerSuite) TestCountrySearchHandler() {
	suite.Run("success", func() {
		countrySearcher := address.NewCountrySearcher()
		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
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
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
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
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
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

	suite.Run("forbidden", func() {
		mockCountrySearcher := mocks.CountrySearcher{}
		mockCountrySearcher.On("SearchCountries",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*string"),
		).Return(nil, apperror.QueryError{})

		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		notOfficeUser := factory.BuildUser(nil, nil, nil)
		req = suite.AuthenticateUserRequest(req, notOfficeUser)
		params := addressop.SearchCountriesParams{
			HTTPRequest: req,
			Search:      nil,
		}

		handler := SearchCountriesHandler{
			HandlerConfig:   suite.NewHandlerConfig(),
			CountrySearcher: &mockCountrySearcher,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.SearchCountriesForbidden{}, response)
	})
}
