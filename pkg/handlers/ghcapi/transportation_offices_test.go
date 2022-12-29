package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	"github.com/transcom/mymove/pkg/models"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetTransportationOfficesHandler() {
	transportationOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Name: "LRC Fort Knox",
		},
	})

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "LRC Fort Know",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesOK)

	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
	suite.Equal(transportationOffice.Name, *responsePayload.Payload[0].Name)
	suite.Equal(transportationOffice.Address.ID.String(), responsePayload.Payload[0].Address.ID.String())
	suite.Equal(transportationOffice.Gbloc, responsePayload.Payload[0].Gbloc)

}

func (suite *HandlerSuite) TestNoTransportationOfficesHandler() {

	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "LRC Fort Know",
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload, ok := response.(*transportationofficeop.GetTransportationOfficesOK)

	suite.True(ok)
	suite.NotNil(responsePayload, "Response should not be nil")

}
