package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
)

func (suite *HandlerSuite) TestGetTransportationOfficesHandler() {
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
	}

	handler := GetTransportationOfficesHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesOK)

	// Validate outgoing payload
	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
}
