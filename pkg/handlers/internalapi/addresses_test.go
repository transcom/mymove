package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	addressop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func fakeAddressPayload(country *models.Country) *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: models.StringPointer("An address"),
		StreetAddress2: models.StringPointer("Apt. 2"),
		StreetAddress3: models.StringPointer("address line 3"),
		City:           models.StringPointer("NICHOLASVILLE"),
		State:          models.StringPointer("AL"),
		PostalCode:     models.StringPointer("40356"),
		County:         models.StringPointer("JESSAMINE"),
		IsOconus:       models.BoolPointer(false),
		CountryID:      strfmt.UUID(country.ID.String()),
		Country:        payloads.Country(country),
	}
}

func (suite *HandlerSuite) TestShowAddressHandler() {

	suite.Run("successful lookup", func() {
		fetchedUsPostRegionCity, err := models.FindByZipCodeAndCity(suite.DB(), "12345", "SCHENECTADY")
		suite.NoError(err)
		address := models.Address{
			StreetAddress1:     "some address",
			City:               fetchedUsPostRegionCity.USPostRegionCityNm,
			State:              "NY",
			PostalCode:         fetchedUsPostRegionCity.UsprZipID,
			County:             models.StringPointer("JESSAMINE"),
			IsOconus:           models.BoolPointer(false),
			UsPostRegionCityID: &fetchedUsPostRegionCity.ID,
		}
		suite.MustSave(&address)

		requestUser := factory.BuildUser(nil, nil, nil)

		fakeUUID, _ := uuid.FromString("not-valid-uuid")

		tests := []struct {
			ID        uuid.UUID
			hasResult bool
			resultID  string
		}{
			{ID: address.ID, hasResult: true, resultID: address.ID.String()},
			{ID: fakeUUID, hasResult: false, resultID: ""},
		}

		for _, ts := range tests {
			req := httptest.NewRequest("GET", "/addresses/"+ts.ID.String(), nil)
			req = suite.AuthenticateUserRequest(req, requestUser)

			params := addressop.ShowAddressParams{
				HTTPRequest: req,
				AddressID:   *handlers.FmtUUID(ts.ID),
			}

			handler := ShowAddressHandler{suite.NewHandlerConfig()}
			res := handler.Handle(params)

			response := res.(*addressop.ShowAddressOK)
			payload := response.Payload

			if ts.hasResult {
				suite.NotNil(payload, "Should have address record")
				suite.Equal(payload.ID.String(), ts.resultID, "Address ID doest match")
				suite.Equal(payload.UsPostRegionCitiesID.String(), fetchedUsPostRegionCity.ID.String())
			} else {
				suite.Nil(payload, "Should not have address record")
			}
		}
	})

}

func (suite *HandlerSuite) TestGetLocationByZipCityHandler() {
	suite.Run("successful zip city lookup", func() {
		zip := "90210"
		var fetchedVLocation models.VLocation
		err := suite.DB().Where("uspr_zip_id = $1", zip).First(&fetchedVLocation)

		suite.NoError(err)
		suite.Equal(zip, fetchedVLocation.UsprZipID)

		vLocationServices := address.NewVLocation()
		move := factory.BuildMove(suite.DB(), nil, nil)
		req := httptest.NewRequest("GET", "/addresses/zip_city_lookup/"+zip, nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		params := addressop.GetLocationByZipCityStateParams{
			HTTPRequest: req,
			Search:      zip,
		}

		handler := GetLocationByZipCityStateHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			VLocation:     vLocationServices}

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

		vLocationServices := address.NewVLocation()
		move := factory.BuildMove(suite.DB(), nil, nil)
		req := httptest.NewRequest("GET", "/addresses/zip_city_lookup/"+zip, nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		params := addressop.GetLocationByZipCityStateParams{
			HTTPRequest:    req,
			Search:         zip,
			IncludePOBoxes: models.BoolPointer(false),
		}

		handler := GetLocationByZipCityStateHandler{
			HandlerConfig: suite.NewHandlerConfig(),
			VLocation:     vLocationServices}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityStateOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityStateOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(0, len(responsePayload.Payload))
	})
}

func (suite *HandlerSuite) TestCountrySearchHandler() {
	suite.Run("success", func() {
		countrySearcher := address.NewCountrySearcher()
		req := httptest.NewRequest("GET", "/addresses/countries", nil)
		serviceMember := factory.BuildDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, serviceMember)
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
		serviceMember := factory.BuildDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, serviceMember)
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
		serviceMember := factory.BuildDefaultUser(suite.DB())
		req = suite.AuthenticateUserRequest(req, serviceMember)
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
		suite.Assertions.IsType(&addressop.SearchCountriesForbidden{}, response)
	})
}
