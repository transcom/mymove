package primeapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
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

	suite.Run("returns no results for a PO box zip when PO boxes are excluded", func() {
		zip := "00929" // PO Box ZIP in PR
		var fetchedVLocation models.VLocation
		err := suite.DB().Where("uspr_zip_id = $1", zip).First(&fetchedVLocation)

		suite.NoError(err)
		suite.Equal(zip, fetchedVLocation.UsprZipID)

		vLocationService := address.NewVLocation()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/addresses/zip_city_lookup/"+zip, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := addressop.GetLocationByZipCityStateParams{
			HTTPRequest:    req,
			Search:         zip,
			IncludePOBoxes: models.BoolPointer(false),
		}

		handler := GetLocationByZipCityStateHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			VLocation:     vLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityStateOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityStateOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(0, len(responsePayload.Payload))
	})
}

func (suite *HandlerSuite) TestGetOconusLocationHandler() {
	suite.Run("successful city name lookup", func() {
		country := "GB"
		city := "LONDON"
		cityResult := "LONDON COLNEY"
		var fetchedVIntlLocation models.VIntlLocation
		err := suite.DB().Where("city_name = $1", city).First(&fetchedVIntlLocation)

		suite.NoError(err)
		suite.Equal(city, *fetchedVIntlLocation.CityName)

		vIntlLocationService := address.NewVIntlLocation()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/addresses/oconus_lookup/"+country+"/"+city, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := addressop.GetOconusLocationParams{
			HTTPRequest: req,
			Country:     country,
			Search:      city,
		}

		handler := GetOconusLocationHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			VIntlLocation: vIntlLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetOconusLocationOK{}, response)
		responsePayload := response.(*addressop.GetOconusLocationOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(cityResult, responsePayload.Payload[0].City)
	})
}
