package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	addressop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/addresses"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
)

func fakeAddressPayload() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: models.StringPointer("An address"),
		StreetAddress2: models.StringPointer("Apt. 2"),
		StreetAddress3: models.StringPointer("address line 3"),
		City:           models.StringPointer("Happytown"),
		State:          models.StringPointer("AL"),
		PostalCode:     models.StringPointer("40356"),
		County:         models.StringPointer("JESSAMINE"),
		IsOconus:       models.BoolPointer(false),
	}
}

func (suite *HandlerSuite) TestShowAddressHandler() {

	suite.Run("successful lookup", func() {
		address := models.Address{
			StreetAddress1: "some address",
			City:           "city",
			State:          "state",
			PostalCode:     "12345",
			County:         models.StringPointer("JESSAMINE"),
			IsOconus:       models.BoolPointer(false),
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

			handler := ShowAddressHandler{suite.HandlerConfig()}
			res := handler.Handle(params)

			response := res.(*addressop.ShowAddressOK)
			payload := response.Payload

			if ts.hasResult {
				suite.NotNil(payload, "Should have address record")
				suite.Equal(payload.ID.String(), ts.resultID, "Address ID doest match")
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
			HandlerConfig: suite.HandlerConfig(),
			VLocation:     vLocationServices}

		response := handler.Handle(params)
		suite.Assertions.IsType(&addressop.GetLocationByZipCityStateOK{}, response)
		responsePayload := response.(*addressop.GetLocationByZipCityStateOK)
		suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
		suite.Equal(zip, responsePayload.Payload[0].PostalCode)
	})
}

func (suite *HandlerSuite) TestGetOconusLocationHandler() {
	suite.Run("successful city name lookup", func() {
		country := "GB"
		city := "LONDON"
		var fetchedVIntlLocation models.VIntlLocation
		err := suite.DB().Where("city_name = $1", city).First(&fetchedVIntlLocation)

		suite.NoError(err)
		suite.Equal(city, *fetchedVIntlLocation.CityName)

		vIntlLocationService := address.NewVIntlLocation()

		move := factory.BuildMove(suite.DB(), nil, nil)
		req := httptest.NewRequest("GET", "/addresses/oconus_lookup/"+country+"/"+city, nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
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
		suite.Equal(city, responsePayload.Payload[0].City)
	})
}
