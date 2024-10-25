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
		params := addressop.GetLocationByZipCityParams{
			HTTPRequest: req,
			Search:      zip,
		}

		handler := GetLocationByZipCityHandler{
			HandlerConfig: suite.HandlerConfig(),
			VLocation:     vLocationService}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(zip, responsePayload.Payload[0].PostalCode)
	})
}
