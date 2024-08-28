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
		usPostRegionCity := factory.BuildDefaultUsPostRegionCity(suite.DB())
		suite.MustSave(&usPostRegionCity)

		var fetchedUsPostRegionCity models.UsPostRegionCity
		err := suite.DB().Where("uspr_zip_id = $1", usPostRegionCity.UsprZipID).First(&fetchedUsPostRegionCity)

		suite.NoError(err)
		suite.Equal(usPostRegionCity.UsprZipID, fetchedUsPostRegionCity.UsprZipID)

		usPostRegionCityService := address.NewUsPostRegionCity()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/addresses/zip_city_lookup/"+usPostRegionCity.UsprZipID, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		params := addressop.GetLocationByZipCityParams{
			HTTPRequest: req,
			Search:      usPostRegionCity.UsprZipID,
		}

		handler := GetLocationByZipCityHandler{
			HandlerConfig:    suite.HandlerConfig(),
			UsPostRegionCity: usPostRegionCityService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(usPostRegionCity.UsprZipID, responsePayload.Payload[0].PostalCode)
	})
}
