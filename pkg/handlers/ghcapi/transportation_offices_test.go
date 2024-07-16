package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	"github.com/transcom/mymove/pkg/models"
	transportationofficeservice "github.com/transcom/mymove/pkg/services/transportation_office"
)

func (suite *HandlerSuite) TestGetTransportationOfficesHandler() {
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "LRC Fort Knox",
				ProvidesCloseout: true,
			},
		},
	}, nil)
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices", nil)
	params := transportationofficeop.GetTransportationOfficesParams{
		HTTPRequest: req,
		Search:      "LRC Fort Knox",
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
		Search:      "LRC Fort Knox",
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

func (suite *HandlerSuite) TestGetTransportationOfficesGBLOCsHandler() {
	transportationOffice1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "LRC Fort Knox",
				ProvidesCloseout: true,
				Gbloc:            "AGFM",
			},
		},
	}, nil)
	transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:             "NSF End of Alphabet",
				ProvidesCloseout: true,
				Gbloc:            "WXYZ",
			},
		},
	}, nil)
	fetcher := transportationofficeservice.NewTransportationOfficesFetcher()

	req := httptest.NewRequest("GET", "/transportation_offices/gblocs", nil)
	params := transportationofficeop.GetTransportationOfficesGBLOCsParams{
		HTTPRequest: req,
	}

	handler := GetTransportationOfficesGBLOCsHandler{
		HandlerConfig:                suite.HandlerConfig(),
		TransportationOfficesFetcher: fetcher,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&transportationofficeop.GetTransportationOfficesGBLOCsOK{}, response)
	responsePayload := response.(*transportationofficeop.GetTransportationOfficesGBLOCsOK)

	suite.NoError(responsePayload.Payload.Validate(strfmt.Default))
	suite.Equal(transportationOffice1.Gbloc, responsePayload.Payload[0])
	suite.Equal(transportationOffice2.Gbloc, responsePayload.Payload[1])
}
