package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	addressop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/addresses"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
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
			HandlerConfig: suite.HandlerConfig(),
			VLocation:     vLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityStateOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityStateOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(zip, responsePayload.Payload[0].PostalCode)
	})
}

func (suite *HandlerSuite) TestGetOconusLocationHandler() {
	country := "GB"
	city := "LONDON"

	suite.Run("successful city name lookup", func() {
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
			HandlerConfig: suite.HandlerConfig(),
			VIntlLocation: vIntlLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetOconusLocationOK{}, response)
		responsePayload := response.(*addressop.GetOconusLocationOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(cityResult, responsePayload.Payload[0].City)
	})

	suite.Run("forbidden", func() {
		vIntlLocationService := address.NewVIntlLocation()

		req := httptest.NewRequest("GET", "/addresses/oconus_lookup/"+country+"/"+city, nil)
		notOfficeUser := factory.BuildUser(nil, nil, nil)
		req = suite.AuthenticateUserRequest(req, notOfficeUser)
		params := addressop.GetOconusLocationParams{
			HTTPRequest: req,
			Country:     country,
			Search:      city,
		}

		handler := GetOconusLocationHandler{
			HandlerConfig: suite.HandlerConfig(),
			VIntlLocation: vIntlLocationService,
		}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetOconusLocationForbidden{}, response)
	})
}
